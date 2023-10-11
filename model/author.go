package model

type Author struct {
	ID			uint64 `gorm:"primaryKey"`
	FirstName	string `json:"first_name" form:"first-name" validate:"required,max=255"`
	MiddleName	string `json:"middle_name" form:"middle-name" validate:"max=255"`
	LastName	string `json:"last_name" form:"last-name" validate:"required,max=255"`
}

