package middleware

import (
	"bytes"
	config "core-ledger/configs"
	"core-ledger/model/dto"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func WrapResponse(c *gin.Context) {
	blw := &CustomResponseWriter{
		body:           bytes.NewBuffer([]byte{}),
		ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()
	var original map[string]any
	var originalArr []map[string]any
	var originalArrInterface []any
	var message string
	if err := json.Unmarshal(blw.body.Bytes(), &original); err != nil {
		log.Error(err)
	}
	if original == nil {
		_ = json.Unmarshal(blw.body.Bytes(), &originalArrInterface)
	}
	if originalArrInterface == nil {
		_ = json.Unmarshal(blw.body.Bytes(), &originalArr)
	}
	if originalArr == nil {
		_ = json.Unmarshal(blw.body.Bytes(), &message)
	}
	var data any = original
	if original == nil {
		data = originalArr
		if originalArr == nil {
			data = originalArrInterface
		}
	}
	newResponse := dto.BaseResponse{
		Code: c.Writer.Status(),
		PreResponse: dto.PreResponse{
			Data: data,
			Error: &dto.ResponseError{
				Message: message,
			},
		},
		System: dto.System{
			Name: config.GetConfig().Common.Name,
			Mode: dto.AppMode(config.GetConfig().Common.Mode),
			Version: dto.Version{
				Code: config.GetConfig().Version.Code,
				Name: config.GetConfig().Version.Name,
				Path: config.GetConfig().Version.Path,
			},
			Timezone: time.Now().In(time.UTC),
		},
	}
	b, _ := json.Marshal(newResponse)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(c.Writer.Status())
	blw.ResponseWriter.Write(b)
}
