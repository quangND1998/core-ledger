package validate

import (
	model "core-ledger/model/core-ledger"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ProvideValidator đăng ký custom validator với DB
func ProvideValidator(db *gorm.DB) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.RegisterValidation("unique_ruleCode", UniqueCodeValidator(db))
	}
	return nil
}
func UniqueCodeValidator(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		code := fl.Field().String()
		if code == "" {
			return true // required đã check riêng
		}

		var count int64
		// Giả sử bảng của bạn tên "rule_values" và cột "code"
		if err := db.Model(&model.RuleValue{}).Where("code = ?", code).Count(&count).Error; err != nil {
			return false
		}
		return count == 0
	}
}
func ValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "email":
		return "Invalid email address"
	case "gte":
		return fe.Field() + " must be greater than or equal to " + fe.Param()
	case "lte":
		return fe.Field() + " must be less than or equal to " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}
