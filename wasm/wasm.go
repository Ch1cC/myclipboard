// main.go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"syscall/js"
	"time"
)

var
// 加载公钥
publicKeyPEM = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA64gDlAFUlVNc/Fm1dN4k
njxrok2Y6C4mmnt0FDrC/jYO+pxzZ6mPkVS/JQuweHmrYVkQ6RJSXKew8I+2ukJc
Ny+N43ZuSPDqCHVECQlQkClTAug139cGBiMaUNnSWj2/d6R8DLXVYgfDuqPWBaCp
PJ9+jjy9WYGheoE/n5MPAhNSMqP4PDqt+auVJcWgVCrizeO/GuUn84Fm4J98Ln9s
9CqWcg/JSTGF1Za55FG9BfThW5bM0L+SpTKXzGco7jQ4QDF+bPFCzbzbUR638AbE
lHCBT+jGRhwwzWBJ8Z9bWY7NaYHJUv6OiQ+5J3OhcNkMa8rf8tIfCovNN3BPmJhl
9wIDAQAB
-----END PUBLIC KEY-----
`)

func encryptedFunc(this js.Value, args []js.Value) interface{} {

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		fmt.Println("Error decoding PEM block")
		return nil
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return nil
	}
	// 转换为 RSA 公钥对象
	rsaPubKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Failed to convert to RSA public key")
		return nil
	}

	// 获取当前的Unix时间戳
	unixTimestamp := time.Now().Unix()

	// 将Unix时间戳转换为字节数组
	data := make([]byte, 8) // 64位整数需要8个字节
	binary.BigEndian.PutUint64(data, uint64(unixTimestamp))
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, data, nil)
	if err != nil {
		fmt.Println("Failed to encrypt data:", err)
		return nil
	}
	encryptedData := hex.EncodeToString(ciphertext)
	return js.ValueOf(encryptedData)
}

func main() {
	done := make(chan int, 0)
	js.Global().Set("encryptedFunc", js.FuncOf(encryptedFunc))
	<-done
}
