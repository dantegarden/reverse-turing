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

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoytleH/hu+FNZ4WA0z/LGqvVK3jfe0uU1OKa6Nd6/lVT+dyi
kgLHxFKjbDLKJIMC7+iAteCUgzzqHmGJuKT7qja4qPchCJpNb8fAkR3OOHCSCXkP
Ar4jW4QarQevmCMYRpJPD+7XCY/v0lkxuxOyrMJUAA7zry/vlmUvI8kv2bhMus2X
3XdWpUL83rXOyFd9UbtptubDcKqs1Eq/gSWAYZy7g/BoWpjWaaRkxRNY5QS8GHpH
6h3jM5DLTlDFJrHYPmrdh+f2qOqwj9d8L19kz7i5YOo955sghnt+BT8k4Wuwg/SQ
dyUdnI1t52en7iYqR6VCUVp/O/nskufpOZRu6QIDAQABAoIBAFFf1x8dR8qXNi8m
mXTBH92RTKJ9iZbHvtXcnTz6GdC1ZUf7DOicklwKio3vVniXDePvpCEQe4Bn5Kp6
ImD/hrMAz18UwFi1+2B/0j2NC8eB/JLU2POZN7DwVQ1uA9hvyC+Jz/w2NPAD5KqW
6QJPdJBL6fCNhGIeGfJ7S+Mg9sgxIaCAWqz9OROvp87YmDrV8f+La2TLS87/vnZ3
e1uL30Ael2nIuhAkwgidLadmq2npbQ0mF6MkOXNhusfiJU27vR+nhUGgG8QWcBVR
oX5+6M/IOTWkwt4gVQmwRFxblxULf3ym//lgAonqYZxAm/6DuDks4IqYwRg8LOJf
djznoHECgYEA1gI8Wlnl5y4Z2uv4ENSms13reWDZtXEogBoXwDoiuQsib3Don3BH
kRhQLQSCWAucE7Qb+FF4g4Jd09Sa90/xdtqxIVPvPDyqUCUvPHxjGoJbczPzxwkp
1GvguJo3LnqVJg6RZDp1Mz+CB6Tj/BWqhcNL7MCnddVej2heeZBOAQ8CgYEAwy97
UMWE6x59ek+59VEpx69YK9ywB77QnUej7b1yCsK9qTdnX/OxyEyKpvBfuuJU1XTK
zTieohWOdzJNFr7WmkJCfB+gpdZZmKAH2M82AZA/GGMuGlv1RGF+ruL+yYAZ9aNs
aQvX5MWGH7746EnmpUMMR5CQW0JZB2UnNfqwIIcCgYEAnqLhogfph8iAmes44yD+
wQ7psfu85eaPowW0fWWav5glWn5TsXxFUKS5KeWhySox89kasqORtco5SwDaLmEw
GG5bxty1Be3iQa6OqUN7IvdmWqs0FWIRg8jDt5N5PBbZ4HAEDkvW/Loi5Q+xf/d/
g9AUw9a7S3lystMm1O3HJLkCgYEAnztoNb+9rPZl1VefVFOPaxlQLBBRBzCTDgx4
3qWTmNXAVoEdc3jii45t+rzUzCiCntU18XAEciR14iYGH802VAhEJvDCZShWVZ4Q
aL66x1G/N40J+nUUxWFoMRJ8WzSHeQ6Gjbgcu8Sso89vTkmjwSTOqr90FQ8uhErw
TyghcZcCgYBAhThu70sLAE29C1sUnQAQXVwXrdjtOpWyTg6yPbV7oKHiPZGBVscs
tCAaIBRCF9nSz1W/jgCaWhM0dk/fNgA6gKscxznXO5Fci79wSbvYV69J7u80Wu7B
53ruvh0k9mq0Uap6FsOFw+O2pObORxaFH6+Dic8xOSt+LCF3uuXQzw==
-----END RSA PRIVATE KEY-----`

func TestMW(t *testing.T) {
	encrypt, err := RsaEncrypt("silicon_hello", publicKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(encrypt)
}

func TestDecrpyt(t *testing.T) {
	s := "ZBrwoe2+le6RVpLvUb3pvchVLacRpa/wyvhU4yY/fR9RPA1VUw6E9XjQ+Z2Osm2/xfm4BbBHOgkXnmlmGsSh2z07Ag8bNXV+TzuL97t6w2Ow1Rk3GUVgqAY1dNApzbVD0YbOy4DlUHB9oGZLL4Wu/CR1dlAoP412TrFEk6S1rRW6VMA8uD0pK5Q4aSyCWbbyzkalIBATD2IHJJWZYfH/yP9z3OFdJarTyaBCzhu2jEAAuGxCgOTHSLfKh2gWPILFYw3vy279pe73/7iD7lYeM/vjtVtyC7mm81srDukG1bU97WlItu4Pp3Np4BqIo+NC2GzmkWlJkdgXWZt7pqEMxQ=="
	block, _ := pem.Decode([]byte(privateKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	encryptedData, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatal(err)
		return
	}
	decryptedData, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedData)
	fmt.Println(string(decryptedData))
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
