package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type WebhookLog struct {
	ID               string          `json:"id"`
	EventName        string          `json:"event_name"`
	WebhookConfigID  string          `json:"webhook_config_id"`
	Status           string          `json:"status"`
	ErrorMsg         string          `json:"error_msg"`
	HeaderXTimestamp string          `json:"header_x_timestamp"`
	HeaderXSignature string          `json:"header_x_signature"`
	Request          json.RawMessage `json:"body"`
	Response         json.RawMessage `json:"response"`
	HttpResponseCode int             `json:"http_response_code"`
	MessageID        string          `json:"message_id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func (w *WebhookLog) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	return nil
}

func (w *WebhookLog) BeforeSave(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	return nil
}

type WebhookConfiguration struct {
	ID         string    `json:"id"`
	OwnerID    int64     `json:"owner_id"`
	Name       string    `json:"name"`
	TargetURL  string    `json:"target_url"`
	PrivateKey string    `json:"private_key"`
	PublicKey  string    `json:"public_key"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (WebhookConfiguration) TableName() string {
	return "webhook_configuration"
}

type EventTrigger struct {
	EventName   string    `gorm:"uniqueIndex" json:"event_name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (EventTrigger) TableName() string {
	return "event_triggers"
}

type WebhookRegisterEvent struct {
	WebhookID    string    `gorm:"primaryKey" json:"webhook_id"`
	EventName    string    `gorm:"primaryKey" json:"event_name"`
	SubscribedAt time.Time `json:"subscribed_at"`
}

func (WebhookRegisterEvent) TableName() string {
	return "webhook_register_events"
}
