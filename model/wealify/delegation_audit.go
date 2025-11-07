package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Enum cho Event
type DelegationAuditEvent string

const (
	// Account Management
	CreateVirtualAccount DelegationAuditEvent = "CREATE_VIRTUAL_ACCOUNT"
	UpdateVirtualAccount DelegationAuditEvent = "UPDATE_VIRTUAL_ACCOUNT"
	DeleteVirtualAccount DelegationAuditEvent = "DELETE_VIRTUAL_ACCOUNT"

	// Transaction Events
	MakeTransaction   DelegationAuditEvent = "MAKE_TRANSACTION"
	TopUp             DelegationAuditEvent = "TOP_UP"
	Withdrawal        DelegationAuditEvent = "WITHDRAWAL"
	Transfer          DelegationAuditEvent = "TRANSFER"
	ReceiveTransfer   DelegationAuditEvent = "RECEIVE_TRANSFER"
	ReceiveTopUp      DelegationAuditEvent = "RECEIVE_TOP_UP"
	ReceiveWithdrawal DelegationAuditEvent = "RECEIVE_WITHDRAWAL"

	// Account Status Changes
	AccountActivated   DelegationAuditEvent = "ACCOUNT_ACTIVATED"
	AccountDeactivated DelegationAuditEvent = "ACCOUNT_DEACTIVATED"
	AccountSuspended   DelegationAuditEvent = "ACCOUNT_SUSPENDED"

	// Permission Changes
	PermissionGranted DelegationAuditEvent = "PERMISSION_GRANTED"
	PermissionRevoked DelegationAuditEvent = "PERMISSION_REVOKED"
	RoleChanged       DelegationAuditEvent = "ROLE_CHANGED"

	// Security Events
	LoginSuccess      DelegationAuditEvent = "LOGIN_SUCCESS"
	LoginFailed       DelegationAuditEvent = "LOGIN_FAILED"
	PasswordChanged   DelegationAuditEvent = "PASSWORD_CHANGED"
	TwoFactorEnabled  DelegationAuditEvent = "TWO_FACTOR_ENABLED"
	TwoFactorDisabled DelegationAuditEvent = "TWO_FACTOR_DISABLED"

	// Delegation Management
	DelegationCreated  DelegationAuditEvent = "DELEGATION_CREATED"
	DelegationUpdated  DelegationAuditEvent = "DELEGATION_UPDATED"
	DelegationRevoked  DelegationAuditEvent = "DELEGATION_REVOKED"
	InvitationSent     DelegationAuditEvent = "INVITATION_SENT"
	InvitationAccepted DelegationAuditEvent = "INVITATION_ACCEPTED"
	InvitationExpired  DelegationAuditEvent = "INVITATION_EXPIRED"
)

// Enum cho Status
type DelegationAuditStatus string

const (
	StatusSuccess   DelegationAuditStatus = "SUCCESS"
	StatusFailed    DelegationAuditStatus = "FAILED"
	StatusPending   DelegationAuditStatus = "PENDING"
	StatusCancelled DelegationAuditStatus = "CANCELLED"
)

type DelegationAudit struct {
	ID                  string                `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DelegationAccountID string                `gorm:"type:varchar(255);not null;index:idx_delegation_account_created_at"`
	Event               DelegationAuditEvent  `gorm:"type:varchar(255);not null;index:idx_event_created_at;index:idx_refid_event"`
	Status              DelegationAuditStatus `gorm:"type:varchar(50);not null;default:SUCCESS"`

	Data         datatypes.JSONMap `gorm:"type:jsonb"` // map[string]interface{}
	RefID        *string           `gorm:"type:varchar(255);index:idx_refid_event"`
	RefType      *string           `gorm:"type:varchar(100)"`
	IPAddress    *string           `gorm:"type:varchar(45)"`
	UserAgent    *string           `gorm:"type:text"`
	Metadata     datatypes.JSONMap `gorm:"type:jsonb"`
	ErrorMessage *string           `gorm:"type:text"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (DelegationAudit) TableName() string {
	return "delegation_audit"
}
