package encoding

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// SignJwt 生成jwt
func SignJwt(ctx context.Context, m map[string]any) (string, error) {
	signedString, err2 := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims(m)).SignedString(GetJwtPrivateKey())
	return signedString, err2
}

// VerifyJwt 校验jwt
func VerifyJwt(content string) (bool, error) {
	token, err := jwt.Parse(content, func(token *jwt.Token) (interface{}, error) {
		return GetJwtPublicKey(), nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

// UnpackJwt 解开jwt token内容
func UnpackJwt(content string) (map[string]any, error) {
	token, err := jwt.Parse(content, func(token *jwt.Token) (interface{}, error) {
		return GetJwtPublicKey(), nil
	})
	if err != nil {
		return nil, errors.New("fail to unpack jwt content")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("fail to unpack jwt")
}
