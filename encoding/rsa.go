package encoding

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
)

func RsaEncrypt(plain string) (string, error) {
	// 使用OAEP填充方式加密
	hashFunc := crypto.SHA256.New
	sign, err := rsa.EncryptOAEP(hashFunc(), rand.Reader, GetRsaPublicKey(), []byte(plain), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func RsaDecrypt(plain string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(plain)
	hashFunc := crypto.SHA256.New
	bytes, err := rsa.DecryptOAEP(hashFunc(), rand.Reader, GetRsaPrivateKey(), decodeBytes, nil)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
