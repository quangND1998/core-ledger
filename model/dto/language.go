package dto

type LanguageInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Locale string `json:"locale"`
	Flag   bool   `json:"flag"`
	Status bool   `json:"status"`
}
