package coaaccount

import (
	"bytes"
	"core-ledger/model/dto"
	model "core-ledger/model/wealify"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"core-ledger/pkg/utils"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"
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
	res, err := h.coAccountRepo.Paginate(q)
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
	var b *ExportRequest // khởi tạo trước

	err := c.ShouldBindJSON(&b)
	if err != nil {
		h.logger.Info("ExportRequest err", err)
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(b.Select) == 0 {
		ginhp.RespondError(c, http.StatusBadRequest, "export file must have at least one header key")
		return
	}

	var transactionFilter *dto.ListCoaAccountFilter
	queryBytes, err := json.Marshal(b.Query)
	h.logger.Info("transactionFilter err", err)
	if err != nil {

		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(queryBytes, &transactionFilter)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("transactionFilter", transactionFilter)
	data, err := h.coAccountRepo.Paginate(transactionFilter)

	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(data.Items) == 0 {
		return
	}

	// Create excel file
	f := excelize.NewFile()
	sheet := "Title"
	err = f.SetSheetName("Sheet1", sheet)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	for col, key := range b.Select {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		err = f.SetCellValue(sheet, cell, key)
		if err != nil {
			ginhp.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	rowIndex := 1
	for _, tran := range data.Items {
		for colIndex, head := range b.Select {
			err = h.getDataHeader(f, sheet, b.Select)
			if err != nil {
				ginhp.RespondError(c, http.StatusBadRequest, err.Error())
				return
			}
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			var value any
			switch head {
			case CoaAccountExportKeyIndex:
				value = rowIndex
			case CoaAccountExportKeyCode:
				value = tran.Code
			
			if value == nil {
				value = ""
			}
			_ = f.SetCellValue(sheet, cell, value)
			//if err != nil {
			//	h.logger.Error(rowIndex, colIndex)
			//	ginhp.RespondError(c, http.StatusBadRequest, err.Error())
			//	return
			//}
		}
		_ = f.SetRowHeight(sheet, rowIndex, 35)
		rowIndex++
	}
	lastColName, _ := excelize.ColumnNumberToName(len(b.Select))
	// set col witdh
	for i := 1; i < len(b.Select); i++ {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		text, _ := f.GetCellValue(sheet, fmt.Sprintf("%s%d", colName, 2))
		_ = f.SetColWidth(sheet, colName, colName, float64(len(text))*3)
	}

	// set data style
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
	_ = f.SetCellStyle(sheet, "A2", fmt.Sprintf("%s%d", lastColName, len(data.Items)+2), bodyStyle)

	loc := time.FixedZone("UTC+7", 7*60*60)
	const layout = "02-01-2006"

	// default = today
	timestampFileName := time.Now().In(loc).Format(layout)

	if b.Query.StartDate != nil && b.Query.EndDate != nil {
		startDate, err1 := time.Parse("2006-01-02", *b.Query.StartDate)
		endDate, err2 := time.Parse("2006-01-02", *b.Query.EndDate)
		if err1 == nil && err2 == nil {
			startFormatted := startDate.In(loc).Format(layout)
			endFormatted := endDate.In(loc).Format(layout)

			if startFormatted == endFormatted {
				timestampFileName = startFormatted
			} else {
				timestampFileName = fmt.Sprintf("%s-%s", startFormatted, endFormatted)
			}
		}
	}

	var fileName string
	if b.FileName == "" {
		fileName = "Export-Transaction"
	} else {
		fileName = b.FileName
	}
	fullFileName := fmt.Sprintf("export/transaction/%s-%s.xlsx", fileName, timestampFileName)
	h.logger.Info("Exporting to file: " + fullFileName)
	// Save to buffer
	var buf bytes.Buffer
	if err = f.Write(&buf); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
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
