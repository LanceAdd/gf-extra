package encoding

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type RsaKeys struct {
	sync.RWMutex
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type RsaConfig struct {
	sync.RWMutex
	PrivateKey string `json:"privateKey" v:"required" dc:"私钥"`
	PublicKey  string `json:"publicKey" v:"required" dc:"公钥"`
}

func doLoadRsaConfig(ctx context.Context, config *RsaConfig) {
	rsaConfigMap, _ := g.Cfg().Get(ctx, "rsa")
	err := gconv.Scan(rsaConfigMap, config)
	if err != nil {
		panic(err)
	}
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		panic(err)
	}
}

func doLoadRsaKeys(config *RsaConfig, keys *RsaKeys) {
	publicKeyDecodedBytes, err := base64.StdEncoding.DecodeString(config.PublicKey)
	if err != nil {
		panic(err)
	}
	privateKeyDecodedBytes, err := base64.StdEncoding.DecodeString(config.PrivateKey)
	if err != nil {
		panic(err)
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyDecodedBytes)
	if err != nil {
		panic(err)
	}
	privateKeyInterface, err := x509.ParsePKCS1PrivateKey(privateKeyDecodedBytes)
	if err != nil {
		panic(err)
	}
	keys.PublicKey = publicKeyInterface.(*rsa.PublicKey)
	keys.PrivateKey = privateKeyInterface
}

func GetRsaPrivateKey() *rsa.PrivateKey {
	rsaKeys.RLock()
	defer rsaKeys.RUnlock()
	return rsaKeys.PrivateKey
}

func GetRsaPublicKey() *rsa.PublicKey {
	rsaKeys.RLock()
	defer rsaKeys.RUnlock()
	return rsaKeys.PublicKey
}

func GetRsaPrivateKeyContent() string {
	rsaCfg.RLock()
	defer rsaCfg.RUnlock()
	return rsaCfg.PrivateKey
}

func GetRsaPublicKeyContent() string {
	rsaCfg.RLock()
	defer rsaCfg.RUnlock()
	return rsaCfg.PublicKey
}

func doRsaCfgCompare(source *RsaConfig, target *RsaConfig) bool {
	if source.PublicKey != target.PublicKey {
		return false
	}

	if source.PrivateKey != target.PrivateKey {
		return false

	}
	return true
}
