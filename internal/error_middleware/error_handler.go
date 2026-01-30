package error_middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := err.(*AppError); ok {
			c.JSON(appErr.HTTPStatus, gin.H{
				"http_status": appErr.HTTPStatus,
				"code":        appErr.Code,
				"message":     appErr.Message,
			})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"http_status": http.StatusInternalServerError,
				"code":        CodeInternal,
				"message":     "Lỗi máy chủ",
			})
		}
		c.Abort()
	}
}
