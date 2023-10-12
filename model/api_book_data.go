package model

type ApiBookData struct {
	ID				uint64		`json:"ID"`
	Title			string		`json:"Title"`
	Isbn13			string		`json:"Isbn13" gorm:"column:isbn_13"`
	Isbn10			string		`json:"Isbn10" gorm:"column:isbn_10"`
	PublicationYear	int16		`json:"PublicationYear"`
	PublisherID		uint64		`json:"PublisherID"`
	ImageURL		string		`json:"ImageURL"`
	Edition			string		`json:"Edition"`
	ListPrice		float32		`json:"ListPrice"`
	AuthorIDs		string		`json:"AuthorIDs" gorm:"column:author_ids"`
}
