/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package connect_fabric

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type FabricConfig struct {
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
}

// 默认配置
var defaultConfig = FabricConfig{
	MSPID:        "org1MSP",
	CryptoPath:   "/tmp/hyperledger/org1/admin",
	CertPath:     "/tmp/hyperledger/org1/admin/msp/signcerts",
	KeyPath:      "/tmp/hyperledger/org1/admin/msp/keystore",
	TLSCertPath:  "/tmp/hyperledger/org1/peer1/tls-msp/tlscacerts/tls-0-0-0-0-7052.pem",
	PeerEndpoint: "dns:///localhost:7051",
	GatewayPeer:  "peer1-org1",
}

// 获取合约
func GetContract(config FabricConfig) *client.Contract {
	clientConnection := newGrpcConnection(config)

	id := newIdentity(config)
	sign := newSign(config)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	chaincodeName := "mycc"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	return contract
}

// 创建 gRPC 连接
func newGrpcConnection(config FabricConfig) *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(config.TLSCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, config.GatewayPeer)

	connection, err := grpc.NewClient(config.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// 创建身份
func newIdentity(config FabricConfig) *identity.X509Identity {
	certificatePEM, err := readFirstFile(config.CertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(config.MSPID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// 创建签名函数
func newSign(config FabricConfig) identity.Sign {
	privateKeyPEM, err := readFirstFile(config.KeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// 读取第一个文件
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// 格式化 JSON 数据
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

func queryUser(contract *client.Contract, username string) {
	fmt.Printf("\n--> Evaluate Transaction: ReadUser, 查询用户 %s\n", username)

	result, err := contract.EvaluateTransaction("ReadUser", username)
	if err != nil {
		panic(fmt.Errorf("查询用户失败: %w", err))
	}

	fmt.Printf("*** 用户详情: %s\n", formatJSON(result))
}
