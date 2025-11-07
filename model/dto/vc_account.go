package dto

type CreateVCAccountRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password" binding:"required,min=6"`
	CountryCode string `json:"country_code" binding:"required"`
}

type VCAccountResponse struct {
	ID           int32  `json:"id"`
	CustomerID   string `json:"customer_id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	AccountType  string `json:"account_type"`
	AccountLevel string `json:"account_level"`
	IsVcCustomer bool   `json:"is_vc_customer"`
}
