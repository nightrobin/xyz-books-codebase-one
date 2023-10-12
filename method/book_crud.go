package method

import(
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	
	"xyz-books-codebase-one/model"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DisplayBook struct {
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

type authorIDs struct {
	AuthorIDs []uint64 `form:"AuthorIDs[]"`
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
		
		var books []DisplayBook
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
			Books []DisplayBook
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

		RenderPage(c, "/templates/books/index.html", data)
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
	
	RenderPage(c, "/templates/books/add_form.html", data)

	return
}

func UISubmitAddBookForm(c *gin.Context) {
	var book model.Book
	c.ShouldBind(&book)

	var bookAuthorIDs authorIDs
	c.ShouldBind(&bookAuthorIDs)

	book.AuthorIDs = bookAuthorIDs.AuthorIDs

	var pageData model.PageData
	
	pageData.Errors = FieldValidator(book)

	if len(pageData.Errors) > 0 {
		pageData.Message = "Cannot add the book."
		RenderPage(c, "/templates/books/result.html", pageData)
		return
	}

	var hasIsbn bool = false

	var countIsbn13 int64
	if len(book.Isbn13) > 0 {
		if !IsbnValidator(book.Isbn13) {
			pageData.Message = "Cannot add the book."
			pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN", Message: "The ISBN 13 is invalid."}}
			RenderPage(c, "/templates/books/result.html", pageData)
		
			return
		}
		
		hasIsbn = true
		Db.Table("books").Where("isbn_13 = ?", book.Isbn13).Count(&countIsbn13)
	}

	var countIsbn10 int64
	if len(book.Isbn10) > 0 {
		if !IsbnValidator(book.Isbn10) {
			pageData.Message = "Cannot add the book."
			pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN", Message: "The ISBN 10 is invalid."}}
			RenderPage(c, "/templates/books/result.html", pageData)
		
			return
		}

		hasIsbn = true

		Db.Table("books").Where("isbn_10 = ?", book.Isbn10).Count(&countIsbn10)
	}

	if !hasIsbn {
		pageData.Message = "Cannot add the book."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN", Message: "Atleast ISBN 10 or ISBN 13 must be inputted."}}
		RenderPage(c, "/templates/books/result.html", pageData)
		
		return
	}
	
	if countIsbn13 > 0 {
		pageData.Message = "Cannot add the book."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "Isbn13", Message: "ISBN 13 has already been used in other book."}}
		RenderPage(c, "/templates/books/result.html", pageData)
		
		return
	}

	if countIsbn10 > 0 {
		pageData.Message = "Cannot add the book."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "Isbn10", Message: "ISBN 10 has already been used in other book."}}
		RenderPage(c, "/templates/books/result.html", pageData)
		
		return
	}

	
	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
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
	
	if transactionErr != nil {
		pageData.Message = "Cannot add the book." 
		RenderPage(c, "/templates/books/result.html", pageData)
		return
	}

	pageData.Message = "Successfully added the Book" 

	RenderPage(c, "/templates/books/result.html", pageData)
	
	return
}

func UIUpdateBookForm(c *gin.Context) {
	ID := c.Param("id")
	
	var book model.Book
	Db.Where("id = ?", ID).First(&book)

	var bookAuthors []model.BookAuthor
	Db.Where("book_id = ?", book.ID).Find(&bookAuthors)

	var authors []model.Author
	Db.Table("authors").Find(&authors)

	var publishers []model.Publisher
	Db.Find(&publishers)

	type PageData struct {
		Book model.Book
		Authors []model.Author
		Publishers []model.Publisher
	}

	var data PageData
	data.Book = book
	data.Publishers = publishers
	
	data.Authors = authors
	for i, v := range data.Authors {
		data.Authors[i].IsSelected = checkIfIDIsInExistingAuthorIDs(v.ID, bookAuthors)
    }
	
	RenderPage(c, "/templates/books/update_form.html", data)

	return
}

func UISubmitUpdateBookForm(c *gin.Context) {
	ID := c.Param("id")

	var book model.Book
	c.ShouldBind(&book)

	var bookAuthorIDs authorIDs
	c.ShouldBind(&bookAuthorIDs)

	book.AuthorIDs = bookAuthorIDs.AuthorIDs
	
	var pageData model.PageData
		
	pageData.Errors = FieldValidator(book)

	if len(pageData.Errors) > 0 {
		pageData.Message = "Cannot update the book."
		RenderPage(c, "/templates/books/result.html", pageData)
		return
	}

	var existingBook model.Book
	result := Db.Where("id = ?", ID).First(&existingBook)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		pageData.Message = "Cannot update the book."
		
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN 13", Message: "Book not found with the given ISBN 13."}}
		
		RenderPage(c, "/templates/books/result.html", pageData)
		
		return
	}

	var hasIsbn bool = false

	var countIsbn13 int64
	if len(book.Isbn13) > 0 {
		if !IsbnValidator(book.Isbn13) {
			pageData.Message = "Cannot update the book."
			pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN", Message: "The ISBN 13 is invalid."}}
			RenderPage(c, "/templates/books/result.html", pageData)
		
			return
		}

		hasIsbn = true
		Db.Table("books").Where("id != ? and isbn_13 = ?", existingBook.ID, book.Isbn13).Count(&countIsbn13)
	}

	var countIsbn10 int64
	if len(book.Isbn10) > 0 {
		if !IsbnValidator(book.Isbn10) {
			pageData.Message = "Cannot update the book."
			pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN", Message: "The ISBN 10 is invalid."}}
			RenderPage(c, "/templates/books/result.html", pageData)
		
			return
		}

		hasIsbn = true

		Db.Table("books").Where("id != ? and isbn_10 = ?", existingBook.ID, book.Isbn10).Count(&countIsbn10)
	}

	if !hasIsbn {
		pageData.Message = "Cannot update the book."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN", Message: "Atleast ISBN 10 or ISBN 13 must be inputted."}}
		RenderPage(c, "/templates/books/result.html", pageData)
		
		return
	}
	
	if countIsbn13 > 0 {
		pageData.Message = "Cannot update the book."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "Isbn13", Message: "ISBN 13 has already been used in other book."}}
		RenderPage(c, "/templates/books/result.html", pageData)
		
		return
	}

	if countIsbn10 > 0 {
		pageData.Message = "Cannot update the book."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "Isbn10", Message: "ISBN 10 has already been used in other book."}}
		RenderPage(c, "/templates/books/result.html", pageData)
		
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

	Db.Transaction(func(tx *gorm.DB) error {
	
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

	pageData.Message = "Successfully updated the book."
		
	RenderPage(c, "/templates/books/result.html", pageData)
	
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
	Isbn13 := c.Param("isbn_13")

	var book DisplayBook
	result := Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Where("b.isbn_13 = ?", Isbn13).Group("b.id").First(&book)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		var pageData model.PageData
		pageData.Message = "This Book does not exist."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN 13", Message: "Invalid ISBN 13 given."}}

		RenderPage(c, "/templates/books/result.html", pageData)

		return
	}

	RenderPage(c, "/templates/books/view_one.html", book)

	return
}

func UIDeleteBook(c *gin.Context) {
	Isbn13 := c.Param("isbn_13")
	
	var pageData model.PageData

	var book model.Book
	result := Db.Where("isbn_13 = ?", Isbn13).First(&book)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		pageData.Message = "This Book does not exist."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ISBN 13", Message: "Invalid ISBN 13 given."}}

		RenderPage(c, "/templates/books/result.html", pageData)

		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Table("book_authors").Where("book_id = ?", book.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&book).Error; err != nil {
			return err
		}

		return nil
	})

	pageData.Message = "Successfully deleted the book."

	RenderPage(c, "/templates/books/result.html", pageData)

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

		var books []model.ApiBookData
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
				Message: "No Books yet / Books not found.",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}

		response := model.Response[[]model.ApiBookData]{
			Message: "Successfully retrieved the books.",
			Count: result.RowsAffected,
			Page: int64(page),
			Data:    books,
		}

		c.IndentedJSON(http.StatusOK, response)
		return

	}()
	
	wg.Wait()

	return
}

func GetBook(c *gin.Context) {
	Isbn13 := c.Param("isbn_13")
	
	if !IsbnValidator(Isbn13) {
		response := model.Response[map[string]string]{
			Message: "Invalid ISBN 13",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	var book model.ApiBookData

	result := Db.Table("books b").Select("b.id", "b.title", "CONCAT('[', GROUP_CONCAT('', a.id, ''), ']') author_ids", "b.isbn_13", "b.isbn_10", "b.publication_year", "b.publisher_id", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Where("b.isbn_13 = ?", Isbn13).Group("b.id").First(&book)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Book not found with the given ISBN 13.",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	bookDataJson, _ := json.Marshal(book)
	bookDataJsonStr := string(bookDataJson)

	data := make(map[string]string)
	data["book"] = bookDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully retrieved the book.",
		Count: result.RowsAffected,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)
	
	return
}

func AddBook(c *gin.Context) {
	var book model.Book

	if err := c.BindJSON(&book); err != nil {
		return
	}

	var errors []model.ApiError

	var hasIsbn bool = false
	
	if len(book.Isbn13) > 0 {
		if !IsbnValidator(book.Isbn13) {
			errors = append(errors, model.ApiError{Param:"Isbn13", Message: "Invalid ISBN 13"})
		}

		hasIsbn = true
		var countIsbn13 int64
		Db.Table("books").Where("isbn_13 = ?", book.Isbn13).Count(&countIsbn13)
		if countIsbn13 > 0 {
			errors = append(errors, model.ApiError{Param:"Isbn13", Message: "Duplicate ISBN 13"})
		}
	}
	
	if len(book.Isbn10) > 0 {
		if !IsbnValidator(book.Isbn10) {
			errors = append(errors, model.ApiError{Param:"Isbn10", Message: "Invalid ISBN 10"})
		}

		hasIsbn = true
		var countIsbn10 int64
		Db.Table("books").Where("isbn_10 = ?", book.Isbn10).Count(&countIsbn10)
		if countIsbn10 > 0 {
			errors = append(errors, model.ApiError{Param:"Isbn10", Message: "Duplicate ISBN 10"})
		}
	}

	if !hasIsbn {
		errors = append(errors, model.ApiError{Param:"ISBN", Message: "Atleast ISBN 13 or ISBN 10 is required."})
	}
	
	authorIDsJson, _ := json.Marshal(book.AuthorIDs)
	
	var countAuthor int64
	authorIDsWhereString := string(authorIDsJson)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "[", "(", 1)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "]", ")", 1)
	
	Db.Table("authors").Where("id IN " + authorIDsWhereString).Count(&countAuthor)

	if countAuthor < int64(len(book.AuthorIDs)) {
		errors = append(errors, model.ApiError{Param:"AuthorIDs", Message: "Author ID(s) not valid."})
	}

	errors = append(errors, FieldValidator(book)...)

	if errors != nil {
		response := model.Response[map[string]string]{
			Message: "Field Errors",
		}
		response.Errors = errors

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
		Message: "Successfully added the book.",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)

	return
}

func UpdateBook(c *gin.Context) {
	ID := c.Param("id")
	
	var errors []model.ApiError

	var book model.Book

	if err := c.BindJSON(&book); err != nil {
		return
	}

	var existingBook model.Book
	result := Db.Table("books").Where("id = ?", ID).First(&existingBook)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Book not found with the given ISBN 13.",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	var hasIsbn bool = false
	
	if len(book.Isbn13) > 0 {
		if !IsbnValidator(book.Isbn13) {
			errors = append(errors, model.ApiError{Param:"Isbn13", Message: "Invalid ISBN 13"})
		}

		hasIsbn = true
		var countIsbn13 int64
		Db.Table("books").Where("id != ? and isbn_13 = ?", existingBook.ID, book.Isbn13).Count(&countIsbn13)
		if countIsbn13 > 0 {
			errors = append(errors, model.ApiError{Param:"Isbn13", Message: "Duplicate ISBN 13"})
		}
	}
	
	if len(book.Isbn10) > 0 {
		if !IsbnValidator(book.Isbn10) {
			errors = append(errors, model.ApiError{Param:"Isbn10", Message: "Invalid ISBN 10"})
		}

		hasIsbn = true
		var countIsbn10 int64
		Db.Table("books").Where("id != ? and isbn_10 = ?", existingBook.ID, book.Isbn10).Count(&countIsbn10)
		if countIsbn10 > 0 {
			errors = append(errors, model.ApiError{Param:"Isbn10", Message: "Duplicate ISBN 10"})
		}
	}

	if !hasIsbn {
		errors = append(errors, model.ApiError{Param:"ISBN", Message: "Atleast ISBN 13 or ISBN 10 is required."})
	}
	
	if len(book.AuthorIDs) == 0 {
		errors = append(errors, model.ApiError{Param:"AuthorIDs", Message: "Atleast one valid Author ID is required."})
	}
	
	authorIDsJson, _ := json.Marshal(book.AuthorIDs)
	
	var countAuthor int64
	authorIDsWhereString := string(authorIDsJson)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "[", "(", 1)
	authorIDsWhereString = strings.Replace(authorIDsWhereString, "]", ")", 1)
	
	Db.Table("authors").Where("id IN " + authorIDsWhereString).Count(&countAuthor)

	if countAuthor < int64(len(book.AuthorIDs)) {
		errors = append(errors, model.ApiError{Param:"AuthorIDs", Message: "Author ID(s) not valid."})
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
	
	existingBook.AuthorIDs = book.AuthorIDs

	errors = append(errors, FieldValidator(existingBook)...)

	if errors != nil {
		response := model.Response[map[string]string]{
			Message: "Field Errors",
		}
		response.Errors = errors

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("book_authors").Where("book_id = ?", existingBook.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			return err
		}

		if err := tx.Table("books").Save(&existingBook).Error; err != nil {
			return err
		}

		for _, v := range book.AuthorIDs {
			var bookAuthor model.BookAuthor
			bookAuthor.BookID = existingBook.ID
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
		Message: "Successfully updated the book.",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)

	return
}

func DeleteBook(c *gin.Context) {
	Isbn13 := c.Param("isbn_13")

	var errors []model.ApiError

	if !IsbnValidator(Isbn13) {
		response := model.Response[map[string]string]{
			Message: "Invalid ISBN 13",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	var book model.Book

	result := Db.Where("isbn_13 = ?", Isbn13).First(&book)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Book not found with the given ISBN 13.",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if errors != nil {
		response := model.Response[map[string]string]{
			Message: "Field Errors",
		}
		response.Errors = errors

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Table("book_authors").Where("book_id = ?", book.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			return err
		}

		if err := tx.Table("books").Where("id = ?", book.ID).Unscoped().Delete(&model.Book{}).Error; err != nil {
			return err
		}

		response := model.Response[map[string]string]{
			Message: "Successfully deleted the book.",
			Count: result.RowsAffected,
			Page: int64(1),
		}
		
		c.IndentedJSON(http.StatusOK, response)
	
		return nil
	})
	
	return
}