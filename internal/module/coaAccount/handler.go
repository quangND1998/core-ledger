package coaaccount

import (
	"bytes"
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"core-ledger/pkg/utils"
	"core-ledger/pkg/utils/helper"
	"encoding/json"
	"fmt"
	"time"

	// "encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type CoaAccountHandler struct {
	logger        logger.CustomLogger
	service       *CoaAccountService
	coAccountRepo repo.CoAccountRepo
	dispatcher    queue.Dispatcher
}

func NewCoaAccountHandler(service *CoaAccountService, coAccountRepo repo.CoAccountRepo, dispatcher queue.Dispatcher) *CoaAccountHandler {
	return &CoaAccountHandler{
		logger:        logger.NewSystemLog("CoaAccountHandler"),
		service:       service,
		coAccountRepo: coAccountRepo,
		dispatcher:    dispatcher,
	}
}

func (h *CoaAccountHandler) List(c *gin.Context) {
	// TODO implement me
	q := &dto.ListCoaAccountFilter{}
	err := c.ShouldBindQuery(&q)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("ListCoaAccountFilter request", q)
	res, err := h.coAccountRepo.PaginateWithScopes(c, q)
	h.logger.Info("ListCoaAccountFilter res", res)
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *CoaAccountHandler) GetCoaAccountDetail(c *gin.Context) {
	h.logger.Info("ListCoaAccountFilter request")
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid SystemPaymentID"})
		return
	}
	res, err := h.service.GetCoaAccountDetail(c, id)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *CoaAccountHandler) ExportCoaAccounts(c *gin.Context) {
	// --- Bind JSON request ---
	var req *ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Info("ExportRequest err", err)
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Select) == 0 {
		ginhp.RespondError(c, http.StatusBadRequest, "export file must have at least one header key")
		return
	}

	// --- Parse query filter ---
	var filter *dto.ListCoaAccountFilter
	queryBytes, _ := json.Marshal(req.Query)
	if err := json.Unmarshal(queryBytes, &filter); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("transactionFilter", filter)

	// --- Fetch data ---
	data, err := h.coAccountRepo.Paginate(filter)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(data.Items) == 0 {
		return
	}

	// --- Create Excel file ---
	f := excelize.NewFile()
	const sheet = "Title"
	_ = f.SetSheetName("Sheet1", sheet)

	// Set header
	for col, key := range req.Select {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		_ = f.SetCellValue(sheet, cell, key)
	}

	// Set data rows
	for rowIndex, tran := range data.Items {
		for colIndex, head := range req.Select {
			err = h.getDataHeader(f, sheet, req.Select)
			if err != nil {
				ginhp.RespondError(c, http.StatusBadRequest, err.Error())
				return
			}
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)

			var value any
			switch head {
			case CoaAccountExportKeyIndex:
				value = rowIndex + 1
			case CoaAccountExportKeyCode:
				value = tran.Code
			case CoaAccountExportKeyName:
				value = tran.Name
			case CoaAccountExportKeyType:
				value = tran.Type
			case CoaAccountExportKeyParentCode:
				value = tran.Parent.Code
			case CoaAccountExportKeyStatus:
				value = tran.Status
			case CoaAccountExportKeyProvider:
				value = tran.Provider
			case CoaAccountExportKeyNetwork:
				value = tran.Network
			case CoaAccountExportKeyTags:
				value = helper.FormatJSONForExcel(tran.Tags)
			case CoaAccountExportKeyMetadata:
				value = helper.FormatJSONForExcel(tran.Metadata)
			case CoaAccountExportKeyCreatedAt:
				loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
				value = tran.CreatedAt.In(loc).Format("2006-01-02 15:04:05")
			case CoaAccountExportKeyUpdatedAt:
				loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
				value = tran.UpdatedAt.In(loc).Format("2006-01-02 15:04:05")
			}

			if value == nil {
				value = ""
			}
			_ = f.SetCellValue(sheet, cell, value)
		}
		_ = f.SetRowHeight(sheet, rowIndex+2, 35)
	}

	// --- Adjust column width ---
	for i := 1; i <= len(req.Select); i++ {
		colName, _ := excelize.ColumnNumberToName(i)
		text, _ := f.GetCellValue(sheet, fmt.Sprintf("%s2", colName))
		_ = f.SetColWidth(sheet, colName, colName, float64(len(text))*3)
	}

	// --- Set body style ---
	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11, Color: "000000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "DDDDDD", Style: 1},
			{Type: "top", Color: "DDDDDD", Style: 1},
			{Type: "bottom", Color: "DDDDDD", Style: 1},
			{Type: "right", Color: "DDDDDD", Style: 1},
		},
	})
	lastCol, _ := excelize.ColumnNumberToName(len(req.Select))
	_ = f.SetCellStyle(sheet, "A2", fmt.Sprintf("%s%d", lastCol, len(data.Items)+1), bodyStyle)

	// --- Prepare file name ---
	loc := time.FixedZone("UTC+7", 7*60*60)
	const layout = "02-01-2006"
	timestamp := time.Now().In(loc).Format(layout)

	if req.Query.StartDate != nil && req.Query.EndDate != nil {
		startDate, err1 := time.Parse("2006-01-02", *req.Query.StartDate)
		endDate, err2 := time.Parse("2006-01-02", *req.Query.EndDate)
		if err1 == nil && err2 == nil {
			startFormatted := startDate.In(loc).Format(layout)
			endFormatted := endDate.In(loc).Format(layout)
			if startFormatted == endFormatted {
				timestamp = startFormatted
			} else {
				timestamp = fmt.Sprintf("%s-%s", startFormatted, endFormatted)
			}
		}
	}

	fileName := req.FileName
	if fileName == "" {
		fileName = "Export-Transaction"
	}
	downloadName := fmt.Sprintf("%s-%s.xlsx", fileName, timestamp)
	h.logger.Info("Exporting Excel for download: " + downloadName)
	// --- Write Excel to buffer and return ---
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", downloadName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

func (h *CoaAccountHandler) getDataHeader(f *excelize.File, sheetName string, keys []CoaAccountExportKey) error {
	if len(keys) == 0 {
		// giả sử TransactionExportKeys có thể được định nghĩa sẵn nếu cần
		return nil
	}
	headerData := make([]string, len(keys))
	for i, key := range keys {
		headerData[i] = TransactionExportHeaders[key]
	}

	// Ghi hàng header (bắt đầu từ hàng 1)
	rowIndex := 1
	for i, header := range headerData {
		cell, _ := excelize.CoordinatesToCellName(i+1, rowIndex)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style cho header
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
			WrapText:   true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1, // solid fill
			Color:   []string{"#D3D3D3"},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// Áp dụng style cho toàn hàng
	colCount := len(headerData)
	startCell, _ := excelize.CoordinatesToCellName(1, rowIndex)
	endCell, _ := excelize.CoordinatesToCellName(colCount, rowIndex)
	f.SetCellStyle(sheetName, startCell, endCell, style)

	// Đặt chiều cao hàng (40)
	f.SetRowHeight(sheetName, rowIndex, 40)

	return nil
}
