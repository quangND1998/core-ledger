package ginhp

import (
	"core-ledger/internal/core"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"core-ledger/model/dto"
)

type Version struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Path string `json:"path"`
}
type System struct {
	Name     string  `json:"name"`
	Mode     string  `json:"mode"`
	Version  Version `json:"version"`
	Timezone string  `json:"timezone"`
}
type Response struct {
	Status  bool              `json:"status"`
	Error   *core.AppError    `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Data    interface{}       `json:"data"`
	System  System            `json:"system"`
}

func GetAccountReq(c *gin.Context) *AccountRequest {
	return c.MustGet(ContextKeyAccountRequest.String()).(*AccountRequest)
}

func GetCustomerReq(c *gin.Context) *CustomerRequest {
	return c.MustGet(ContextKeyCustomerRequest.String()).(*CustomerRequest)
}

func GetEmployeeReq(c *gin.Context) *EmployeeRequest {
	return c.MustGet(ContextKeyEmployeeRequest.String()).(*EmployeeRequest)
}

func ReturnBadRequestError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, dto.PreResponse{
		Error: &dto.ResponseError{
			Message: err.Error(),
		},
	})
}

func ReturnInternalError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, dto.PreResponse{
		Error: &dto.ResponseError{
			Message: err.Error(),
		},
	})
}

func ReturnSuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: data,
	})
}
func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  true,
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		System: System{
			Name: "wealify",
			Mode: "production",
			Version: Version{
				Code: 2,
				Name: "v2.0.0",
				Path: "/v2",
			},
		},
	})
}
func RespondOKPagination(c *gin.Context, data interface{}, page, limit int, total int64) {
	c.JSON(http.StatusOK, ResponsePagination{
		Status:  true,
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Page:    page,
		Limit:   limit,
		Total:   total,
		System: System{
			Name: "wealify",
			Mode: "production",
			Version: Version{
				Code: 2,
				Name: "v2.0.0",
				Path: "/v2",
			},
		},
	})
}

type ResponsePagination struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int64       `json:"total"`
	Data    interface{} `json:"data"`
	System  System      `json:"system"`
}

func RespondOKWithError(c *gin.Context, err error) {
	var appErr *core.AppError
	if !errors.As(err, &appErr) {
		appErr = &core.AppError{
			Message:     appErr.Error(),
			Description: appErr.Description,
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Status:  true,
		Code:    http.StatusOK,
		Message: "error",
		Error:   appErr,
		System: System{
			Name: "wealify",
			Mode: "production",
			Version: Version{
				Code: 2,
				Name: "v2.0.0",
				Path: "/v2",
			},
		},
	})
}
func RespondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, Response{
		Status:  false,
		Code:    code,
		Message: message,
		System: System{
			Name: os.Getenv("APP_NAME"),
			Mode: os.Getenv("MODE"),
			Version: Version{
				Code: 2,
				Name: "v2.0.0",
				Path: "/v2",
			},
		},
	})
}

func RespondErrorValidate(c *gin.Context, code int, message string, errors map[string]string) {
	c.AbortWithStatusJSON(code, Response{
		Status:  false,
		Code:    code,
		Errors:  errors,
		Message: message,
		System: System{
			Name: os.Getenv("APP_NAME"),
			Mode: os.Getenv("MODE"),
			Version: Version{
				Code: 2,
				Name: "v2.0.0",
				Path: "/v2",
			},
		},
	})
}

const (
	BadRequest          = "bad request"
	InternalServerError = "internal server error"
	NotFound            = "not found"
)

var ErrAlreadyRegistered = errors.New("email already registered")
var ErrInternalServer = errors.New("internal server error")

var ErrMessageMap = map[string]map[error]string{
	"vi": {
		ErrAlreadyRegistered: "Email này đã được đăng ký. Vui lòng chọn một email Google khác để tiếp tục.",
		ErrInternalServer:    "Có lỗi xảy ra, vui lòng thử lại sau.",
	},
	"en": {
		ErrAlreadyRegistered: "This email has already been registered. Please choose a different Google email to continue.",
		ErrInternalServer:    "An error has been occur, please try again.",
	},
}

func ToMessageResp(acceptLang string, err error) string {
	acceptLang = strings.ToLower(acceptLang)
	if acceptLang == "vi" {
		msg, ok := ErrMessageMap["vi"][err]
		if !ok {
			return ErrMessageMap["vi"][ErrInternalServer]
		}
		return msg
	} else {
		msg, ok := ErrMessageMap["en"][err]
		if !ok {
			return ErrMessageMap["en"][ErrInternalServer]
		}
		return msg
	}
}
