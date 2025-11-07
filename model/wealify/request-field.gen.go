package model

const TableNameRequestField = "request-field"

// RequestField mapped from table <request-field>
type RequestField struct {
	RequestID string `gorm:"column:request_id;primaryKey" json:"request_id"`
	FieldID   string `gorm:"column:field_id;primaryKey" json:"field_id"`
}

// TableName RequestField's table name
func (*RequestField) TableName() string {
	return TableNameRequestField
}
