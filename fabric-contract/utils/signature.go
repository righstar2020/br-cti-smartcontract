package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


// 解析公钥 PEM 数据
func ParseECPublicKeyFromPEM(pemBytes []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key PEM")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not an ECDSA public key")
	}

	return ecdsaPubKey, nil
}

// 验签
func Verify(ctx contractapi.TransactionContextInterface, txData []byte, publicKeyPEM []byte, signature []byte) (bool, error) {
	// 解析公钥
	pubKey, err := ParseECPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return false, err
	}

	// 解码签名
	derBytes, err := base64.StdEncoding.DecodeString(string(signature)) //utf-8解码	
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	var sig struct {
		R, S *big.Int
	}
	_, err = asn1.Unmarshal(derBytes, &sig)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal signature: %v", err)
	}

	// 计算消息的哈希值(链下也需要保证msg使用的是sha256)
	hash := sha256.Sum256(txData)

	// 验签
	return ecdsa.Verify(pubKey, hash[:], sig.R, sig.S), nil
}

// 获取和更新CTIID计数器
func GenerateNextCTIID(ctx contractapi.TransactionContextInterface) (string, error) {
	// 获取当前CTIID计数器的值
	ctiIDCounterKey := "CTIID_COUNTER"
	currentCountAsBytes, err := ctx.GetStub().GetState(ctiIDCounterKey)
	if err != nil {
		return "", fmt.Errorf("failed to get CTIID counter: %v", err)
	}

	var currentCount int
	if currentCountAsBytes == nil {
		// 如果计数器不存在，则初始化为0
		currentCount = 0
	} else {
		// 将存储的计数器值转换为整数
		err = json.Unmarshal(currentCountAsBytes, &currentCount)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal current count: %v", err)
		}
	}

	// 生成新的CTIID（通过递增计数器）
	newCTIID := fmt.Sprintf("CTI_%d", currentCount+1)

	// 更新计数器
	currentCount++
	newCountAsBytes, err := json.Marshal(currentCount)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated count: %v", err)
	}
	err = ctx.GetStub().PutState(ctiIDCounterKey, newCountAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update CTIID counter: %v", err)
	}

	// 返回生成的CTIID
	return string(newCTIID), nil
}


// 获取和更新ModelID计数器
func GenerateNextModelID(ctx contractapi.TransactionContextInterface) (string, error) {
	// 获取当前ModelID计数器的值
	modelIDCounterKey := "MODELID_COUNTER"
	currentCountAsBytes, err := ctx.GetStub().GetState(modelIDCounterKey)
	if err != nil {
		return "", fmt.Errorf("failed to get ModelID counter: %v", err)
	}

	var currentCount int
	if currentCountAsBytes == nil {
		// 如果计数器不存在，则初始化为0
		currentCount = 0
	} else {
		// 将存储的计数器值转换为整数
		err = json.Unmarshal(currentCountAsBytes, &currentCount)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal current count: %v", err)
		}
	}

	// 生成新的ModelID（通过递增计数器）
	newModelID := fmt.Sprintf("MODEL_%d", currentCount+1)
	// 更新计数器
	currentCount++
	newCountAsBytes, err := json.Marshal(currentCount)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated count: %v", err)
	}
	err = ctx.GetStub().PutState(modelIDCounterKey, newCountAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to update ModelID counter: %v", err)
	}

	// 返回生成的ModelID
	return newModelID, nil
}
