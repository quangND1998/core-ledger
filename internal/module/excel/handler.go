package excel

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/utils"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

type ExcelHandler struct {
	logger  logger.CustomLogger
	service *excelService
}
type FileExcelRequest struct {
	File *multipart.FileHeader `form:"file"  validate:"required|mimeTypes:image/jpeg,image/jpg,image/png|maxSize:2MB"`
}

func NewExcelHandler(service *excelService) *ExcelHandler {
	return &ExcelHandler{
		logger:  logger.NewSystemLog("TransactionHandler"),
		service: service,
	}
}

func (h *ExcelHandler) ImportCoAccounts(c *gin.Context) {

	file, _ := c.FormFile("file")
	h.logger.Info("ImportCoAccounts request: %+v", file.Filename)
	req := FileExcelRequest{File: file}
	h.logger.Info("ImportCoAccounts request: %+v", req)
	v := validate.Struct(req)

	v.AddMessages(map[string]string{

		"File.required":  "Trường {field} là bắt buộc",
		"File.mimeTypes": "Chỉ chấp nhận file định dạng xlsx,xls,csv",
		"File.maxSize":   "Kích thước file không được vượt quá 2MB",
	})
	if !v.Validate() {

		errsFlatten := utils.ErrsFlatten(v.Errors.All())

		c.JSON(http.StatusBadRequest, gin.H{"errors": errsFlatten})
		return
	}
	tmpPath := fmt.Sprintf("/tmp/%s", file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save file"})
		return
	}
	h.service.ImportCoAccounts(tmpPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "File queued for import",
	})
}
