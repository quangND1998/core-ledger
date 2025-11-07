package dto

type CreateWebhookRequest struct {
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint" binding:"required"`
	Events   []string `json:"events" binding:"required"`
}
type UpdateWebhookRequest struct {
	ID       string   `json:"-"`
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint"`
	Events   []string `json:"events"`
}

type CreateWebhookResponse struct {
	ID        string `json:"id"`
	Endpoint  string `json:"endpoint"`
	PublicKey string `json:"public_key"`
}
type GetTriggerEvents struct {
	Total  int64    `json:"total"`
	Events []string `json:"items"`
}
type ListWebhookRequest struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}
type ListWebhookResponse struct {
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	TotalPage int               `json:"total_page"`
	Items     []WebhookItemResp `json:"items"`
}

type WebhookItemResp struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Endpoint         string `json:"endpoint"`
	PublicKey        string `json:"public_key"`
	SubscribedEvents []struct {
		EventName    string `json:"event_name"`
		SubscribedAt string `json:"subscribed_at"`
	} `json:"subscribed_events"`
}

type ListWebhookLogRequest struct {
	HookURL   string `json:"hook_url"`
	EventName string `json:"event_name"`
	Status    string `json:"status"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	FromUnix  int64  `json:"from_unix"`
	ToUnix    int64  `json:"to_unix"`
}
type ListWebhookLogResponse struct {
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	TotalPage int         `json:"total_page"`
	Limit     int         `json:"limit"`
	Items     interface{} `json:"items"`
}
