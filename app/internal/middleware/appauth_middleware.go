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
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"
	"net/http"
	"reverse-turing/common/consts"
	"strings"
)

type AppAuthMiddleware struct {
	PrivateKey string
}

func NewAppAuthMiddleware(privateKey string) *AppAuthMiddleware {
	return &AppAuthMiddleware{PrivateKey: privateKey}
}

func (m *AppAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w, r, errors.New("Token不能为空"))
			return
		}

		tokenString := strings.Split(authHeader, consts.TokenType)
		if len(tokenString) < 2 {
			unauthorized(w, r, errors.New("Bearer token格式错误"))
			return
		}

		apiKey := tokenString[1]
		token, err := rsaDecryptWithPEM(apiKey, m.PrivateKey)
		if err != nil {
			unauthorized(w, r, errors.New("无效的token"))
			return
		}

		if !strings.HasPrefix(token, consts.AuthPrefix) {
			unauthorized(w, r, errors.New("无效的token, 请检查token是否正确"))
			return
		}
		// 截掉前缀，放到ctx里
		token = token[len(consts.AuthPrefix):]
		ctx := r.Context()
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(
			"device_id", token,
		))

		next(w, r.WithContext(ctx))
	}
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
		"code": 41401,
		"msg":  err.Error(),
	})
}

// rsaDecryptWithPEM 使用PEM格式的私钥进行解密
func rsaDecryptWithPEM(encrypted string, privateKeyPem string) (string, error) {
	// 解码私钥
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 将base64编码的字符串解码回字节切片
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	// 进行解密
	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		encryptedBytes,
		nil,
	)
	if err != nil {
		return "", err
	}

	// 返回解密后的字符串
	return string(decrypted), nil
}
