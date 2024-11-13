package utils

import (
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/sha256"
    "crypto/x509"
    "encoding/asn1"
    "encoding/base64"
    "encoding/pem"
    "errors"
	"math/big"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/ethereum/go-ethereum/crypto"
)

// 解析私钥 PEM 数据
func ParseECPrivateKeyFromPEM(pemBytes []byte) (*ecdsa.PrivateKey, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil || block.Type != "PRIVATE KEY" {
        return nil, errors.New("invalid private key PEM")
    }

    privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %v", err)
    }

    ecdsaPrivKey, ok := privKey.(*ecdsa.PrivateKey)
    if !ok {
        return nil, errors.New("parsed key is not an ECDSA private key")
    }

    return ecdsaPrivKey, nil
}

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

// 从签名和消息推断公钥(secp256k1)
func RecoverPublicKey(message string, signatureBase64 string) (*ecdsa.PublicKey, error) {
    // 计算消息的哈希值
    hash := sha256.Sum256([]byte(message))

    // 解码签名
    derBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
    if err != nil {
        return nil, fmt.Errorf("failed to decode signature: %v", err)
    }

    var sig struct {
        R, S *big.Int
    }
    _, err = asn1.Unmarshal(derBytes, &sig)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal signature: %v", err)
    }

    // 将 R 和 S 转换为字节切片
    rBytes := sig.R.Bytes()
    sBytes := sig.S.Bytes()

    // 确保 R 和 S 的长度为 32 字节
    if len(rBytes) < 32 {
        rBytes = append(make([]byte, 32-len(rBytes)), rBytes...)
    }
    if len(sBytes) < 32 {
        sBytes = append(make([]byte, 32-len(sBytes)), sBytes...)
    }

    // 合并 R 和 S 为 64 字节的签名
    sigBytes := append(rBytes, sBytes...)

    // 恢复公钥
    pubBytes, err := crypto.Ecrecover(hash[:], sigBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to recover public key: %v", err)
    }

    // 将恢复的公钥转换为 ecdsa.PublicKey
    pubKey, err := crypto.UnmarshalPubkey(pubBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %v", err)
    }

    return pubKey, nil
}


// 签名
func Sign(ctx contractapi.TransactionContextInterface, message string, privateKeyPEM []byte) (string, error) {
    // 解析私钥
    privKey, err := ParseECPrivateKeyFromPEM(privateKeyPEM)
    if err != nil {
        return "", err
    }

    // 计算消息的哈希值
    hash := sha256.Sum256([]byte(message))

    // 签名
    r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
    if err != nil {
        return "", fmt.Errorf("failed to sign message: %v", err)
    }

    // 将签名结果编码为 ASN.1 DER 格式
    derBytes, err := asn1.Marshal(asn1.RawValue{Tag: 0x30, Class: 0, IsCompound: true, Bytes: append(r.Bytes(), s.Bytes()...)})
    if err != nil {
        return "", fmt.Errorf("failed to marshal signature: %v", err)
    }

    return base64.StdEncoding.EncodeToString(derBytes), nil
}


// 验签
func Verify(ctx contractapi.TransactionContextInterface, message string, publicKeyPEM []byte, signatureBase64 string) (bool, error) {
    // 解析公钥
    pubKey, err := ParseECPublicKeyFromPEM(publicKeyPEM)
    if err != nil {
        return false, err
    }

    // 解码签名
    derBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
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

    // 计算消息的哈希值
    hash := sha256.Sum256([]byte(message))

    // 验签
    return ecdsa.Verify(pubKey, hash[:], sig.R, sig.S), nil
}
