package method

import(
	// "fmt"
	// "encoding/json"
	"net/http"
	"strconv"
	"sync"
	
	"xyz-books/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UIAddRecords(c *gin.Context) {
	numberOfRecords, _ := strconv.ParseUint(c.DefaultQuery("number", "100000"), 10, 64)

	var wg sync.WaitGroup
	
	wg.Add(1)

	go func () {
		defer wg.Done()

		Db.Transaction(func(tx *gorm.DB) error {

			for  i := uint64(0); i < numberOfRecords; i++ {
				var book model.Book
				book.ID = i+1
				book.Title = "Book " + strconv.FormatUint(book.ID, 10)
				book.Isbn13 = strconv.FormatUint(1000000000000 + i, 10)
				book.Isbn10 = strconv.FormatUint(1000000000 + i, 10)
				book.PublicationYear = 2001
				book.PublisherID = 3
				book.ImageURL = "https://www.collinsdictionary.com/images/full/book_181404689_1000.jpg"
				book.Edition = "Book " + book.Isbn10
				book.ListPrice = float32(1000 + i)

				if err := tx.Table("books").Create(&book).Error; err != nil {
	
					c.IndentedJSON(http.StatusOK, "Book record NOT SAVED")
		
					return err
				}
			}

			for  i := uint64(0); i < numberOfRecords; i++ {
				if i % 2 == 0 {
					var bookAuthor1 model.BookAuthor
					bookAuthor1.BookID = i+1
					bookAuthor1.AuthorID = 1

					if err := tx.Table("book_authors").Create(&bookAuthor1).Error; err != nil {
			
						c.IndentedJSON(http.StatusOK, "Book Author 1 record NOT SAVED")
			
						return err
					}

					var bookAuthor2 model.BookAuthor
					bookAuthor2.BookID = i+1
					bookAuthor2.AuthorID = 2

					if err := tx.Table("book_authors").Create(&bookAuthor2).Error; err != nil {
			
						c.IndentedJSON(http.StatusOK, "Book Author 2 record NOT SAVED")
			
						return err
					}
				} else { 
					var bookAuthor3 model.BookAuthor
					bookAuthor3.BookID = i+1
					bookAuthor3.AuthorID = 3

					if err := tx.Table("book_authors").Create(&bookAuthor3).Error; err != nil {
			
						c.IndentedJSON(http.StatusOK, "Book Author 3 record NOT SAVED")
			
						return err
					}
				}
			}

			c.IndentedJSON(http.StatusOK, "Book records SAVED")
			return nil
		})
	}()

	wg.Wait()


	return
}