package router

import (
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/handler"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler *handler.AuthHandler
}

type RouterParams struct {
	Engine   *gin.Engine
	Handlers *Handlers
}

func NewRouter(params RouterParams) *gin.Engine {
	r := params.Engine

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(error_middleware.ErrorHandler())

	// Routes
	auth_api := r.Group("/api/auth")
	{
		auth_api.POST("/register", params.Handlers.AuthHandler.RegisterHandler)
		auth_api.POST("/login", params.Handlers.AuthHandler.LoginHandler)
	}

	// protected_api := r.Group("/api")
	// {
	// 	// Add protected routes here
	// }

	return r
}
