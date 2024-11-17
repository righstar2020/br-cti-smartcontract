package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

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

//验证签名
func VerifySignature(message string, signatureBase64 string, publicKeyPEM []byte) (bool, error) {
    // 解析公钥
    pubKey, err := ParseECPublicKeyFromPEM(publicKeyPEM)
    if err != nil {
        return false, err
    }

    // 计算消息的哈希值
    hash := sha256.Sum256([]byte(message))

    // 解码签名
    r, s, err := DecodeSignature(signatureBase64)
    if err != nil {
        return false, err
    }
    // 验证签名
    return ecdsa.Verify(pubKey, hash[:], r, s), nil
}

// 解码签名
func DecodeSignature(signature string) (r, s *big.Int, err error) {
    // base64解码签名
    derBytes, err := base64.StdEncoding.DecodeString(signature)
    if err != nil {
        return nil, nil, err
    }

    _, err = asn1.Unmarshal(derBytes, &r)
    if err != nil {
        return nil, nil, err
    }

    _, err = asn1.Unmarshal(derBytes[len(r.Bytes()):], &s)
    if err != nil {
        return nil, nil, err
    }

    return r, s, nil
}

