package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

func init() {
	//jwtKey = []byte(os.Getenv("JWT_SECRET"))
	jwtKey = []byte("qingbo1011.top-gin-game-3")
}

type Claims struct {
	Uid int
	jwt.StandardClaims
}

// GenerateToken 生成Token
func GenerateToken(uid int) (string, error) {
	// 过期时间 默认7天
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析token
func ParseToken(tokenStr string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return token, claims, err
}
