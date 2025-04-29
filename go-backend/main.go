package main

import (
	invoke_fabric "backend/fabric-go/call"
	connect_fabric "backend/fabric-go/network"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	users = make(map[string]string) // 内存中的用户数据
	mu    sync.Mutex
)

// 默认配置
var defaultConfig = connect_fabric.FabricConfig{
	MSPID:        "org1MSP",
	CryptoPath:   "/tmp/hyperledger/org1/admin",
	CertPath:     "/tmp/hyperledger/org1/admin/msp/signcerts",
	KeyPath:      "/tmp/hyperledger/org1/admin/msp/keystore",
	TLSCertPath:  "/tmp/hyperledger/org1/peer1/tls-msp/tlscacerts/tls-0-0-0-0-7052.pem",
	PeerEndpoint: "dns:///localhost:7051",
	GatewayPeer:  "peer1-org1",
}

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
	TaskID          string   `json:"ID"`
	Bonus           int      `json:"bonus"`
	RootModelId     string   `json:"rootModelHash"`
	PostedUser      string   `json:"postedUser"`
	AcceptedUsers   []string `json:"acceptedUsers"`
	Models          []string `json:"models"`
	IsComplete      bool     `json:"isComplete"`
	Round           int      `json:"round"`           // 新增：任务的轮次
	NextRoundTaskID string   `json:"nextRoundTaskID"` // 新增：下一轮任务的 ID
}

func main() {
	r := gin.Default()

	// 配置跨域
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

	// 路由配置
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/get_user_info", get_user_info)
	r.POST("/upload_public_key", upload_public_key)
	r.POST("/upload_model", upload_model)
	r.POST("/get_all_task", get_all_task)
	r.POST("/accept_task", accept_task)
	r.POST("/get_all_users", get_all_users)
	r.POST("/delete_user", delete_user)
	r.POST("/verify_user", verify_user)
	r.POST("/new_task", new_task)
	r.POST("/next_task_round", next_task_round)
	r.POST("/delete_task", delete_task)
	r.POST("/finish_task", finish_task)
	r.Run(":8089")
}

// 注册逻辑
func register(ctx *gin.Context) {
	var user User
	contract := connect_fabric.GetContract(defaultConfig)

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}
	print(user.Username)
	print(user.Password)
	print(user.Organization)
	// 调用 QueryUser 并处理返回值
	err := invoke_fabric.CreateNewUser(contract, user.Username, user.Password, user.Organization, "test", 0, false, false, false)
	if err != nil {
		// 返回错误信息到前端
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// 登录逻辑
func login(ctx *gin.Context) {
	var user User
	contract := connect_fabric.GetContract(defaultConfig)

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

	// 检查 IsAccepted 和 IsVerified 状态
	if !queriedUser.IsAccepted || !queriedUser.IsVerified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未被接受或未验证"})
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
	contract := connect_fabric.GetContract(defaultConfig)

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

// 上传公钥
func upload_public_key(ctx *gin.Context) {
	var user User

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	contract := connect_fabric.GetContract(defaultConfig)
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

// 上传模型
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

	contract := connect_fabric.GetContract(defaultConfig)

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

// 获取所有任务
func get_all_task(ctx *gin.Context) {
	contract := connect_fabric.GetContract(defaultConfig)

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

// 接受任务
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

	contract := connect_fabric.GetContract(defaultConfig)

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

// 获取所有用户信息
func get_all_users(ctx *gin.Context) {
	contract := connect_fabric.GetContract(defaultConfig)

	// 调用链码获取所有用户
	result, err := contract.EvaluateTransaction("GetAllUsers")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取用户失败: %s", err.Error())})
		return
	}

	// 如果返回值为空，返回提示信息
	if len(result) == 0 || string(result) == "null" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "没有找到任何用户",
			"users":   []string{},
		})
		return
	}

	// 返回用户信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": "用户获取成功",
		"users":   json.RawMessage(result),
	})
}

func delete_user(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
	}

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	contract := connect_fabric.GetContract(defaultConfig)

	// 调用链码删除用户
	err := invoke_fabric.DeleteUser(contract, request.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除用户失败: %s", err.Error())})
		return
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("用户 %s 已成功删除", request.Username),
	})
}

func verify_user(ctx *gin.Context) {
	var request struct {
		Username   string `json:"username"`
		IsAdmin    bool   `json:"isAdmin"`
		IsAccepted bool   `json:"isAccepted"`
	}
	// 绑定 JSON 数据
	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("绑定 JSON 失败: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}
	print("test2")
	print(request.Username)
	print(request.IsAdmin)
	print(request.IsAccepted)
	contract := connect_fabric.GetContract(defaultConfig)

	// 调用链码更新用户的 isAdmin 和 isAccepted 状态
	err := invoke_fabric.ManageUser(contract, request.Username, request.IsAdmin, true, request.IsAccepted)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("更新用户状态失败: %s", err.Error())})
		return
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("用户 %s 的状态已成功更新", request.Username),
	})
}

func new_task(c *gin.Context) {
	var requestBody struct {
		Username    string `json:"username"`
		Bonus       int    `json:"bonus"`
		RootModelId string `json:"rootModelId"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取 Fabric 合约实例
	contract := connect_fabric.GetContract(defaultConfig)

	// 调用 createNewTask 函数
	round := 1            // 初始轮数为 1
	nextRoundTaskID := "" // 初始任务没有下一轮任务 ID
	invoke_fabric.CreateNewTask(contract, requestBody.Bonus, requestBody.RootModelId, requestBody.Username, round, nextRoundTaskID)

	c.JSON(http.StatusOK, gin.H{"message": "任务创建成功"})
}

func next_task_round(c *gin.Context) {
	var requestBody struct {
		TaskID      string `json:"taskId"`
		RootModelId string `json:"rootModelId"`
		Username    string `json:"username"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		fmt.Printf("绑定 JSON 失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取 Fabric 合约实例
	contract := connect_fabric.GetContract(defaultConfig)

	// 调用 next_round 函数
	err := invoke_fabric.Next_round(contract, requestBody.TaskID, requestBody.RootModelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("任务轮次更新失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务轮次更新成功"})
}

func delete_task(ctx *gin.Context) {
	var request struct {
		TaskID string `json:"taskId"`
	}

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	contract := connect_fabric.GetContract(defaultConfig)

	// 调用链码删除任务
	err := invoke_fabric.DeleteTask(contract, request.TaskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除任务失败: %s", err.Error())})
		return
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("任务 %s 已成功删除", request.TaskID),
	})
}

func finish_task(ctx *gin.Context) {
	var request struct {
		TaskID string `json:"taskId"`
	}

	// 绑定 JSON 数据
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的 JSON 数据"})
		return
	}

	contract := connect_fabric.GetContract(defaultConfig)
	get_task, err := invoke_fabric.QueryTask(contract, request.TaskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取任务数据失败: %s", err.Error())})
		return
	}
	// 调用链码读取任务信息
	acceptedUsers, err := invoke_fabric.Finish_Task(contract, request.TaskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("修改任务时失败;%s", err.Error())})
		return
	}
	// 提取接受任务的用户名单

	if len(acceptedUsers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"success": "没有用户接受该任务"})
		return
	}

	// 向接受任务的用户转账
	for _, user := range acceptedUsers {
		err := invoke_fabric.TransferTokens(contract, "task_owner", user, get_task.Bonus)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("向用户 %s 转账失败: %s", user, err.Error())})
			return
		}
	}

	// 返回成功信息到前端
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("任务 %s 已成功完成，奖励已发放:", request.TaskID),
	})
}
