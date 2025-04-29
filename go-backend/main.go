package main

import (
	invoke_fabric "backend/fabric-go/call"
	connect_fabric "backend/fabric-go/network"
	"encoding/json"
	"fmt"
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
	Accepted     []string `json:"accepted"`
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
	r.POST("/get_user_info", get_user_info)
	r.POST("/upload_public_key", upload_public_key)
	r.POST("/upload_model", upload_model)
	r.POST("/get_all_task", get_all_task)
	r.POST("/accept_task", accept_task)

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

// 用户信息查询逻辑
func get_user_info(ctx *gin.Context) {
	var user User
	contract := connect_fabric.GetContract()

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	// 打印接收到的用户信息
	fmt.Printf("查询用户信息: 用户名=%s, 组织=%s\n", user.Username, user.Organization)

	// 调用链码查询用户信息
	queriedUser, err := invoke_fabric.Get_one_User(contract, user.Username)
	if err != nil {
		// 返回错误信息到前端
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("查询用户失败: %s", err.Error())})
		return
	}

	// 返回用户信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": "用户信息查询成功",
		"user": gin.H{
			"username":     queriedUser.Username,
			"organization": queriedUser.Organization,
			"pubkeyhash":   queriedUser.Pubkeyhash,
			"token":        queriedUser.Token,
			"posted":       queriedUser.Posted,
			"accepted":     queriedUser.Accepted,
		},
	})
}

func upload_public_key(ctx *gin.Context) {
	var user User

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	contract := connect_fabric.GetContract()
	// 打印接收到的用户信息
	fmt.Printf("上传公钥: 用户名=%s, 公钥=%s\n", user.Username, user.Pubkeyhash)

	//调用链码上传公钥
	err := invoke_fabric.UploadPublicKey(contract, user.Username, user.Pubkeyhash)
	if err != nil {
		// 返回错误信息到前端
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传公钥失败: %s", err.Error())})
		return
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": "公钥上传成功",
	})
}

func upload_model(ctx *gin.Context) {
	var model struct {
		Username  string `json:"username"`
		Signature string `json:"signature"`
		CID       string `json:"cid"`
	}

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	// 打印接收到的 JSON 数据
	fmt.Printf("接收到的模型数据: 用户名=%s, 签名=%s, CID=%s\n", model.Username, model.Signature, model.CID)

	// 注释掉调用链码及其后续逻辑

	contract := connect_fabric.GetContract()

	// 调用链码上传模型
	err := invoke_fabric.CreateNewModel(contract, model.Username, model.CID, model.Signature)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("上传模型失败: %s", err.Error())})
		return
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": "模型上传成功",
	})

	// 返回调试信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": "接收到模型数据，已打印到终端供调试",
	})
}

func get_all_task(ctx *gin.Context) {
	contract := connect_fabric.GetContract()

	// 调用链码获取所有任务
	tasks, err := invoke_fabric.GetAllTasks(contract)
	if err != nil {
		// 返回错误信息到前端
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取任务失败: %s", err.Error())})
		return
	}

	// 返回任务信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": "任务获取成功",
		"tasks":   tasks,
	})
}

func accept_task(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
		TaskID   string `json:"taskID"`
	}

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	contract := connect_fabric.GetContract()

	// 检查用户是否已经接受了该任务
	user, err := invoke_fabric.Get_one_User(contract, request.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("查询用户信息失败: %s", err.Error())})
		return
	}

	for _, acceptedTask := range user.Accepted {
		if acceptedTask == request.TaskID {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("用户 %s 已经接受了任务 %s", request.Username, request.TaskID)})
			return
		}
	}

	// 调用链码将任务添加到用户的 Accepted 字段
	fmt.Printf("接受任务: 用户名=%s, 任务ID=%s\n", request.Username, request.TaskID)
	err = invoke_fabric.AddToAccepted(contract, request.Username, request.TaskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("添加任务到用户的 Accepted 字段失败: %s", err.Error())})
		return
	}

	// 调用链码将用户添加到任务的接受用户列表中
	err = invoke_fabric.AddUserToTask(contract, request.TaskID, request.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("将用户添加到任务的接受用户列表失败: %s", err.Error())})
		return
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("任务 %s 已成功被用户 %s 接受", request.TaskID, request.Username),
	})
}
