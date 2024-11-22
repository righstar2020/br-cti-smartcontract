package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/sha256"
	"encoding/base64"
	"fmt"
	"bytes"
)

// 使用ECC公钥加密
func Encrypt(message string, publicKey *ecdsa.PublicKey) (string, error) {
    // 生成临时私钥
    ephemeralPrivateKey, err := ecdsa.GenerateKey(publicKey.Curve, rand.Reader)
    if err != nil {
        return "", fmt.Errorf("failed to generate ephemeral key: %v", err)
    }

    // 计算共享密钥
    sx, _ := publicKey.Curve.ScalarMult(publicKey.X, publicKey.Y, ephemeralPrivateKey.D.Bytes())
    sharedKey := sha256.Sum256(sx.Bytes())

    // 生成随机IV
    iv := make([]byte, aes.BlockSize)
    if _, err := rand.Read(iv); err != nil {
        return "", fmt.Errorf("failed to generate IV: %v", err)
    }

    // 创建AES密码块
    block, err := aes.NewCipher(sharedKey[:])
    if err != nil {
        return "", fmt.Errorf("failed to create AES cipher: %v", err)
    }

    // 加密消息
    paddedMessage := pkcs7Padding([]byte(message), aes.BlockSize)
    ciphertext := make([]byte, len(paddedMessage))
    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext, paddedMessage)

    // 组合最终的密文（临时公钥 + IV + 密文）
    finalMessage := append(ellipticPointToBytes(ephemeralPrivateKey.PublicKey), iv...)
    finalMessage = append(finalMessage, ciphertext...)

    return base64.StdEncoding.EncodeToString(finalMessage), nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padText...)
}

// 将椭圆曲线点转换为字节
func ellipticPointToBytes(pub ecdsa.PublicKey) []byte {
    return append(pub.X.Bytes(), pub.Y.Bytes()...)
}