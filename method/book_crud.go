package method

import(
	"fmt"
	// "encoding/json"
	"html/template"
	"net/http"
	"log"
	"strconv"
	"strings"
	"sync"
	
	"xyz-books/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type bookDisplay struct {
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

type bookUpdateDisplay struct {
	ID				uint64
	Title			string
	AuthorID		uint64
	Isbn13			string	`gorm:"column:isbn_13"`
	Isbn10			string	`gorm:"column:isbn_10"`
	PublicationYear	int16
	PublisherID		uint64
	Edition			string
	ListPrice		float32
	ImageURL		string
}

func UIBookIndex(c *gin.Context) {

	keyword := c.DefaultQuery("keyword", "")
	keyword = strings.TrimLeft(keyword, " ") 
	keyword = strings.TrimRight(keyword, " ") 
	
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	currentPageNumber := pageNumber
	pageNumber = (pageNumber - 1) * 10000
	limit := 10000

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()
		
		var books []bookDisplay
		// var countResult *gorm.DB 
		var count int64
		
		if (len(keyword) != 0){
		
			havingString := " title LIKE '%" + keyword + "%' " 
			havingString = " isbn_13 LIKE '%" + keyword + "%' " 
			havingString = " isbn_10 LIKE '%" + keyword + "%' " 
			havingString = havingString + " OR author LIKE '%" + keyword + "%' "
			havingString = havingString + " OR publication_year LIKE '%" + keyword + "%' "
			havingString = havingString + " OR publisher_name LIKE '%" + keyword + "%' "
			
			
			Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Having(havingString).Count(&count)
			Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Having(havingString).Limit(limit).Offset(pageNumber).Find(&books)
		
		} else {
			Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Count(&count)
			Db.Table("books b").Select("b.id", "b.title", "GROUP_CONCAT(' ', CONCAT(a.first_name, ' ', IFNULL(a.middle_name, ''), ' ', a.last_name)) author", "b.isbn_13", "b.isbn_10", "b.publication_year", "p.name publisher_name", "b.edition", "b.list_price", "b.image_url").Joins("INNER JOIN book_authors ba ON b.id = ba.book_id").Joins("INNER JOIN authors a ON ba.author_id = a.id").Joins("INNER JOIN publishers p ON b.publisher_id = p.id").Group("b.id").Limit(limit).Offset(pageNumber).Find(&books)
		}


		type PageData struct {
			Books []bookDisplay
			PageNumbers []int
			CurrentPage int
			IsNextEnabled bool
		}
		
		var data PageData
		data.Books = books
		
		numberOfPages := count / 10000

		var pageNumbers []int
		for i := 0; i < int(numberOfPages); i++ {
			pageNumbers = append(pageNumbers, i+1)
		}

		data.PageNumbers = pageNumbers
		data.CurrentPage = currentPageNumber
		data.IsNextEnabled = true
		if (int(numberOfPages) - currentPageNumber) < 10 {
			data.IsNextEnabled = false
		}

		fmt.Println("")
		
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

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("books").Create(&book).Error; err != nil {

			c.IndentedJSON(http.StatusOK, "Book NOT SAVED")

			return err
		}

		for _, v := range bookAuthorIDs.AuthorIDs {
			var bookAuthor model.BookAuthor
			bookAuthor.BookID = book.ID
			bookAuthor.AuthorID = v

			if err := tx.Table("book_authors").Create(&bookAuthor).Error; err != nil {

				c.IndentedJSON(http.StatusOK, "Books Author NOT SAVED")
	
				return err
			}

		}

		return nil
	})

	c.IndentedJSON(http.StatusOK, "OK")


	return
}

func UIUpdateBookForm(c *gin.Context) {
	isbn_13 := c.Param("isbn_13")
	
	var book model.Book
	Db.Where("isbn_13 = ?", isbn_13).First(&book)

	var bookAuthors []model.BookAuthor
	Db.Where("book_id = ?", book.ID).Find(&bookAuthors)

	var authors []model.Author
	Db.Find(&authors)

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

	Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Save(&existingBook).Error; err != nil {

			c.IndentedJSON(http.StatusOK, "Book Details NOT UPDATED")
			return err

		}

		if err := tx.Table("book_authors").Where("book_id = ?", book.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			
			c.IndentedJSON(http.StatusOK, "Existing Books Authors NOT DELETED")
			return err

		}

		for _, v := range bookAuthorIDs.AuthorIDs {

			var bookAuthor model.BookAuthor
			bookAuthor.BookID = book.ID
			bookAuthor.AuthorID = v

			if err := tx.Table("book_authors").Create(&bookAuthor).Error; err != nil {

				c.IndentedJSON(http.StatusOK, "Books Author NOT UPDATED")
				return err

			}
		}

		return nil
	})

	c.IndentedJSON(http.StatusOK, "OK")


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

	var book bookDisplay

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

	if result.Error == gorm.ErrRecordNotFound {
		c.IndentedJSON(http.StatusNotFound, "Book NOT found")
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Table("book_authors").Where("book_id = ?", book.ID).Unscoped().Delete(&model.BookAuthor{}).Error; err != nil {
			
			c.IndentedJSON(http.StatusNotFound, "Existing Books Authors NOT found")
			return err

		}

		if err := tx.Unscoped().Delete(&book).Error; err != nil {

			c.IndentedJSON(http.StatusNotFound, "Book NOT DELETED")
			return err

		}

		return nil
	})

	c.IndentedJSON(http.StatusOK, "OK")
}

func AddBook(c *gin.Context) {
	data := model.Book{ID: 1, Title: "Book 1"}
	response := model.Response[model.Book]{
		Message: "Successfully added a book",
		Data:    data,
	}

	c.IndentedJSON(http.StatusOK, response)
	return
}