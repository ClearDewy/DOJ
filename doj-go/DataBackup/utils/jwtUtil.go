/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Token struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

const (
	secretKey = "零露漙兮~"
	// ExpireTime token过期时间
	ExpireTime = 24 * time.Hour
)

func JwtGenerate(uid string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Token{
		uid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpireTime)), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                 // 生效时间
		},
	}).SignedString([]byte(secretKey))
}

func JwtParse(token string) (*Token, error) {
	t, err := jwt.ParseWithClaims(token, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*Token); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
