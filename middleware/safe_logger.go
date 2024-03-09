package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goflet/util/log"
)

// SafeLogger returns a safe logger middleware
func SafeLogger() gin.HandlerFunc {
	fun := gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// strip the token from the header
		param.Request.Header.Del("Authorization")

		// strip the token from the query
		query := param.Request.URL.Query()
		query.Del("token")
		param.Request.URL.RawQuery = query.Encode()

		if param.Request.URL.RawQuery != "" {
			param.Request.URL.RawQuery = "?" + param.Request.URL.RawQuery
		}

		param.Path = param.Request.URL.Path + param.Request.URL.RawQuery

		str := fmt.Sprintf("%s GIN: %3d | %13v | %15s | %-7s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)

		log.RawPrintf(str)
		return ""
	})
	return fun
}
