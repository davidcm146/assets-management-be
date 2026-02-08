package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/policy"
	"github.com/davidcm146/assets-management-be.git/internal/service"
	"github.com/davidcm146/assets-management-be.git/internal/validator"
	"github.com/gin-gonic/gin"
)

type LoanSlipHandler struct {
	loanSlipService *service.LoanSlipService
	uploader        service.Uploader
	policy          policy.LoanSlipPolicy
}

func NewLoanSlipHandler(loanSlipService *service.LoanSlipService, uploader service.Uploader) *LoanSlipHandler {
	return &LoanSlipHandler{
		loanSlipService: loanSlipService,
		uploader:        uploader,
	}
}

func (h *LoanSlipHandler) LoanSlipsListHandler(c *gin.Context) {
	var query dto.LoanSlipQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       error_middleware.CodeBadRequest,
			Message:    "Yêu cầu không hợp lệ",
		})
		return
	}

	result, err := h.loanSlipService.LoanSlipsListService(c.Request.Context(), &query)
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       error_middleware.CodeInternal,
			Message:    "Lỗi hệ thống",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *LoanSlipHandler) CreateLoanSlipHandler(c *gin.Context) {
	var req dto.CreateLoanSlipRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusUnprocessableEntity,
			Code:       error_middleware.CodeValidationFailed,
			Message:    err.Error(),
			Details:    validator.HandleValidationError(err, &req),
		})
		return
	}

	user := c.MustGet("user").(*dto.AuthUser)

	loan, err := h.loanSlipService.CreateLoanSlipService(
		c.Request.Context(),
		user.ID,
		&req,
	)
	if err != nil || loan == nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       error_middleware.CodeInternal,
			Message:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

func extractUpdateFields(updateDTO *dto.UpdateLoanSlipRequest) []string {
	fields := []string{}

	if updateDTO.Name != nil {
		fields = append(fields, "name")
	}
	if updateDTO.BorrowerName != nil {
		fields = append(fields, "borrower_name")
	}
	if updateDTO.Department != nil {
		fields = append(fields, "department")
	}
	if updateDTO.Position != nil {
		fields = append(fields, "position")
	}
	if updateDTO.Description != nil {
		fields = append(fields, "description")
	}
	if updateDTO.Status != nil {
		fields = append(fields, "status")
	}
	if updateDTO.SerialNumber != nil {
		fields = append(fields, "serial_number")
	}
	if updateDTO.BorrowedDate != nil {
		fields = append(fields, "borrowed_date")
	}
	if updateDTO.ReturnedDate != nil {
		fields = append(fields, "returned_date")
	}
	return fields
}

func (h *LoanSlipHandler) UpdateLoanSlipHandler(c *gin.Context) {
	var req dto.UpdateLoanSlipRequest

	if err := c.ShouldBind(&req); err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusUnprocessableEntity,
			Code:       error_middleware.CodeValidationFailed,
			Message:    err.Error(),
			Details:    validator.HandleValidationError(err, &req),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       error_middleware.CodeBadRequest,
			Message:    "ID không hợp lệ",
		})
		return
	}

	user := c.MustGet("user").(*dto.AuthUser)
	fields := extractUpdateFields(&req)

	forbidden := h.policy.ForbiddenFields(user.Role, fields)
	fmt.Println(forbidden)
	if len(forbidden) > 0 {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusForbidden,
			Code:       error_middleware.CodeForbidden,
			Message:    "Bạn không có quyền cập nhật một số trường",
			Details: gin.H{
				"fields": forbidden,
			},
		})
		return
	}

	loan_slip, err := h.loanSlipService.UpdateLoanSlipService(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       error_middleware.CodeInternal,
			Message:    "Lỗi hệ thống",
		})
		return
	}

	c.JSON(http.StatusOK, loan_slip)
}

func (h *LoanSlipHandler) LoanSlipDetailHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       error_middleware.CodeBadRequest,
			Message:    "ID không hợp lệ",
		})
		return
	}

	result, err := h.loanSlipService.LoanSlipDetailService(c.Request.Context(), id)
	if err != nil {
		c.Error(&error_middleware.AppError{
			HTTPStatus: http.StatusNotFound,
			Code:       error_middleware.CodeNotFound,
			Message:    "Không tìm thấy",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
