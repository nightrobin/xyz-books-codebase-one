package model

type Book struct {
	ID				uint64	`gorm:"primaryKey"`
	Title			string	`json:"title"`
	Isbn13			string	`gorm:"column:isbn_13" json:"isbn_13"`
	Isbn10			string	`gorm:"column:isbn_10" json:"isbn_10"`
	PublicationYear	int16 	`json:"publication_year"`
	PublisherID		uint64	`json:"publisher_id"`
	ImageURL		string	`json:"image_url"`
	Edition			string	`json:"edition"`
	ListPrice		float32	`json:"list_price"`
}
