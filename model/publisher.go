package model

type Publisher struct {
	ID		uint64 `gorm:"primaryKey"`
	Name	string `json:"name" form:"name"`
}
