package encoding

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Encrypt(plain string) string {
	data := []byte(plain)
	hash := sha256.New()
	hash.Write(data)
	// 计算哈希值
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
