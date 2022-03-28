package rsa

import (
	_ "embed"
	"encoding/base64"
)
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

//go:embed rsa_public_key.pem
var pubkey []byte

//go:embed rsa_private_key.pem
var prikey []byte

func RSA_Encrypt(plainText []byte) []byte {

	//pem解码
	block, _ := pem.Decode(pubkey)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText
}

func RSA_Decrypt(cipherText []byte) []byte {

	block, _ := pem.Decode(prikey)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	//返回明文
	return plainText
}
func Decodestr(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		decrypt := RSA_Decrypt(decoded)
		return string(decrypt)
	}
	return ""
}
func DecodeByte(b []byte) []byte {
	l := len(b)
	SIZE := 128
	var result []byte
	var decrypt []byte
	for POINTER := 0; POINTER < l; POINTER += SIZE {
		if POINTER+SIZE <= l {
			decrypt = RSA_Decrypt(b[POINTER : POINTER+SIZE])
			result = append(result, decrypt...)
		} else {
			decrypt = RSA_Decrypt(b[POINTER:l])
			result = append(result, decrypt...)
		}
	}
	return result

}
