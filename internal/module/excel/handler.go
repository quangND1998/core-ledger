package excel

import (
	"core-ledger/pkg/logger"

	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type ExcelHandler struct {
	logger  logger.CustomLogger
	service *ExcelService
}

func NewExcelHandler(service *ExcelService) *ExcelHandler {
	return &ExcelHandler{
		logger:  logger.NewSystemLog("TransactionHandler"),
		service: service,
	}
}

func (h *ExcelHandler) ImportCoAccounts(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Vui lòng chọn file để tải lên",
		})
		return
	}
	h.logger.Info("ImportCoAccounts file:", file)
	rules := govalidator.MapData{
		"file:file": []string{
			"required",
			"ext:xlsx,xls",
			"mime:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-excel,application/zip",
			"size:2048000", // giới hạn 2MB
		},
	}

	messages := govalidator.MapData{
		"file:file": []string{
			"required:Vui lòng chọn file để tải lên",
			"ext:Chỉ chấp nhận file xlsx và xls",
			"mime:Định dạng file không hợp lệ",
			"size:Kích thước file không được vượt quá 2MB",
		},
	}

	opts := govalidator.Options{
		Request:  c.Request,
		Rules:    rules,
		Messages: messages,
	}

	v := govalidator.New(opts)
	e := v.Validate()
	if len(e) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": e})
		return
	}

	tmpPath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu file"})
		return
	}

	h.logger.Info("File uploaded:", file.Filename)
	h.service.ImportCoAccounts(c, tmpPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "File đã được tải lên và đang chờ xử lý",
	})
}
