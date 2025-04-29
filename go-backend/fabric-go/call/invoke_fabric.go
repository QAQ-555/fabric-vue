/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package invoke_fabric

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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

// 初始化用户账本
func InitUserLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, 初始化用户数据 \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("初始化账本失败: %w", err))
	}

	fmt.Printf("*** 用户数据初始化成功\n")
}

// 查询用户
func QueryUser(contract *client.Contract, username string, password string) (*User, error) {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询用户 %s\n", username)

	// 调用链码查询用户信息
	result, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 将查询结果解析为 User 结构体
	var user User
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 验证密码
	if user.Password != password {
		return nil, fmt.Errorf("密码错误")
	}

	// 检查 IsVerified 和 IsAccepted 字段
	if !user.IsVerified || !user.IsAccepted {
		return nil, fmt.Errorf("用户未通过验证或未被接受")
	}

	// 如果所有条件都满足，返回用户信息
	return &user, nil
}

// 注册用户
func CreateNewUser(contract *client.Contract, username, password, org, pubkeyhash string, token int, isAdmin, isVerified, isAccepted bool) error {
	fmt.Printf("\n--> Submit Transaction: CreateUser, 创建新用户 %s\n", username)

	_, err := contract.SubmitTransaction("CreateUser",
		username, password, org, pubkeyhash,
		fmt.Sprintf("%d", token),
		fmt.Sprintf("%t", isAdmin),
		fmt.Sprintf("%t", isVerified),
		fmt.Sprintf("%t", isAccepted))
	if err != nil {
		return fmt.Errorf("注册失败")
	}
	return nil
}

// 格式化 JSON 数据
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

// 查询单个用户
func Get_one_User(contract *client.Contract, username string) (*User, error) {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询用户 %s\n", username)

	// 调用链码查询用户信息
	result, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 将查询结果解析为 User 结构体
	var user User
	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 如果所有条件都满足，返回用户信息
	return &user, nil
}

// 上传公钥
func UploadPublicKey(contract *client.Contract, username string, pubkeyhash string) error {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询用户 %s\n", username)

	// 调用链码查询用户信息
	result, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 将查询结果解析为 User 结构体
	var user User
	err = json.Unmarshal(result, &user)
	if err != nil {
		return fmt.Errorf("解析用户信息失败: %w", err)
	}

	// 更新用户的公钥哈希
	user.Pubkeyhash = pubkeyhash
	print(pubkeyhash)
	print("test")
	print(user.Pubkeyhash)
	fmt.Printf("\n--> Submit Transaction: UpdateUser, 更新用户 %s 的公钥哈希\n", username)

	// 调用链码更新用户信息
	_, err = contract.SubmitTransaction(
		"UpdateUser",
		user.Username,
		user.Password,
		user.Organization,
		user.Pubkeyhash,
		fmt.Sprintf("%d", user.Token),
		fmt.Sprintf("%t", user.IsAdmin),
		fmt.Sprintf("%t", user.IsVerified),
		fmt.Sprintf("%t", user.IsAccepted),
	)
	if err != nil {
		return fmt.Errorf("更新用户公钥哈希失败: %w", err)
	}

	fmt.Printf("*** 用户 %s 的公钥哈希已成功更新\n", username)
	return nil
}

// 上传模型
func CreateNewModel(contract *client.Contract, modelowner, modelhash, modelsign string) error {
	fmt.Printf("\n--> Submit Transaction: CreateModel, 创建新模型 %s\n", modelhash)

	// 调用链码的 CreateModel 方法
	result, err := contract.SubmitTransaction("CreateModel", modelowner, modelhash, modelsign)
	if err != nil {
		return fmt.Errorf("创建模型失败: %w", err)
	}

	modelID := string(result)
	fmt.Printf("*** 模型创建成功, 模型ID: %s\n", modelID)

	// 使用 add_2_posed 函数将模型 ID 添加到用户的 Posted 列表
	fmt.Printf("\n--> Submit Transaction: AddToPosted, 更新用户 %s 的 Posted 列表\n", modelowner)
	add_2_posed(contract, modelowner, modelID)

	fmt.Printf("*** 用户 %s 的 Posted 列表已成功更新\n", modelowner)
	return nil
}

// 添加到 Posted 列表
func add_2_posed(contract *client.Contract, username, modelid string) {
	result, err := contract.SubmitTransaction("AddToPosted", username, modelid)
	if err != nil {
		fmt.Printf("Failed to add to Posted: %v\n", err)
		return
	}
	fmt.Printf("Successfully added to Posted: %s\n", string(result))
	if err != nil {
		fmt.Printf("Failed to add to Posted: %v\n", err)
	}
}

func GetAllTasks(contract *client.Contract) ([]map[string]interface{}, error) {
	fmt.Println("\n--> Evaluate Transaction: GetAllTasks, 查询所有任务")

	result, err := contract.EvaluateTransaction("GetAllTasks")
	if err != nil {
		return nil, fmt.Errorf("查询所有任务失败: %w", err)
	}

	// 如果返回值为空，直接返回提示信息
	if len(result) == 0 || string(result) == "null" {
		return nil, fmt.Errorf("没有找到任何任务")
	}

	// 将结果解析为 JSON 对象
	var tasks []map[string]interface{}
	err = json.Unmarshal(result, &tasks)
	if err != nil {
		return nil, fmt.Errorf("解析任务 JSON 失败: %w", err)
	}

	fmt.Printf("*** 所有任务: %s\n", formatJSON(result))
	return tasks, nil
}

// 添加任务到用户的 Accepted 字段
func AddToAccepted(contract *client.Contract, username, taskID string) error {
	fmt.Printf("\n--> Submit Transaction: AddToAccepted, 将任务 %s 添加到用户 %s 的 Accepted 字段\n", taskID, username)

	// 调用链码的 AddToAccepted 方法
	_, err := contract.SubmitTransaction("AddToAccepted", username, taskID)
	if err != nil {
		return fmt.Errorf("添加任务到 Accepted 字段失败: %w", err)
	}

	fmt.Printf("*** 成功将任务 %s 添加到用户 %s 的 Accepted 字段\n", taskID, username)
	return nil
}

// 将用户添加到任务的接受用户列表中
func AddUserToTask(contract *client.Contract, taskID, username string) error {
	fmt.Printf("\n--> Submit Transaction: AddUserToTask, 将用户 %s 添加到任务 %s 的接受用户列表中\n", username, taskID)

	// 调用链码的 AddUserToTask 方法
	_, err := contract.SubmitTransaction("AddUserToTask", taskID, username)
	if err != nil {
		return fmt.Errorf("将用户添加到任务的接受用户列表失败: %w", err)
	}

	fmt.Printf("*** 成功将用户 %s 添加到任务 %s 的接受用户列表中\n", username, taskID)
	return nil
}

// 删除用户
func DeleteUser(contract *client.Contract, username string) error {
	fmt.Printf("\n--> Submit Transaction: DeleteUser, 删除用户 %s\n", username)

	// 调用链码的 DeleteUser 方法
	_, err := contract.SubmitTransaction("DeleteUser", username)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	fmt.Printf("*** 用户 %s 已成功删除\n", username)
	return nil
}

// 管理用户
func ManageUser(contract *client.Contract, username string, isAdmin, isVerified, isAccepted bool) error {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询用户 %s\n", username)

	// 调用链码查询用户信息
	result, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 将查询结果解析为 User 结构体
	var user User
	err = json.Unmarshal(result, &user)
	if err != nil {
		return fmt.Errorf("解析用户信息失败: %w", err)
	}

	// 更新用户的状态
	user.IsAdmin = isAdmin
	user.IsVerified = isVerified
	user.IsAccepted = isAccepted

	fmt.Printf("\n--> Submit Transaction: UpdateUser, 更新用户 %s 的状态\n", username)

	// 调用链码更新用户信息
	_, err = contract.SubmitTransaction(
		"UpdateUser",
		user.Username,
		user.Password,
		user.Organization,
		user.Pubkeyhash,
		fmt.Sprintf("%d", user.Token),
		fmt.Sprintf("%t", user.IsAdmin),
		fmt.Sprintf("%t", user.IsVerified),
		fmt.Sprintf("%t", user.IsAccepted),
	)
	if err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	fmt.Printf("*** 用户 %s 的状态已成功更新\n", username)
	return nil
}

// 创建新任务
func CreateNewTask(contract *client.Contract, bonus int, rootModelId, postedUser string, round int, nextRoundTaskID string) string {
	fmt.Printf("\n--> Submit Transaction: CreateTask, 创建新任务\n")

	// 调用链码的 CreateTask 方法
	result, err := contract.SubmitTransaction("CreateTask",
		"", // taskID 由链码生成
		fmt.Sprintf("%d", bonus),
		rootModelId,
		postedUser,
		fmt.Sprintf("%d", round),
		nextRoundTaskID)
	if err != nil {
		panic(fmt.Errorf("创建任务失败: %w", err))
	}

	// 解析返回的 taskID
	taskID := string(result)
	fmt.Printf("*** 任务创建成功, 任务ID: %s\n", taskID)
	return taskID
}

func DeleteTask(contract *client.Contract, Taskid string) error {
	fmt.Printf("\n--> Submit Transaction: DeleteUser, 删除任务 %s\n", Taskid)

	_, err := contract.SubmitTransaction("DeleteTask", Taskid)

	if err != nil {
		return fmt.Errorf("删除任务失败: %w", err)
	}

	fmt.Printf("*** 任务 %s 已成功删除\n", Taskid)
	return nil
}

func Next_round(contract *client.Contract, taskID string, rootModelID string) error {
	// 调用链码读取任务信息
	result, err := contract.EvaluateTransaction("ReadTask", taskID)
	if err != nil {
		return fmt.Errorf("读取任务失败: %v", err)
	}
	var task Task
	err = json.Unmarshal(result, &task)
	if err != nil {
		return fmt.Errorf("解析任务信息失败: %w", err)
	}
	print(task.TaskID + "\n")
	print(task.RootModelId + "\n")
	print(task.PostedUser + "\n")
	// 更新任务信息
	newRound := task.Round + 1

	// 创建新任务
	nexttaskid := CreateNewTask(contract,
		task.Bonus,
		rootModelID,
		task.PostedUser,
		newRound,
		"")
	print(nexttaskid + "\n")
	/*
		UpdateTask(ctx contractapi.TransactionContextInterface,
			taskID string,
			bonus int,
			rootModelId string,
			postedUser string,
			isComplete bool,
			round int,
			nextRoundTaskID string
		)
	*/
	// 调用链码更新原任务
	_, err = contract.SubmitTransaction("UpdateTask",
		task.TaskID,
		fmt.Sprintf("%d", task.Bonus),
		task.RootModelId,
		task.PostedUser,
		fmt.Sprintf("%t", true),
		fmt.Sprintf("%d", task.Round),
		nexttaskid,
	)

	if err != nil {
		return fmt.Errorf("更新原任务失败: %v", err)
	}

	return nil
}

func TransferTokens(contract *client.Contract, sender, receiver string, amount int) error {
	// 查询接收者信息
	receiverResult, err := contract.EvaluateTransaction("ReadUser", receiver)
	if err != nil {
		return fmt.Errorf("查询接收者失败: %w", err)
	}

	// 解析接收者数据
	var receiverData User
	if err := json.Unmarshal(receiverResult, &receiverData); err != nil {
		return fmt.Errorf("解析接收者数据失败: %w", err)
	}

	// 增加接收者余额
	receiverData.Token += amount

	// 更新接收者信息
	_, err = contract.SubmitTransaction("UpdateUser",
		receiverData.Username,
		receiverData.Password,
		receiverData.Organization,
		receiverData.Pubkeyhash,
		fmt.Sprintf("%d", receiverData.Token),
		fmt.Sprintf("%t", receiverData.IsAdmin),
		fmt.Sprintf("%t", receiverData.IsVerified),
		fmt.Sprintf("%t", receiverData.IsAccepted),
	)
	if err != nil {
		return fmt.Errorf("更新接收者信息失败: %w", err)
	}

	// 返回成功信息
	return nil
}

func Finish_Task(contract *client.Contract, taskID string) ([]string, error) {
	// 调用链码读取任务信息
	result, err := contract.EvaluateTransaction("ReadTask", taskID)
	if err != nil {
		return nil, fmt.Errorf("读取任务失败: %v", err)
	}
	var task Task
	err = json.Unmarshal(result, &task)
	if err != nil {
		return nil, fmt.Errorf("解析任务信息失败: %w", err)
	}
	/*
		UpdateTask(ctx contractapi.TransactionContextInterface,
			taskID string,
			bonus int,
			rootModelId string,
			postedUser string,
			isComplete bool,
			round int,
			nextRoundTaskID string
		)
	*/
	// 调用链码更新原任务
	_, err = contract.SubmitTransaction("UpdateTask",
		task.TaskID,
		fmt.Sprintf("%d", task.Bonus),
		task.RootModelId,
		task.PostedUser,
		fmt.Sprintf("%t", true),
		fmt.Sprintf("%d", task.Round),
		task.NextRoundTaskID,
	)
	if err != nil {
		return nil, fmt.Errorf("更新原任务失败: %v", err)
	}

	return task.AcceptedUsers, nil
}

func QueryTask(contract *client.Contract, taskID string) (*Task, error) {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询任务 %s\n", taskID)

	// 调用链码查询用户信息
	result, err := contract.EvaluateTransaction("ReadTask", taskID)
	if err != nil {
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	// 将查询结果解析为 User 结构体
	var task Task
	err = json.Unmarshal([]byte(result), &task)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	// 如果所有条件都满足，返回用户信息
	return &task, nil
}
