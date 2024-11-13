package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "encoding/pem"
    "crypto/sha256"
    "encoding/asn1"
    "encoding/base64"
	"errors"
    "fmt"
    "os"
)

func main() {
    GenECCAccount()
}
func GenECCAccount(){
	// 生成椭圆曲线密钥对
    privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        fmt.Println("Error generating key pair:", err)
        return
    }

    // 将私钥保存到文件
    privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
    if err != nil {
        fmt.Println("Error marshaling private key:", err)
        return
    }
    privPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PRIVATE KEY",
        Bytes: privBytes,
    })
    err = os.WriteFile("private_key.pem", privPEM, 0600)
    if err != nil {
        fmt.Println("Error writing private key to file:", err)
        return
    }

    // 将公钥保存到文件
    pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
    if err != nil {
        fmt.Println("Error marshaling public key:", err)
        return
    }
    pubPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: pubBytes,
    })
    err = os.WriteFile("public_key.pem", pubPEM, 0600)
    if err != nil {
        fmt.Println("Error writing public key to file:", err)
        return
    }

    fmt.Println("Keys generated and saved successfully.")
}
// 签名
func Sign(message string, privateKeyPEM []byte) (string, error) {
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