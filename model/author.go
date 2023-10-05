package model

type Author struct {
	ID			uint64 `gorm:"primaryKey"`
	FirstName	string `json:"first_name"`
	MiddleName	string `json:"middle_name"`
	LastName	string `json:"last_name"`
}
