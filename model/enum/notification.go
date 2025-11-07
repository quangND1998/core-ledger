package enum

type NotificationChannel string

const (
	NotificationChannelInApp NotificationChannel = "IN_APP"
	NotificationChannelEmail NotificationChannel = "EMAIL"
	NotificationChannelLark  NotificationChannel = "LARK"
)

type NotificationStatus string

// Notification Status
const (
	NotificationStatusCreated   NotificationStatus = "CREATED"
	NotificationStatusSent      NotificationStatus = "SENT"
	NotificationStatusDelivered NotificationStatus = "DELIVERED"
	NotificationStatusRead      NotificationStatus = "READ"
	NotificationStatusFailed    NotificationStatus = "FAILED"
)

func (n NotificationStatus) String() string {
	return string(n)
}

type NotificationGroup string

const (
	NotificationGroupSystem    NotificationGroup = "SYSTEM"
	NotificationGroupImportant NotificationGroup = "IMPORTANT"
)

func (n NotificationGroup) String() string {
	return string(n)
}

func NotificationGroupValues() []NotificationGroup {
	return []NotificationGroup{
		NotificationGroupSystem,
		NotificationGroupImportant,
	}
}

type NotificationType string

const (
	NotificationTypeWelcome                     NotificationType = "WELCOME"
	NotificationTypeNewDeviceLogin              NotificationType = "NEW_DEVICE_LOGIN"
	NotificationTypeUpdateAccountLevel          NotificationType = "UPDATE_ACCOUNT_LEVEL"
	NotificationTypeUpdateAccountType           NotificationType = "UPDATE_ACCOUNT_TYPE"
	NotificationTypeUpgradeTier                 NotificationType = "UPGRADE_TIER"
	NotificationTypeDowngradeTier               NotificationType = "DOWNGRADE_TIER"
	NotificationTypeEnableWealifyWalletFeature  NotificationType = "ENABLE_WEALIFY_WALLET_FEATURE"
	NotificationTypeDisableWealifyWalletFeature NotificationType = "DISABLE_WEALIFY_WALLET_FEATURE"
	NotificationTypeAdjustBalance               NotificationType = "ADJUST_BALANCE"
	NotificationTypeRejectKyc                   NotificationType = "REJECT_KYC"
	NotificationTypeApproveKyc                  NotificationType = "APPROVE_KYC"
	NotificationTypeRejectKyb                   NotificationType = "REJECT_KYB"
	NotificationTypeApproveKyb                  NotificationType = "APPROVE_KYB"
	NotificationTypeApprovePayment              NotificationType = "APPROVE_PAYMENT"
	NotificationTypeRejectPayment               NotificationType = "REJECT_PAYMENT"
	NotificationTypeCreateVirtualAccount        NotificationType = "CREATE_VIRTUAL_ACCOUNT"
	NotificationTypeApproveWithdraw             NotificationType = "APPROVE_WITHDRAW"
	NotificationTypeRejectWithdraw              NotificationType = "REJECT_WITHDRAW"
	NotificationTypeApproveTopUp                NotificationType = "APPROVE_TOP_UP"
	NotificationTypeApproveAdjust               NotificationType = "APPROVE_ADJUST"
)

func (n NotificationType) String() string {
	return string(n)
}
