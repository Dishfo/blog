package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"time"
)

/**
提供令牌支持
*/
const (
	issuer             = "Dishfo"
	tokenZet           = "tokens"
	tokenCleanSchedule = "schedule"
)

var (
	hmacSampleSercret = []byte("dishfo.159357ghj")
)

func init() {
	randStr := time.Now().String()
	hmacSampleSercret = []byte(randStr)
}

func GenerateToken(admin string) string {
	issueAt := time.Now().UnixNano()
	expiresAt := issueAt + int64(time.Hour*12)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		Audience:  admin,
		IssuedAt:  issueAt,
		ExpiresAt: expiresAt,
		Id:        "",
		NotBefore: issueAt,
		Subject:   admin,
	})

	tokenString, _ := token.SignedString(hmacSampleSercret)
	cacheToken(tokenString, expiresAt)
	return tokenString
}

func cacheToken(token string, expiresAt int64) {
	conn := client.Get()
	_, _ = conn.Do("ZADD", tokenZet, token)
}

//移除超时的token
func cleanExpiredToken() {
	conn := client.Get()
	_, _ = conn.Do("ZREMRANGEBYSCORE", tokenZet, 0, time.Now().UnixNano())
}

func CheckToken(token string) bool {
	conn := client.Get()
	r, _ := conn.Do("ZSCORE", tokenZet, token)
	if r != nil {
		expire, _ := redis.Int64(r, nil)
		if expire < time.Now().UnixNano() {
			return false
		}
		return true
	}
	return false
}
