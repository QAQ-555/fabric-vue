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
	TaskID        string   `json:"ID"`
	Bonus         int      `json:"bonus"`
	RootModelId   string   `json:"rootModelHash"`
	PostedUser    string   `json:"postedUser"`
	AcceptedUsers []string `json:"acceptedUsers"`
	Models        []string `json:"models"`
	IsComplete    bool     `json:"isComplete"`
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

// 注册用户register
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

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

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
