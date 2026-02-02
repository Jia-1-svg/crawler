package middleWare

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "token不能为空",
				"code": 500,
			})
			c.Abort()
			return
		}
		getToken, err := GetToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "jwt解密失败",
				"code": 500,
			})
			c.Abort()
			return
		}
		//res := config.Rdb.SIsMember(context.Background(), "list:k1", getToken["userId"]).Val()
		//if res {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"msg":  "用户黑名单",
		//		"code": 500,
		//	})
		//	c.Abort()
		//	return
		//}
		c.Set("userId", getToken["userId"].(string))
		c.Next()
	}
}

//handler, err := middleWare.TokenHandler(strconv.FormatInt(login.UserId, 10), time.Now().Add(time.Hour*time.Duration(1)).Unix())
