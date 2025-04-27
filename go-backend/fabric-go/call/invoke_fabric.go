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
func QueryUser(contract *client.Contract, username string) {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询用户 %s\n", username)

	result, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		panic(fmt.Errorf("查询用户失败: %w", err))
	}

	fmt.Printf("*** 用户详情: %s\n", formatJSON(result))
}

func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
