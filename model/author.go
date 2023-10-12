package model

type Author struct {
	ID			uint64	`gorm:"primaryKey" json:"ID"`
	FirstName	string	`form:"FirstName" json:"FirstName" validate:"required,max=255"`
	MiddleName	string	`form:"MiddleName" json:"MiddleName" validate:"max=255"`
	LastName	string	`form:"LastName" json:"LastName" validate:"required,max=255"`
	IsSelected	bool	`gorm:"-" json:"-"`
}

