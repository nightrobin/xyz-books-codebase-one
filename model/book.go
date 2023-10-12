package model

type Book struct {
	ID				uint64		`gorm:"primaryKey"`
	Title			string		`form:"title" json:"title"`
	Isbn13			string		`form:"isbn-13" gorm:"column:isbn_13" json:"isbn_13"`
	Isbn10			string		`form:"isbn-10" gorm:"column:isbn_10" json:"isbn_10"`
	PublicationYear	int16 		`form:"publication-year" json:"publication_year"`
	PublisherID		uint64		`form:"publisher-id" json:"publisher_id"`
	ImageURL		string		`form:"image-url" json:"image_url"`
	Edition			string		`form:"edition" json:"edition"`
	ListPrice		float32		`form:"list-price" json:"list_price"`
	AuthorIDs		[]uint64	`gorm:"-" json:"author_ids" validate:"required"`
}
