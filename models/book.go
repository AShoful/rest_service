package models

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" validate:"required,min=2"`
	Author string `json:"author" validate:"required,min=2"`
	UserId uint   `json:"userid"`
}
