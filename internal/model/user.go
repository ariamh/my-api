package model

type User struct {
	Base
	Name     string `json:"name" gorm:"size:100;not null"`
	Email    string `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:255;not null"`
	Role     string `json:"role" gorm:"size:20;default:user"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

func (User) TableName() string {
	return "users"
}