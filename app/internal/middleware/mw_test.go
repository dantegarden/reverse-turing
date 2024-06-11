package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoytleH/hu+FNZ4WA0z/L
GqvVK3jfe0uU1OKa6Nd6/lVT+dyikgLHxFKjbDLKJIMC7+iAteCUgzzqHmGJuKT7
qja4qPchCJpNb8fAkR3OOHCSCXkPAr4jW4QarQevmCMYRpJPD+7XCY/v0lkxuxOy
rMJUAA7zry/vlmUvI8kv2bhMus2X3XdWpUL83rXOyFd9UbtptubDcKqs1Eq/gSWA
YZy7g/BoWpjWaaRkxRNY5QS8GHpH6h3jM5DLTlDFJrHYPmrdh+f2qOqwj9d8L19k
z7i5YOo955sghnt+BT8k4Wuwg/SQdyUdnI1t52en7iYqR6VCUVp/O/nskufpOZRu
6QIDAQAB
-----END PUBLIC KEY-----`

func TestMW(t *testing.T) {
	encrypt, err := RsaEncrypt("silicon_hello", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(encrypt)
}

var sha256Hash = sha256.New()

func RsaEncrypt(plainText, publicKey string) (string, error) {
	// 解码公钥PEM
	publicKeyBytes := []byte(publicKey)
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return "", errors.New("无效的公钥")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	// 对明文进行加密
	ciphertext, err := rsa.EncryptOAEP(sha256Hash, rand.Reader, pub, []byte(plainText), nil)
	if err != nil {
		return "", err
	}

	// 将加密结果编码为base64字符串
	encryptedData := base64.StdEncoding.EncodeToString(ciphertext)
	return encryptedData, nil
}
