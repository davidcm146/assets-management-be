package validator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func HandleValidationError(err error) map[string]string {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, fieldError := range validationError {
			switch fieldError.Tag() {
			case "required":
				errors[fieldError.Field()] = fieldError.Field() + " là bắt buộc"
			case "min":
				errors[fieldError.Field()] = fieldError.Field() + " phải có độ dài tối thiểu " + fieldError.Param()
			case "max":
				errors[fieldError.Field()] = fieldError.Field() + " phải có độ dài tối đa " + fieldError.Param()
			case "oneof":
				errors[fieldError.Field()] = fieldError.Field() + " phải là một trong các giá trị: " + fieldError.Param()
			case "gte":
				errors[fieldError.Field()] = fieldError.Field() + " phải lớn hơn hoặc bằng " + fieldError.Param()
			case "lte":
				errors[fieldError.Field()] = fieldError.Field() + " phải nhỏ hơn hoặc bằng " + fieldError.Param()
			default:
				errors[fieldError.Field()] = "Giá trị không hợp lệ cho " + fieldError.Field()
			}
		}
		return errors
	}
	return nil
}

func IsValidImageType(mime string) bool {
	file, _, err := strings.Cut(mime, "/")
	if err != true || file != "image" {
		return false
	}
	allowedTypes := []string{"jpeg", "png", "jpg"}

	for _, allowedType := range allowedTypes {
		if strings.HasSuffix(mime, allowedType) {
			return true
		}
	}
	return false
}

func ValidateImage() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return fmt.Errorf("Lỗi khởi tạo trình xác thực")
	}

	v.RegisterValidation("images", func(fl validator.FieldLevel) bool {
		images, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}
		for _, img := range images {
			if len(img) == 0 {
				return false
			}
			if !IsValidImageType(img) {
				return false
			}
		}
		return true
	})
	return nil
}
