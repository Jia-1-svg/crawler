package middleWare

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NegotiationCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentVersion := "v1.0.0"
		etag := fmt.Sprintf(`W/"%x"`, md5.Sum([]byte(contentVersion)))

		// 2. 检查客户端携带的 ETag（核心校验）
		clientETag := c.GetHeader("If-None-Match")

		// 3. 缓存命中则返回 304（核心逻辑）
		if clientETag == etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}

		// 4. 设置响应头（让浏览器缓存 ETag）
		c.Header("ETAG", etag)
		c.Header("Cache-Control", "no-cache") // 可选保留，强制走协商缓存

		c.Next()
	}
}
