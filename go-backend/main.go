package main

import (
	invoke_fabric "backend/fabric-go/call"
	connect_fabric "backend/fabric-go/network"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	users = make(map[string]string) // 内存中的用户数据
	mu    sync.Mutex
)

type User struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Organization string   `json:"organization"`
	Pubkeyhash   string   `json:"pubkeyhash"`
	Token        int      `json:"token"`
	Posted       []string `json:"posted"`
	Accepted     []string `json:"posted"`
	IsAdmin      bool     `json:"isAdmin"`
	IsVerified   bool     `json:"isVerified"`
	IsAccepted   bool     `json:"isAccepted"`
}

type Model struct {
	Modelid    string `json:"Modelid"`
	Modelowner string `json:"Modelowner"`
	Modelhash  string `json:"Modelhash"`
	Modelsign  string `json:"Modelsign"`
}
type Task struct {
	TaskID        string   `json:"ID"`
	Bonus         int      `json:"bonus"`
	RootModelId   string   `json:"rootModelHash"`
	PostedUser    string   `json:"postedUser"`
	AcceptedUsers []string `json:"acceptedUsers"`
	Models        []string `json:"models"`
	IsComplete    bool     `json:"isComplete"`
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
func register(ctx *gin.Context) {
	var user User
	contract := connect_fabric.GetContract()

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}
	print(user.Username)
	print(user.Password)
	print(user.Organization)
	// 调用 QueryUser 并处理返回值
	err := invoke_fabric.CreateNewUser(contract, user.Username, user.Password, user.Organization, "test", 0, true, true, true)
	if err != nil {
		// 返回错误信息到前端
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
	})
}

// 登录逻辑
func login(ctx *gin.Context) {
	var user User
	contract := connect_fabric.GetContract()

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	// 调用 QueryUser 并处理返回值
	queriedUser, err := invoke_fabric.QueryUser(contract, user.Username, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 模拟生成一个令牌（实际应使用 JWT 或其他方式）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": User{
			Username:     queriedUser.Username,
			Organization: queriedUser.Organization,
			IsAdmin:      queriedUser.IsAdmin,
		},
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	})
	print(token)

	// 返回用户信息和令牌
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"username":     queriedUser.Username,
			"organization": queriedUser.Organization,
			"isadmin":      queriedUser.IsAdmin,
		},
	})
}
