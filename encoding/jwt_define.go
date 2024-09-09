package encoding

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

type JwtKeys struct {
	sync.RWMutex
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

type JwtConfig struct {
	sync.RWMutex
	PrivateKeyContent string `json:"privateKey" v:"required" dc:"私钥"`
	PublicKeyContent  string `json:"publicKey" v:"required" dc:"公钥"`
}

func doLoadJwtConfig(ctx context.Context, config *JwtConfig) {
	jwtConfigMap, _ := g.Cfg().Get(ctx, "jwt")
	err := gconv.Scan(jwtConfigMap, config)
	if err != nil {
		panic(err)
	}
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		panic(err)
	}
}

func doLoadJwt(config *JwtConfig, keys *JwtKeys) {
	privateKey, publicKey := doLoadJwtKeys(config)
	keys.PrivateKey = privateKey
	keys.PublicKey = publicKey
}

func doLoadJwtKeys(config *JwtConfig) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(config.PrivateKeyContent)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(config.PublicKeyContent)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}
	return privateKey, publicKey
}
func doJwtConfigCompare(source *JwtConfig, target *JwtConfig) bool {
	if source.PublicKeyContent != target.PublicKeyContent {
		return false
	}

	if source.PrivateKeyContent != target.PrivateKeyContent {
		return false

	}
	return true
}
func GetJwtPrivateKey() *ecdsa.PrivateKey {
	jwtKeys.RLock()
	defer jwtKeys.RUnlock()
	return jwtKeys.PrivateKey
}

func GetJwtPublicKey() *ecdsa.PublicKey {
	jwtKeys.RLock()
	defer jwtKeys.RUnlock()
	return jwtKeys.PublicKey
}

func GetJwtPrivateKeyContent() string {
	jwtCfg.RLock()
	defer jwtCfg.RUnlock()
	return jwtCfg.PrivateKeyContent
}

func GetJwtPublicKeyContent() string {
	jwtCfg.RLock()
	defer jwtCfg.RUnlock()
	return jwtCfg.PublicKeyContent
}
