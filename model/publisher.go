package model

type Publisher struct {
	ID		uint64 `gorm:"primaryKey"`
	Name	string `form:"Name" json:"Name" validate:"required,max=255"`
}