package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" gorm:"unique" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}
