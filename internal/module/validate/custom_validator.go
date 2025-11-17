package validate

import (
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ProvideValidator đăng ký custom validator với DB
func ProvideValidator(db *gorm.DB) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.RegisterValidation("unique_ruleCode", UniqueRuleCodeIgnoreID(db))
	}
	return nil
}
func UniqueRuleCodeIgnoreID(db *gorm.DB) validator.Func {
	// To keep track of values in current payload
	payloadValues := map[string]bool{}

	return func(fl validator.FieldLevel) bool {
		item, ok := fl.Parent().Interface().(dto.RuleValueRequest)
		if !ok {
			return false
		}

		// Skip deleted items
		if item.IsDelete != nil && *item.IsDelete {
			return true
		}

		// Check in payload first
		if _, exists := payloadValues[item.Value]; exists {
			return false
		}
		payloadValues[item.Value] = true

		// Check in DB only for items not marked deleted
		var count int64
		query := db.Model(&model.RuleValue{}).Where("value = ? ", item.Value)
		if item.ID != 0 {
			query = query.Where("id != ?", item.ID)
		}
		if err := query.Count(&count).Error; err != nil {
			return false
		}
		return count == 0
	}
}
func FormatErrorMessage(req interface{}, err error) map[string]string {
	out := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		reqType := reflect.TypeOf(req)
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}

		for _, fe := range ve {
			ns := fe.Namespace() // e.g. SaveRuleValueRequest.Data[0].Name
			parts := strings.Split(ns, ".")

			// Bỏ tên struct gốc
			if len(parts) > 0 && parts[0] == reqType.Name() {
				parts = parts[1:]
			}

			currType := reqType
			for i, part := range parts {
				idx := ""
				if strings.Contains(part, "[") && strings.Contains(part, "]") {
					// lấy index
					idxStart := strings.Index(part, "[")
					idxEnd := strings.Index(part, "]")
					idx = part[idxStart+1 : idxEnd] // chỉ lấy số
					part = part[:idxStart]          // Data
				}

				// Map field name -> json tag
				if currType.Kind() == reflect.Struct {
					if field, ok := currType.FieldByName(part); ok {
						tag := field.Tag.Get("json")
						if tag != "" && tag != "-" {
							part = strings.Split(tag, ",")[0]
						}
						currType = field.Type
						if currType.Kind() == reflect.Ptr {
							currType = currType.Elem()
						}
						if currType.Kind() == reflect.Slice {
							currType = currType.Elem()
							if currType.Kind() == reflect.Ptr {
								currType = currType.Elem()
							}
						}
					}
				}

				// Nếu có index, nối bằng "." thay vì [i]
				if idx != "" {
					part = part + "." + idx
				}

				parts[i] = part
			}

			key := strings.Join(parts, ".")
			out[key] = ValidationErrorMessage(fe)
		}
	} else {
		out["error"] = err.Error()
	}

	return out
}

func ValidationErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "email":
		return "Invalid email address"
	case "gte":
		return field + " must be greater than or equal to " + fe.Param()
	case "lte":
		return field + " must be less than or equal to " + fe.Param()
	case "unique":
		return field + " must be unique"
	case "boolean":
		return field + " must be of type " + fe.Param()
	case "unique_ruleCode":
		return field + " is already in use"
	default:
		return field + " is invalid"
	}
}
