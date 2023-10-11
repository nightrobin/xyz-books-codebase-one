package model

type Author struct {
	ID			uint64 `gorm:"primaryKey"`
	FirstName	string `json:"first_name" form:"first-name"`
	MiddleName	string `json:"middle_name" form:"middle-name"`
	LastName	string `json:"last_name" form:"last-name"`
}

