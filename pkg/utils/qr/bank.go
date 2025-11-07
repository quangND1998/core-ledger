package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/sigurn/crc16"
	"github.com/skip2/go-qrcode"
)

type QRInput struct {
	MerchantID string // Đúng định dạng như bạn gửi: 25 ký tự (đã bao gồm prefix + bankBin + account)
	BankBin    string
	BankNumber string
	Amount     string // VD: "100000"
	Message    string // VD: "Test"
	IsStaticQR bool   // true: QR tĩnh, false: QR động
}

func formatTag(tag string, value string) string {
	return fmt.Sprintf("%s%02d%s", tag, len(value), value)
}

func buildEMVCoString(input QRInput) string {
	payloadFormat := formatTag("00", "01")
	method := "12"
	if input.IsStaticQR {
		method = "11"
	}
	initiationMethod := formatTag("01", method)

	gui := formatTag("00", "A000000727")
	accountID := formatTag("01", "0006"+input.BankBin+"0111"+input.BankNumber)
	serviceCode := formatTag("02", "QRIBFTTA")
	merchantAccount := formatTag("38", gui+accountID+serviceCode)

	currency := formatTag("53", "704")
	amount := ""
	if input.Amount != "" {
		amount = formatTag("54", input.Amount)
	}
	country := formatTag("58", "VN")

	additional := ""
	if input.Message != "" {
		message := formatTag("08", input.Message)
		additional = formatTag("62", message)
	}

	qrBase := payloadFormat + initiationMethod + merchantAccount + currency + amount + country + additional
	crcBase := qrBase + "6304"

	crcTable := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum([]byte(crcBase), crcTable)
	crcStr := fmt.Sprintf("%04X", crc)

	return crcBase + strings.ToUpper(crcStr)
}

func main() {
	input := QRInput{
		MerchantID: "", // Độ dài 25 ký tự y chang mẫu bạn đưa
		BankBin:    "970436",
		BankNumber: "0711000270243",
		Amount:     "100000",
		Message:    "Test 1 2 3 4 5",
		IsStaticQR: false,
	}

	qrString := buildEMVCoString(input)
	fmt.Println("✅ EMVCo QR string:")
	fmt.Println(qrString)

	err := qrcode.WriteFile(qrString, qrcode.High, 512, "emvco_qr.png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ QR đã được lưu tại: emvco_qr.png")
}
