package method

import(
	"fmt"
	"encoding/json"
	"html/template"
	"net/http"
	"log"
	"strconv"
	"strings"
	"sync"

	
	"xyz-books/model"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type BookDisplay struct {
	ID				uint64
	Title			string
	Author			string
	Isbn13			string	`gorm:"column:isbn_13"`
	Isbn10			string	`gorm:"column:isbn_10"`
	PublicationYear	int16
	PublisherName	string
	Edition			string
	ListPrice		float32
	ImageURL		string
}

type AuthorDisplay struct {
	ID			uint64 `gorm:"primaryKey"`
	FirstName	string
	MiddleName	string
	LastName	string
	IsSelected  bool
}

func UIBookIndex(c *gin.Context) {

	keyword := c.DefaultQuery("keyword", "")
	keyword = strings.TrimLeft(keyword, " ") 
	keyword = strings.TrimRight(keyword, " ") 
	keyword = strings.NewReplacer(`'`, `\'`, `"`, `\"`).Replace(keyword) 
	
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	var currentPageNumber int64 = int64(pageNumber)
	pageNumber = (pageNumber - 1) * recordLimitPerPage
	limit := recordLimitPerPage

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()
		
		var books []BookDisplay
		var count int64
		
		if (len(keyword) != 0){
		
			havingString := " title LIKE '%" + keyword + "%' " 
			havingString = havingString + " or isbn_13 LIKE '%" + keyword + "%' " 
			havingString = havingString + " or isbn_10 LIKE '%" + keyword + "%' " 
			havingString = havingString + " OR author LIKE '%" + keyword + "%' "
			havingString = havingString + " OR publication_year LIKE '%" + keyword + "%' "
			havingString = havingString + " OR publisher_name LIKE '%" + keyword + "%' "
			
			Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Having(havingString).Limit(limit).Offset(pageNumber).Find(&books)
			Db.Table("books b").Select("Count(*) count", "b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Having(havingString).Count(&count)
		
		} else {
			Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Limit(limit).Offset(pageNumber).Find(&books)
			Db.Table("books b").Select("Count(*) count", "b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Count(&count)
		}


		type PageData struct {
			Keyword string
			Books []BookDisplay
			PageNumbers []int64
			CountShownPageNumber int64
			CurrentPage int64
			PreviousPageNumber int64
			NextPageNumber int64
			MaxPageNumber int64
			IsNextEnabled bool
		}
		
		var data PageData
		data.Books = books
		
		numberOfPages := count / int64(recordLimitPerPage)

		var pageNumbers []int64
		var countShownPageNumber int64 = int64(countShownPageNumber)
		var minPageNumber int64 = currentPageNumber - (countShownPageNumber / 2) 
		var maxPageNumber int64 = currentPageNumber + (countShownPageNumber / 2)

		if minPageNumber < 1 {
			minPageNumber = 1
			maxPageNumber = countShownPageNumber
		}

		if numberOfPages < countShownPageNumber {
			minPageNumber = 1
			maxPageNumber = numberOfPages
		}

		if maxPageNumber > numberOfPages {
			minPageNumber = numberOfPages - countShownPageNumber
			maxPageNumber = numberOfPages
		} 

		for i := minPageNumber; i <= maxPageNumber; i++ {
			pageNumbers = append(pageNumbers, i)
		}

		data.Keyword = keyword
		data.PageNumbers = pageNumbers
		data.CountShownPageNumber = countShownPageNumber
		data.CurrentPage = currentPageNumber
		data.PreviousPageNumber = currentPageNumber - 1
		data.NextPageNumber = currentPageNumber + 1
		data.MaxPageNumber = maxPageNumber
		data.IsNextEnabled = true

		if (numberOfPages - currentPageNumber) <= (countShownPageNumber / 2) {
			data.IsNextEnabled = false
		}

		w := c.Writer

		parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(parsedIndexTemplate, err)
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}()
	
	wg.Wait()

	return
}

func UIAddBookForm(c *gin.Context) {
	
	var authors []model.Author
	Db.Find(&authors)

	var publishers []model.Publisher
	Db.Find(&publishers)

	type PageData struct {
		Authors []model.Author
		Publishers []model.Publisher
	}

	var data PageData
	data.Authors = authors
	data.Publishers = publishers
	
	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/add_form.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func UISubmitAddBookForm(c *gin.Context) {
	var book model.Book
	c.ShouldBind(&book)

	type authorIDs struct {
		AuthorIDs []uint64 `form:"author-ids[]"`
	}

	var bookAuthorIDs authorIDs
	c.ShouldBind(&bookAuthorIDs)

	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully added a Book"
	pageData.HasError = false

	var countIsbn13 int64
	Db.Table("books").Where("isbn_13 = ?", book.Isbn13).Count(&countIsbn13)
	if countIsbn13 > 0 {
		pageData.Message = "Cannot add a book. Duplicate ISBN 13"
		pageData.HasError = true
	}
	
	var countIsbn10 int64
	Db.Table("books").Where("isbn_10 = ?", book.Isbn10).Count(&countIsbn10)
	if countIsbn10 > 0 {
		pageData.Message = "Cannot add a book. Duplicate ISBN 10"
		pageData.HasError = true
	}

	if !pageData.HasError {
		transactionErr := Db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Table("books").Create(&book).Error; err != nil {
				return err
			}

			for _, v := range bookAuthorIDs.AuthorIDs {
				var bookAuthor model.BookAuthor
				bookAuthor.BookID = book.ID
				bookAuthor.AuthorID = v

				if err := tx.Table("book_authors").Create(&bookAuthor).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if transactionErr != nil {
			pageData.Message = "Cannot add a Book." 
			pageData.HasError = true
		}
	}

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/result.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func UIUpdateBookForm(c *gin.Context) {
	isbn_13 := c.Param("isbn_13")
	
	var book model.Book
	Db.Where("isbn_13 = ?", isbn_13).First(&book)

	var bookAuthors []model.BookAuthor
	Db.Where("book_id = ?", book.ID).Find(&bookAuthors)

	var authors []AuthorDisplay
	Db.Table("authors").Find(&authors)

	var publishers []model.Publisher
	Db.Find(&publishers)

	type PageData struct {
		Book model.Book
		Authors []AuthorDisplay
		Publishers []model.Publisher
	}

	var data PageData
	data.Book = book
	data.Publishers = publishers
	
	data.Authors = authors
	for i, v := range data.Authors {
		data.Authors[i].IsSelected = checkIfIDIsInExistingAuthorIDs(v.ID, bookAuthors)
    }
	
	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/update_form.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func UISubmitUpdateBookForm(c *gin.Context) {
	isbn_13 := c.Param("isbn_13")

	var book model.Book
	c.ShouldBind(&book)

	type authorIDs struct {
		AuthorIDs []uint64 `form:"author-ids[]"`
	}

	var bookAuthorIDs authorIDs
	c.ShouldBind(&bookAuthorIDs)
	
	var existingBook model.Book
	Db.Where("isbn_13 = ?", isbn_13).First(&existingBook)

	existingBook.Title = book.Title
	existingBook.PublicationYear = book.PublicationYear
	existingBook.PublisherID = book.PublisherID
	existingBook.ImageURL = book.ImageURL
	existingBook.Edition = book.Edition
	existingBook.ListPrice = book.ListPrice

	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully updated Book"
	pageData.HasError = false

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Save(&existingBook).Error; err != nil {
			return err
		}

		if err := tx.Table("book_authors").Where("book_id = ?", existingBook.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			return err
		}

		for _, v := range bookAuthorIDs.AuthorIDs {

			var bookAuthor model.BookAuthor
			bookAuthor.BookID = existingBook.ID
			bookAuthor.AuthorID = v

			if err := tx.Table("book_authors").Create(&bookAuthor).Error; err != nil {
				return err
			}
			
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot update Book." 
		pageData.HasError = true
	}

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/result.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func checkIfIDIsInExistingAuthorIDs(authorID uint64, bookAuthors []model.BookAuthor) bool {
	var isFound bool = false
	for _, v := range bookAuthors {
        if v.AuthorID == authorID {
            isFound = true
        }
    }
    return isFound
}

func UIViewBook(c *gin.Context) {
	isbn_13 := c.Param("isbn_13")

	var book BookDisplay

	Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Where("b.isbn_13 = ?", isbn_13).Group("b.id").First(&book)

	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/view_one.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func UIDeleteBook(c *gin.Context) {
	isbn_13 := c.Param("isbn_13")

	var book model.Book
	result := Db.Where("isbn_13 = ?", isbn_13).First(&book)
	
	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully deleted Book" 
	pageData.HasError = false

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/books/result.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if result.Error == gorm.ErrRecordNotFound {
		pageData.Message = "Book not found." 
		pageData.HasError = true
		if err := tmpl.Execute(w, pageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Table("book_authors").Where("book_id = ?", book.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&book).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot delete this Book." 
		pageData.HasError = true
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func GetBooks(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	keyword = strings.TrimLeft(keyword, " ") 
	keyword = strings.TrimRight(keyword, " ")
	keyword = strings.NewReplacer(`'`, `\'`, `"`, `\"`).Replace(keyword) 
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if limit < 1 {
		limit = recordLimitPerPage
	}

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()

		type BookData struct {
			ID				uint64		`json:"id"`
			Title			string		`json:"title"`
			Isbn13			string		`json:"isbn_13"`
			Isbn10			string		`json:"isbn_10"`
			PublicationYear	int16		`json:"publication_year"`
			PublisherID		uint64		`json:"publisher_id"`
			ImageURL		string		`json:"image_url"`
			Edition			string		`json:"edition"`
			ListPrice		float32		`json:"list_price"`
			AuthorIDs		string		`json:"author_ids" gorm:"column:author_ids"`
		}

		var books []BookData
		var result *gorm.DB
		
		if (len(keyword) != 0){
		
			havingString := " title LIKE '%" + keyword + "%' " 
			havingString = havingString + " or isbn_13 LIKE '%" + keyword + "%' " 
			havingString = havingString + " or isbn_10 LIKE '%" + keyword + "%' " 
			havingString = havingString + " OR author LIKE '%" + keyword + "%' "
			havingString = havingString + " OR publication_year LIKE '%" + keyword + "%' "
			havingString = havingString + " OR publisher_name LIKE '%" + keyword + "%' "
			
			result = Db.Table("books b").Select("b.id", "b.title", "CONCAT('[', GROUP_CONCAT('', a.id, ''), ']') author_ids", "b.isbn_13", "b.isbn_10", "b.publication_year", "b.publisher_id", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Group("b.id").Having(havingString).Limit(limit).Offset(page - 1).Find(&books)
		} else {
			result = Db.Table("books b").Select("b.id", "b.title", "CONCAT('[', GROUP_CONCAT('', a.id, ''), ']') author_ids", "b.isbn_13", "b.isbn_10", "b.publication_year", "b.publisher_id", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Group("b.id").Limit(limit).Offset(page - 1).Find(&books)
		}

		if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
			response := model.Response[map[string]string]{
				Message: "No Books yet / Books not found",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		
		booksDataJson, _ := json.Marshal(books)
		booksDataJsonStr := string(booksDataJson) 
		
		data := make(map[string]string)
		data["books"] = booksDataJsonStr

		response := model.Response[map[string]string]{
			Message: "Successfully retrieved books",
			Count: result.RowsAffected,
			Page: int64(page),
			Data:    data,
		}

		c.IndentedJSON(http.StatusOK, response)
		return

	}()
	
	wg.Wait()

	return
}

func GetBook(c *gin.Context) {
	Isbn13 := c.Param("isbn_13")
	
	var book model.Book
	book.Isbn13 = Isbn13

	err := Validate.Struct(book)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.StructField()] = err.Tag()
		}

		response := model.Response[map[string]string]{
			Message: "Field Errors",
			Data:    errors,
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	result := Db.Where("isbn_13 = ?", book.Isbn13).First(&book)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Book not found",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	bookDataJson, _ := json.Marshal(book)
	bookDataJsonStr := string(bookDataJson)

	data := make(map[string]string)
	data["book"] = bookDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully retrieved book",
		Count: result.RowsAffected,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)
	
	return
}

func AddBook(c *gin.Context) {
	type bookData struct {
		ID				uint64	`gorm:"primaryKey"`
		Title			string	`json:"title"`
		Isbn13			string	`gorm:"column:isbn_13" json:"isbn_13"`
		Isbn10			string	`gorm:"column:isbn_10" json:"isbn_10"`
		PublicationYear	int16 	`json:"publication_year"`
		PublisherID		uint64	`json:"publisher_id"`
		ImageURL		string	`json:"image_url"`
		Edition			string	`json:"edition"`
		ListPrice		float32	`json:"list_price"`
		AuthorIDs		[]uint64	`gorm:"-" json:"author_ids" validate:"required"`
	}

	var book bookData

	if err := c.BindJSON(&book); err != nil {
		return
	}

	err := Validate.Struct(book)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.StructField()] = err.Tag()
		}

		response := model.Response[map[string]string]{
			Message: "Field Errors",
			Data:    errors,
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	var hasIsbn bool = false
	
	if len(book.Isbn13) > 0 {
		hasIsbn = true
		var countIsbn13 int64
		Db.Table("books").Where("isbn_13 = ?", book.Isbn13).Count(&countIsbn13)
		if countIsbn13 > 0 {
			response := model.Response[map[string]string]{
				Message: "Duplicate ISBN 13",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
	}
	
	if len(book.Isbn10) > 0 {
		hasIsbn = true
		var countIsbn10 int64
		Db.Table("books").Where("isbn_10 = ?", book.Isbn10).Count(&countIsbn10)
		if countIsbn10 > 0 {
			response := model.Response[map[string]string]{
				Message: "Duplicate ISBN 10",
			}
	
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
	}

	if !hasIsbn {
		response := model.Response[map[string]string]{
			Message: "Atleast ISBN 13 or ISBN 10 is required",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if len(book.AuthorIDs) == 0 {
		response := model.Response[map[string]string]{
			Message: "Atleast one valid Author ID is required",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	authorIDsJson, _ := json.Marshal(book.AuthorIDs)
	
	var countAuthor int64
	authorIDsWhereString := string(authorIDsJson)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "[", "(", 1)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "]", ")", 1)
	
	Db.Table("authors").Where("id IN " + authorIDsWhereString).Count(&countAuthor)

	if countAuthor < int64(len(book.AuthorIDs)) {
		response := model.Response[map[string]string]{
			Message: "Author ID(s) not valid",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("books").Create(&book).Error; err != nil {
			return err
		}

		for _, v := range book.AuthorIDs {
			var bookAuthor model.BookAuthor
			bookAuthor.BookID = book.ID
			bookAuthor.AuthorID = v

			if err := tx.Table("book_authors").Create(&bookAuthor).Error; err != nil {
				return err
			}
		}

		return nil
	})

	bookDataJson, _ := json.Marshal(book)
	bookDataJsonStr := string(bookDataJson)

	data := make(map[string]string)
	data["book"] = bookDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully added an Book",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)

	return
}

func UpdateBook(c *gin.Context) {
	Isbn13 := c.Param("isbn_13")
	
	type bookData struct {
		ID				uint64	`gorm:"primaryKey"`
		Title			string	`json:"title"`
		Isbn13			string	`gorm:"column:isbn_13" json:"isbn_13"`
		Isbn10			string	`gorm:"column:isbn_10" json:"isbn_10"`
		PublicationYear	int16 	`json:"publication_year"`
		PublisherID		uint64	`json:"publisher_id"`
		ImageURL		string	`json:"image_url"`
		Edition			string	`json:"edition"`
		ListPrice		float32	`json:"list_price"`
		AuthorIDs		[]uint64	`gorm:"-" json:"author_ids" validate:"required"`
	}

	var book bookData

	if err := c.BindJSON(&book); err != nil {
		return
	}

	err := Validate.Struct(book)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.StructField()] = err.Tag()
		}

		response := model.Response[map[string]string]{
			Message: "Field Errors",
			Data:    errors,
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	var existingBook model.Book
	result := Db.Table("books").Where("isbn_13 = ?", Isbn13).First(&existingBook)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Book not found with the given ISBN-13",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if len(book.Isbn13) > 0 {
		var countIsbn13 int64
		Db.Table("books").Where("id != ? and isbn_13 = ?", existingBook.ID, book.Isbn13).Count(&countIsbn13)
		if countIsbn13 > 0 {
			response := model.Response[map[string]string]{
				Message: "Duplicate ISBN 13",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
	}
	
	if len(book.Isbn10) > 0 {
		var countIsbn10 int64
		Db.Table("books").Where("id != ? and isbn_10 = ?", existingBook.ID, book.Isbn10).Count(&countIsbn10)
		if countIsbn10 > 0 {
			response := model.Response[map[string]string]{
				Message: "Duplicate ISBN 10",
			}
	
			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
	}

	if len(book.AuthorIDs) == 0 {
		response := model.Response[map[string]string]{
			Message: "Atleast one valid Author ID is required",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	authorIDsJson, _ := json.Marshal(book.AuthorIDs)
	
	var countAuthor int64
	authorIDsWhereString := string(authorIDsJson)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "[", "(", 1)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "]", ")", 1)
	
	Db.Table("authors").Where("id IN " + authorIDsWhereString).Count(&countAuthor)

	if countAuthor < int64(len(book.AuthorIDs)) {
		response := model.Response[map[string]string]{
			Message: "Author ID(s) not valid",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if existingBook.Title != book.Title {
		existingBook.Title = book.Title
	}

	if existingBook.Isbn13 != book.Isbn13 {
		existingBook.Isbn13 = book.Isbn13
	}

	if existingBook.Isbn10 != book.Isbn10 {
		existingBook.Isbn10 = book.Isbn10
	}
	
	if existingBook.PublicationYear != book.PublicationYear {
		existingBook.PublicationYear = book.PublicationYear
	}

	if existingBook.PublisherID != book.PublisherID {
		existingBook.PublisherID = book.PublisherID
	}

	if existingBook.ImageURL != book.ImageURL {
		existingBook.ImageURL = book.ImageURL
	}
	
	if existingBook.Edition != book.Edition {
		existingBook.Edition = book.Edition
	}

	if existingBook.ListPrice != book.ListPrice {
		existingBook.ListPrice = book.ListPrice
	}

	//TODO Revalidate the data of the book to be saved

	// fmt.Println(book)
	// fmt.Println(existingBook)
	// return
	// {0 Book 1 9781891830858 1891830859 2001 1 https://asd.com Book 1 123.12 [1 2 3]}
	// {1 American Elf 9781891830853 1891830856 2004 1 https://www.collinsdictionary.com/images/full/book_181404689_1000.jpg Book 2 1000}
	
	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("book_authors").Where("book_id = ?", existingBook.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			fmt.Println(err)
			return err
		}

		if err := tx.Table("books").Save(&existingBook).Error; err != nil {
			fmt.Println(err)
			return err
		}

		for _, v := range book.AuthorIDs {
			var bookAuthor model.BookAuthor
			bookAuthor.BookID = existingBook.ID
			bookAuthor.AuthorID = v

			if err := tx.Table("book_authors").Create(&bookAuthor).Error; err != nil {
				fmt.Println(err)
				return err
			}
		}
		
		return nil
	})

	bookDataJson, _ := json.Marshal(book)
	bookDataJsonStr := string(bookDataJson)

	data := make(map[string]string)
	data["book"] = bookDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully updated the Book",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)

	return
}

func DeleteBook(c *gin.Context) {
	Isbn13 := c.Param("isbn_13")

	var book model.Book

	result := Db.Where("isbn_13 = ?", Isbn13).First(&book)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Book not found",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Table("book_authors").Where("book_id = ?", book.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			return err
		}

		if err := tx.Table("books").Where("id = ?", book.ID).Unscoped().Delete(&model.Book{}).Error; err != nil {
			response := model.Response[map[string]string]{
				Message: "Cannot delete this book.",
			}
	
			c.IndentedJSON(http.StatusBadRequest, response)
			return err
		}

		response := model.Response[map[string]string]{
			Message: "Successfully deleted book",
			Count: result.RowsAffected,
			Page: int64(1),
		}
		
		c.IndentedJSON(http.StatusOK, response)
	
		return nil
	})
	
	
	return
}