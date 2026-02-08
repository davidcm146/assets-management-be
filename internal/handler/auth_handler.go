package handler

import (
	"fmt"
	"net/http"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/davidcm146/assets-management-be.git/internal/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepository
}

func NewAuthHandler(authService *service.AuthService, userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{authService: authService, userRepo: userRepo}
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.HandleValidationError(err, &req)
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusUnprocessableEntity,
			Code:       error_middleware.CodeUnprocessableEntity,
			Message:    "Dữ liệu không hợp lệ",
			Details:    errors,
		})
		return
	}

	token, err := h.authService.LoginService(c, req.Username, req.Password)
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
		errors := validator.HandleValidationError(err, &req)
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusUnprocessableEntity,
			Code:       error_middleware.CodeUnprocessableEntity,
			Message:    "Dữ liệu không hợp lệ",
			Details:    errors,
		})
		return
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}
	err := h.authService.RegisterService(c, user)
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

func (h *AuthHandler) GetMe(c *gin.Context) {
	userCtx, exists := c.Get("user")
	userID, ok := userCtx.(*dto.AuthUser).ID, true
	if !exists || !ok {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusUnauthorized,
			Code:       error_middleware.CodeUnauthorized,
			Message:    "Không tìm thấy thông tin người dùng",
		})
		return
	}

	user, err := h.userRepo.GetByID(c, userID)
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusNotFound,
			Code:       error_middleware.CodeNotFound,
			Message:    "Người dùng không tồn tại",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
