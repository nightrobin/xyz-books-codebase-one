package model

type Book struct {
	ID				uint64		`gorm:"primaryKey" json:"ID"`
	Title			string		`form:"Title" json:"Title" validate:"required"`
	Isbn13			string		`form:"Isbn13" gorm:"column:isbn_13" json:"Isbn13"`
	Isbn10			string		`form:"Isbn10" gorm:"column:isbn_10" json:"Isbn10"`
	PublicationYear	int16 		`form:"PublicationYear" json:"PublicationYear" validate:"required,number,min=1000,max=2023"`
	PublisherID		uint64		`form:"PublisherID" json:"PublisherID" validate:"required"`
	ImageURL		string		`form:"ImageURL" json:"ImageURL"`
	Edition			string		`form:"Edition" json:"Edition"`
	ListPrice		float32		`form:"ListPrice" json:"ListPrice" validate:"required"`
	AuthorIDs		[]uint64	`gorm:"-" json:"AuthorIDs" validate:"-"`
}
