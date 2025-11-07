package model

type User struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	FullName string `gorm:"column:full_name;not null" json:"full_name"`
	Email    string `gorm:"column:email;not null" json:"email"`
}
