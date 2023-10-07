package model

type Book struct {
	ID				uint64	`gorm:"primaryKey"`
	Title			string	`json:"title" form:"title"`
	Isbn13			string	`gorm:"column:isbn_13" json:"isbn_13" form:"isbn-13"`
	Isbn10			string	`gorm:"column:isbn_10" json:"isbn_10" form:"isbn-10"`
	PublicationYear	int16 	`json:"publication_year" form:"publication-year"`
	PublisherID		uint64	`json:"publisher_id" form:"publisher-id"`
	ImageURL		string	`json:"image_url" form:"image-url"`
	Edition			string	`json:"edition" form:"edition"`
	ListPrice		float32	`json:"list_price" form:"list-price"`
}
