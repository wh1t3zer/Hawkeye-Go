package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/utils"
)

// RecoveryMiddleware 捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//先做一下日志记录
				fmt.Println(string(debug.Stack()))
				utils.ComLogNotice(c, "_com_panic", map[string]interface{}{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				})

				if lib.ConfBase.DebugMode != "debug" {
					ResponseError(c, 500, errors.New("内部错误"))
					return
				}
				ResponseError(c, 500, errors.New(fmt.Sprint(err)))
			}
		}()
		c.Next()
	}
}
