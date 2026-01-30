package handler

import (
	"fmt"
	"net/http"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       error_middleware.CodeBadRequest,
			Message:    "Yêu cầu không hợp lệ",
		})
		return
	}

	token, err := h.authService.Login(c, req.Username, req.Password)
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusUnauthorized,
			Code:       error_middleware.CodeUnauthorized,
			Message:    "Tên người dùng hoặc mật khẩu không đúng",
		})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       error_middleware.CodeBadRequest,
			Message:    "Yêu cầu không hợp lệ",
		})
		return
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}
	err := h.authService.Register(c, user)
	fmt.Println("Registering user:", err)
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       error_middleware.CodeInternal,
			Message:    "Đăng ký không thành công",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tạo tài khoản thành công"})
}
