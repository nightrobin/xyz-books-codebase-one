package model

type Publisher struct {
	ID		uint64 `gorm:"primaryKey" json:"ID"`
	Name	string `form:"Name" json:"Name" validate:"required,max=255"`
}