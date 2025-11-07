package dto

import "time"

type TimeRangeFilter struct {
	From time.Time `form:"from" time_format:"2006-01-02"`
	To   time.Time `form:"to" time_format:"2006-01-02"`
}

type PreResponse struct {
	Data  interface{}    `json:"data"`
	Error *ResponseError `json:"error"`
}

type ResponseError struct {
	Message  string `json:"message"`
	Property string `json:"property"`
}

type BaseResponse struct {
	PreResponse
	Code   int    `json:"code"`
	System System `json:"system"`
}

type AppMode string

const (
	AppModeProduction  AppMode = "production"
	AppModeStaging     AppMode = "staging"
	AppModeDevelopment AppMode = "development"
)

type Version struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type System struct {
	Name     string    `json:"name"`
	Mode     AppMode   `json:"mode"`
	Version  Version   `json:"version"`
	Timezone time.Time `json:"timezone"`
}
