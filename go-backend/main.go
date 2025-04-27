package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	users = make(map[string]string) // 内存中的用户数据
	mu    sync.Mutex
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	// 启动时加载用户数据
	loadUsersFromFile()
}

func main() {
	r := gin.Default()

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/register", register)
	r.POST("/login", login)

	r.Run(":8089")
}

// 加载用户数据
func loadUsersFromFile() {
	data, err := os.ReadFile("users.json")
	if err == nil {
		json.Unmarshal(data, &users)
	}
}

// 保存用户数据到文件
func saveUsersToFile() {
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile("users.json", data, 0644)
}

// 注册逻辑
func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[user.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	users[user.Username] = string(hashedPassword)
	saveUsersToFile()

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// 登录逻辑
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	storedPassword, exists := users[user.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}
