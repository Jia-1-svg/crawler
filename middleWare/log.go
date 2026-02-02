package middleWare

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogData struct {
	Timestamp  string `json:"timestamp"`   // 日志时间
	TokenExist bool   `json:"token_exist"` // token是否存在
	UserID     string `json:"user_id"`     // 解析成功的用户ID
	IsError    bool   `json:"is_error"`    // 是否出错
}

func JWTLog(logData LogData) {
	// 序列化为JSON字符串（忽略错误，避免日志逻辑影响主流程）
	logJSON, err := json.Marshal(logData)
	if err != nil {
		fmt.Printf("[JWT_LOG_ERROR] 日志序列化失败: %s\n", err.Error())
		return
	}
	// 输出日志（可替换为写入日志文件/日志组件，如zap/logrus）
	fmt.Printf("[JWT_MIDDLEWARE] %s\n", string(logJSON))
}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ========== 初始化埋点日志数据 ==========
		startTime := time.Now()

		logData := LogData{
			Timestamp:  startTime.Format("2006-01-02 15:04:05"),
			TokenExist: false,
			IsError:    false,
		}
		//startTime := time.Now()
		//
		//logData := LogData{
		//    Timestamp:  startTime.Format("2006-01-02 15:04:05.000"),
		//    TokenExist: false,
		//    IsError:    false,
		//}

		// 延迟执行日志输出（即使c.Abort()也能执行，保证埋点不丢失）
		defer func() {
			// 输出结构化日志
			JWTLog(logData)
		}()

		// ========== 原有核心逻辑 ==========
		token := c.Request.Header.Get("token")
		if token == "" {
			// 埋点：token不存在
			logData.TokenExist = false
			logData.IsError = true
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "token不能为空",
			})
			c.Abort()
			return
		}

		// 标记token存在
		logData.TokenExist = true

		// 解析token
		parseToken, err := GetToken(token)
		if err != nil {
			// 埋点：解析token失败
			logData.IsError = true
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "解析token失败",
			})
			c.Abort()
			return
		}

		// 埋点：解析成功，记录用户ID
		userId := parseToken["userId"].(string)
		logData.UserID = userId
		c.Set("userId", userId)

		// 继续执行后续中间件/路由
		c.Next()
	}
}

//		//获取黑名单key
//		limitUser, _ := config.Rdb.Get(context.Background(), "list:user").Result()
//		res := config.Rdb.SIsMember(context.Background(),  limitUser, uid).Val()
//		if res {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code": 400,
//				"msg":  "被加入黑名单",
//			})
//			c.Abort()
//			return
//		}
