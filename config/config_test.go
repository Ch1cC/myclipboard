// math_test.go

package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) { // JavaScript 加密后的数据（Base64 编码的字符串）
	encryptedDataHex := "26115993bec522f9ae913a14719feb448d9b42b0aeaa078f98dbac3ed5bb4f802842312cba719150ad5bbc2a1adb67c90997779337d268004603ca88b62fad7ca84faf0c1c85a6c15d6fdfce55f3cfbbfe2e442bc83008f6263aee46d99f44f75cf8e739b5ee4efeed138824b316588d210848e2e1ce1b0be97d9499134eb96203eb9493db97bcc5ae5939c84c305c40c05e7d52adcd646de5543318b47c38a1a4d96bdb8188e839a164fd6c1fde9eda27dee9d922dd6e1c2c6a3575f5071a12e46ed52cb6d38bf8fd0212c407fd713f79d05a032f00dec0d111505ab82b9ccf1f41b54c4264aa59d167d63182aae5449cb995f12ec2d4c024d6dee1d8011686" // 将此处替换为实际的加密后的数据

	// 将十六进制格式的字符串解码为字节数组
	encryptedData, err := hex.DecodeString(encryptedDataHex)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	// 加载私钥
	privateKeyPEM := []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA64gDlAFUlVNc/Fm1dN4knjxrok2Y6C4mmnt0FDrC/jYO+pxz
Z6mPkVS/JQuweHmrYVkQ6RJSXKew8I+2ukJcNy+N43ZuSPDqCHVECQlQkClTAug1
39cGBiMaUNnSWj2/d6R8DLXVYgfDuqPWBaCpPJ9+jjy9WYGheoE/n5MPAhNSMqP4
PDqt+auVJcWgVCrizeO/GuUn84Fm4J98Ln9s9CqWcg/JSTGF1Za55FG9BfThW5bM
0L+SpTKXzGco7jQ4QDF+bPFCzbzbUR638AbElHCBT+jGRhwwzWBJ8Z9bWY7NaYHJ
Uv6OiQ+5J3OhcNkMa8rf8tIfCovNN3BPmJhl9wIDAQABAoIBAF60ojuxDUG983XE
3KHRcQfvY5gve3xgkkUrvtEjj6SIOB2tKvpTw9a8LN/Ig3zt72uFVLdjFpsLEqqo
HuFmLY5DINqFlHLf7DrabDD1/d75vtc8Z+1JmLJ/CuXXsC15BrQk/3xc4SA53wn9
NuGsaGBoAYtQARMsfLG1ps+TYkWCmeEjx+uvwTdf5t4pXnu32TnGL8pZgJFM39dX
ScPKgSkQO8fZiqhu9+K5tH+sG6fi+sP1Tw2dj+Xwq0ewboaMhQ8DZf8i59jYo3QW
g8doqSY1ZBZGwnPhuoPxtnnvhU16/GUlAG5h+PjURSLltvLHoEz07ffjbK0KiUm/
N21jg4ECgYEA9Jlaw3ohPyjNZD05+O5rjjzfpbRt5TNjXji4tYqZhxX6GskdJkIn
6CxylogC0E7+52TFmrZ0h0wbg75HZAwnDIwgxWNADYqhTv5bsv/kDx1/urs4Hnss
XXDp4LL6dn6UFelPKXEO/gl8TQcMiGGcwS0IdjdzzfhNxLwC/EDEC/UCgYEA9oJ1
ukgCTaOgMJ3Hh1QfL2bmKbczinTBseH27isVkW66xbk3eEOQN/2GxMP6Gl+ewIqi
YXtRIWK4t80xvTLs4hItHls1aiaokBXj5Rgd4o0uidkhkj1Ls+iLmrcMHBBKW8sa
emPa+ccaNBnE4mtr+3aCbGe7ONGqxolVYJaEwrsCgYAkbTI3KlkJLupnozae++LI
rAgihVxYZe7GeWwInTuAAqXcl1bf+7o1uWjXQiopG5qam0dSYxm3jH4MgKnhHG40
UCoRO1aurZaYTQka/0DXf20mQft5jp5szAQIkp76Rp+HI9fGNDAnZQI99m7HYMIX
gr1f3aJBalkqk1Vee8a2gQKBgBikC86WGhzWqVGSw/okD4X2fDVZSb6iUyZL1xoR
lnNWJTdUf1X3MvhV0F3k1SBDxKOsd/TUldSHgL1mtn0aFRG4DWiGZ3135cuZVJF2
6q3VvPwshy2OEP4n1aSefYhknHo2gCwRxTbIjzb8CHE1mbKmMv1RFSbl1nNIEQ+5
nAW3AoGAKQZ0u1nd8YAUH1T2rHR+21z7BUneKnhXqZ21oYynxDAfpgdhHYSQFv9P
5fKjqW2pnVjO4Uyk3H3YD6g1sBH5uMb9D6czrRcp2T9iKYuRefzSdn6bM3w8vzvw
C5zQCQYrICnLj6xQ3qla8aX2fLqy6pObVLG3ajZ0att+Rdfip8U=
-----END RSA PRIVATE KEY-----
	`)
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		fmt.Println("Error decoding PEM block")
		return
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return
	}

	// 解密数据
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedData, nil)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return
	}

	// 打印解密后的数据
	fmt.Println("Decrypted data:", decryptedData)
}

func TestGen(t *testing.T) {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA private key:", err)
		return
	}

	// 将私钥转换为 PKCS#1 DER 编码
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// 创建私钥的 PEM 块
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 将私钥 PEM 块进行编码
	privateKeyPEMEncoded := pem.EncodeToMemory(privateKeyPEM)

	// 将私钥打印出来
	fmt.Println("Private Key:")
	fmt.Println(string(privateKeyPEMEncoded))

	// 从私钥中提取公钥
	publicKey := &privateKey.PublicKey

	// 将公钥进行 DER 编码
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error marshaling public key:", err)
		return
	}

	// 创建公钥的 PEM 块
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 将公钥 PEM 块进行编码
	publicKeyPEMEncoded := pem.EncodeToMemory(publicKeyPEM)

	// 将公钥打印出来
	fmt.Println("Public Key:")
	fmt.Println(string(publicKeyPEMEncoded))
}
