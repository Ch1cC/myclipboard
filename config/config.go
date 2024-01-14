package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

var (
	Duration time.Duration
	//PublicKey       rsa.PublicKey
	PrivateKey      *rsa.PrivateKey
	Token           *big.Int
	PublicKeyString string
)

func ConfigRsa() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	// 将publicKey转换为PKIX, ASN.1 DER格式
	if derPkix, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey); err != nil {
		fmt.Println(err)
		fmt.Printf("转换失败\n")
	} else {
		// 设置PEM编码结构
		block := pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: derPkix,
		}
		// 将publicKey以字符串形式返回给javascript
		PublicKeyString = string(pem.EncodeToMemory(&block))
	}
}
func ConfigRandom() {
	random, _ := rand.Int(rand.Reader, big.NewInt(1000))
	Token = random
}
