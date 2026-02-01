package middleWare

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	APP_KEY = "www.topgoer.com"
)

func TokenHandler(userId string) (string, error) {

	// 颁发一个有限期一小时的证书
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":    time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	return tokenString, err
}

func GetToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(APP_KEY), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, err
	} else {
		fmt.Println(err)
	}

	return nil, nil
}

// RefreshToken 刷新JWT令牌
func RefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(APP_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//fmt.Println(claims["foo"], claims["nbf"])
		//payload, err := ParseToken(tokenString)
		// 提取原payload

		payload, ok := claims["payload"].(string)
		if !ok {
			return "", errors.New("invalid payload format")
		}

		// 生成新令牌
		return TokenHandler(payload)
		//return claims["payload"].(string), err
	} else {
		fmt.Println(err)
	}
	return "", errors.New("invalid token structure")
}
