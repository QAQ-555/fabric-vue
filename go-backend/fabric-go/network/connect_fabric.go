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
	"os/exec"
	"path"
	"path/filepath"
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

func GenerateCertificateAndUpdateConfig(username, password, org string, config *FabricConfig) error {
	// 设置 CA 地址，根据 org 动态选择端口
	caURL := fmt.Sprintf("0.0.0.0:%d", map[string]int{"org1": 7054, "org2": 7055}[org])

	// 设置环境变量
	fabricCAClientHome := fmt.Sprintf("./%s/%s/", org, username)
	fmt.Sprintf("./%s/%s/", org, username)
	tlsCertPath := fmt.Sprintf("/tmp/hyperledger/%s/peer1/assets/ca/%s-ca-cert.pem", org, org)
	configtlsCertPath := fmt.Sprintf("/tmp/hyperledger/%s/peer1/tls-msp/tlscacerts/tls-0-0-0-0-7052.pem", org)
	fmt.Printf("设置环境变量:\n")
	fmt.Printf("FABRIC_CA_CLIENT_HOME=%s\n", fabricCAClientHome)
	fmt.Printf("FABRIC_CA_CLIENT_TLS_CERTFILES=%s\n", tlsCertPath)
	fmt.Printf("FABRIC_CA_CLIENT_MSPDIR=msp\n")
	os.Setenv("FABRIC_CA_CLIENT_TLS_CERTFILES", tlsCertPath)
	os.Setenv("FABRIC_CA_CLIENT_HOME", fabricCAClientHome)
	os.Setenv("FABRIC_CA_CLIENT_MSPDIR", "msp")
	// 创建目录
	fmt.Printf("创建目录: %s\n", fabricCAClientHome)
	if err := os.MkdirAll(fabricCAClientHome, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 构建 enroll 命令
	enrollCmdExec := exec.Command("fabric-ca-client", "enroll", "-d",
		"-u", fmt.Sprintf("https://%s:%s@%s", username, password, caURL))
	enrollCmdExec.Stdout = os.Stdout
	enrollCmdExec.Stderr = os.Stderr

	// 执行 enroll 命令
	if err := enrollCmdExec.Run(); err != nil {
		return fmt.Errorf("执行 fabric-ca-client enroll 失败: %w", err)
	}

	// 动态设置 MSPID、PeerEndpoint 和 GatewayPeer
	config.MSPID = fmt.Sprintf("%sMSP", org)
	config.PeerEndpoint = fmt.Sprintf("dns:///localhost:%d", map[string]int{"org1": 7051, "org2": 9051}[org])
	config.GatewayPeer = fmt.Sprintf("peer1-%s", org)

	// 更新 FabricConfig
	config.CryptoPath = fabricCAClientHome
	absPath, err := filepath.Abs(fabricCAClientHome)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %w", err)
	}
	config.CertPath = filepath.Join(absPath, "msp", "signcerts")
	config.KeyPath = filepath.Join(absPath, "msp", "keystore")
	config.TLSCertPath = configtlsCertPath
	fmt.Printf("证书生成成功，存储在: %s\n", fabricCAClientHome)
	return nil
}

func RevokeCertificateAndCleanup(username, org string) error {
	// 设置 CA 地址，根据 org 动态选择端口
	caURL := fmt.Sprintf("0.0.0.0:%d", map[string]int{"org1": 7054, "org2": 7055}[org])

	// 设置环境变量
	tlsCertPath := fmt.Sprintf("/tmp/hyperledger/%s/peer1/assets/ca/%s-ca-cert.pem", org, org)
	fmt.Printf("设置环境变量:\n")
	fmt.Printf("FABRIC_CA_CLIENT_TLS_CERTFILES=%s\n", tlsCertPath)
	fabricCAClientHome := fmt.Sprintf("./%s/%s/", org, username)
	os.Setenv("FABRIC_CA_CLIENT_TLS_CERTFILES", tlsCertPath)
	os.Setenv("FABRIC_CA_CLIENT_HOME", fabricCAClientHome)

	// 获取用户证书路径
	userCertPath := fmt.Sprintf("./%s/%s/msp/signcerts/cert.pem", org, username)
	if _, err := os.Stat(userCertPath); os.IsNotExist(err) {
		return fmt.Errorf("用户证书不存在: %s", userCertPath)
	}

	// 使用 openssl 获取证书的序列号和 AKI
	fmt.Printf("获取证书的序列号和 AKI...\n")
	serialCmd := exec.Command("openssl", "x509", "-in", userCertPath, "-serial", "-noout")
	serialOut, err := serialCmd.Output()
	if err != nil {
		return fmt.Errorf("获取证书序列号失败: %w", err)
	}
	serial := string(bytes.TrimSpace(bytes.Split(serialOut, []byte("="))[1]))

	akiCmd := exec.Command("openssl", "x509", "-in", userCertPath, "-text")
	akiOut, err := akiCmd.Output()
	if err != nil {
		return fmt.Errorf("获取证书 AKI 失败: %w", err)
	}
	aki := ""
	for _, line := range bytes.Split(akiOut, []byte("\n")) {
		if bytes.Contains(line, []byte("keyid")) {
			aki = string(bytes.ToLower(bytes.ReplaceAll(bytes.TrimSpace(bytes.Split(line, []byte(":"))[1]), []byte(":"), []byte(""))))
			break
		}
	}
	if aki == "" {
		return fmt.Errorf("未能解析证书 AKI")
	}

	// 构建 revoke 命令
	fmt.Printf("撤销证书: Serial=%s, AKI=%s\n", serial, aki)
	revokeCmdExec := exec.Command("fabric-ca-client", "revoke", "-d",
		"-u", fmt.Sprintf("https://%s", caURL),
		"-s", serial,
		"-a", aki,
		"-r", "affiliationchange",
		"--tls.certfiles", tlsCertPath)
	revokeCmdExec.Stdout = os.Stdout
	revokeCmdExec.Stderr = os.Stderr

	// 执行 revoke 命令
	if err := revokeCmdExec.Run(); err != nil {
		return fmt.Errorf("执行 fabric-ca-client revoke 失败: %w", err)
	}

	// 清理临时目录
	if fabricCAClientHome != "" {
		if err := os.RemoveAll(fabricCAClientHome); err != nil {
			return fmt.Errorf("清理临时目录失败: %w", err)
		}
	}

	return nil
}

func RegisterIdentity(username, secret, idType, org string) error {
	// 设置 CA 地址，根据 org 动态选择端口
	caURL := fmt.Sprintf("0.0.0.0:%d", map[string]int{"org1": 7054, "org2": 7055}[org])

	// 设置环境变量
	tlsCertPath := fmt.Sprintf("/tmp/hyperledger/%s/peer1/assets/ca/%s-ca-cert.pem", org, org)
	fabricCAClientHome := fmt.Sprintf("/tmp/hyperledger/%s/ca/admin", org)
	fmt.Printf("设置环境变量:\n")
	fmt.Printf("FABRIC_CA_CLIENT_TLS_CERTFILES=%s\n", tlsCertPath)
	fmt.Printf("FABRIC_CA_CLIENT_HOME=%s\n", fabricCAClientHome)
	os.Setenv("FABRIC_CA_CLIENT_TLS_CERTFILES", tlsCertPath)
	os.Setenv("FABRIC_CA_CLIENT_HOME", fabricCAClientHome)

	// 构建 register 命令
	registerCmd := fmt.Sprintf("fabric-ca-client register -d --id.name %s --id.secret %s --id.type %s --id.attrs \"hf.Revoker=true\" -u https://%s", username, secret, idType, caURL)
	fmt.Printf("执行命令: %s\n", registerCmd)
	registerCmdExec := exec.Command("fabric-ca-client", "register", "-d",
		"--id.name", username,
		"--id.secret", secret,
		"--id.type", idType,
		"--id.attrs", "hf.Revoker=true",
		"-u", fmt.Sprintf("https://%s", caURL))
	registerCmdExec.Stdout = os.Stdout
	registerCmdExec.Stderr = os.Stderr

	// 执行 register 命令
	if err := registerCmdExec.Run(); err != nil {
		return fmt.Errorf("执行 fabric-ca-client register 失败: %w", err)
	}

	fmt.Printf("注册成功: 用户名=%s, 密码=%s, 类型=%s\n", username, secret, idType)
	return nil
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
