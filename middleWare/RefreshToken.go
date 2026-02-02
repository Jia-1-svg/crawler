package middleWare

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	APP_KEY = "www.topgoer.com"
)

func TokenHandler(userId string, exp int64) (string, error) {

	// 颁发一个有限期一小时的证书
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    exp,
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
func ReFreShJwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			parseToken, err := GetToken(token)
			// 解析token有错误（如过期、签名错误）才走刷新逻辑
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 400,
					"msg":  "解析token失败",
				})
				c.Abort()
				return // 解析失败直接终止，避免后续无效操作
			}

			// 解析成功，才检查是否需要刷新
			expVal, expOk := parseToken["exp"].(float64)
			uidVal, uidOk := parseToken["userId"].(string)
			// 确保exp和userId存在且类型正确
			if expOk && uidOk {
				// 检查token是否即将过期（1小时内过期则刷新）
				if int64(expVal) < time.Now().Add(1*time.Hour).Unix() {
					newToken, err := TokenHandler(uidVal, time.Now().Add(2*time.Hour).Unix())
					if err == nil && newToken != "" {
						// 将新token放入响应头
						c.Header("token", newToken)
						fmt.Println("刷新后的token：", newToken)
					}
				}
			}
		}
		// 继续执行后续中间件/接口逻辑
		c.Next()
	}
}
