// math_test.go

package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) { // JavaScript 加密后的数据（Base64 编码的字符串）
	encryptedDataHex := "ba2ea204cb15d423ad5e0708288b91b9170c9e278e4ec17a50e5f4795d9626570fea2996e3e994f5128a25057629fabcbeb838db884d81b18d7f6dbb289b15c8633e65cff0b4055049c525e82b4ab55a4e9f00a2aff04529bac2646e79dc3cdb94761cef912306152d930d3129bcd0909f24b099398f18358442c18596f928c21a31a3f93027a09f86a622bac63bbc967c7247c01f46333705dd7ea78c78bd403b4d675cfe12b6b13a25a905a7c74bf9012a4ead8a580ba70316febd656be3c59d74b3369340c6af4309619be7b578114889a59b91de571bec9986cad4838e976f578223cd6c9589a9cbaec082c210515634d7fcc0c1ba6ef6e8a931c998d42f" // 将此处替换为实际的加密后的数据

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
	// 为了验证转换是否正确，可以将字节切片转换回int64
	convertedBack := int64(binary.BigEndian.Uint64(decryptedData))
	// 打印解密后的数据
	fmt.Println("Decrypted data:", convertedBack)
}

func TestGenRsa(t *testing.T) {
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

func TestGenEcc(t *testing.T) {
	// 选择椭圆曲线
	curve := elliptic.P256()

	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}
	derBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Println("Failed to marshal EC private key:", err)
		return
	}

	privateKeyPEMEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: derBytes})
	// 将私钥打印出来
	fmt.Println("Private Key:")
	fmt.Println(string(privateKeyPEMEncoded))
	// 获取公钥
	publicKey := &privateKey.PublicKey

	derBytes, err = x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Failed to marshal PKIX public key:", err)
	}

	publicKeyPEMEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: derBytes})
	// 将公钥打印出来
	fmt.Println("Public Key:")
	fmt.Println(string(publicKeyPEMEncoded))
	// 原始数据
	message := []byte("Hello, world!")

	// 对消息进行哈希
	hash := sha256.Sum256(message)

	// 使用私钥对哈希值进行签名
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		fmt.Println("Error signing:", err)
		return
	}

	// 使用公钥验证签名
	valid := ecdsa.VerifyASN1(publicKey, hash[:], signature)
	fmt.Println("Signature verified:", valid)
}

func TestLoadEcc(t *testing.T) {
	privateKeyPEM := []byte(`
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIDZVje58fmaJGi6n4zXnE28OLEG+0ue5CU763XyTCWCmoAoGCCqGSM49
AwEHoUQDQgAE67BSDTBuDdu957wRaaKp/v/1Hm7dpIKqz7JPET+2USl2/Pc76OtR
jsITCrMMUhuB1e9sxS2ElSAYYYVc1I1CYA==
-----END EC PRIVATE KEY-----
	`)
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		fmt.Println("no PEM data found in file")
		return
	}

	pk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		fmt.Printf("failed to parse EC private key: %s", err)
		return
	}
	fmt.Println(pk)
}
