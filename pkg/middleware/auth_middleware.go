package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUser struct {
	ID   int
	Role model.Role
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Error(&error_middleware.AppError{
				HTTPStatus: http.StatusUnauthorized,
				Code:       error_middleware.CodeUnauthorized,
				Message:    "Yêu cầu xác thực",
			})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == "" {
			c.Error(&error_middleware.AppError{
				HTTPStatus: http.StatusUnauthorized,
				Code:       error_middleware.CodeUnauthorized,
				Message:    "Yêu cầu xác thực",
			})
			return
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.Error(&error_middleware.AppError{
				HTTPStatus: http.StatusUnauthorized,
				Code:       error_middleware.CodeUnauthorized,
				Message:    "Token không hợp lệ",
			})
			return
		}
		user := &AuthUser{
			ID:   int(token.Claims.(jwt.MapClaims)["sub"].(int)),
			Role: model.Role(int(token.Claims.(jwt.MapClaims)["role"].(int))),
		}
		c.Set("user", user)
		c.Set("role", user.Role)

		c.Next()
	}
}
