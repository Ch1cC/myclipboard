package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

var (
	Duration   time.Duration
	privateKey *rsa.PrivateKey
	signature  []byte
	Hash       []byte
	Token      *big.Int
)

//	func ConfigRsa() {
//		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
//		if err != nil {
//			panic(err)
//		}
//		// 将publicKey转换为PKIX, ASN.1 DER格式
//		if derPkix, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey); err != nil {
//			fmt.Println(err)
//			fmt.Printf("转换失败\n")
//		} else {
//			// 设置PEM编码结构
//			block := pem.Block{
//				Type:  "RSA PUBLIC KEY",
//				Bytes: derPkix,
//			}
//			// 将publicKey以字符串形式返回给javascript
//			PublicKeyString = string(pem.EncodeToMemory(&block))
//		}
//	}
func ConfigRandom() {
	random, _ := rand.Int(rand.Reader, big.NewInt(1000))
	Token = random
}

//	func ConfigHash() string {
//		random, _ := rand.Int(rand.Reader, big.NewInt(1000))
//		// 对消息进行哈希
//		hash := sha256.Sum256(random.Bytes())
//		// hash = hash[:]
//		localSignature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
//		if err != nil {
//			fmt.Println("Error signing:", err)
//			return ""
//		}
//		signature = localSignature
//		hashStr := hex.EncodeToString(hash[:])
//		return hashStr
//	}
func init() {
	// configECC()
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

	localPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return
	}
	privateKey = localPrivateKey
}
func Verify(encryptedDataHex string) bool {
	// 将十六进制格式的字符串解码为字节数组
	encryptedData, err := hex.DecodeString(encryptedDataHex)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return false
	}
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedData, nil)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return false
	}
	return decryptedData != nil
}

// func configECC() {
// 	// 选择椭圆曲线
// 	curve := elliptic.P256()
// 	// 生成私钥
// 	localPrivateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	if err != nil {
// 		fmt.Println("Error generating private key:", err)
// 		return
// 	}
// 	privateKey = localPrivateKey
// }
