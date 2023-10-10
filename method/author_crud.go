package method

import(
	// "fmt"
	// "encoding/json"
	"html/template"
	"net/http"
	"log"
	"strings"
	"strconv"
	"sync"
	
	"xyz-books/model"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authorForm struct {
    FirstName string `form:"first-name"`
    MiddleName string `form:"middle-name"`
    LastName string `form:"last-name"`
}

func UIAuthorIndex(c *gin.Context) {

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
		
		var authors []model.Author
		var count int64
		
		if (len(keyword) != 0){
		
			whereString := " first_name LIKE '%" + keyword + "%' " 
			whereString = whereString + " or middle_name LIKE '%" + keyword + "%' " 
			whereString = whereString + " or last_name LIKE '%" + keyword + "%' " 
			
			Db.Table("authors").Select("id", "first_name", "middle_name", "last_name").Where(whereString).Limit(limit).Offset(pageNumber).Find(&authors)
			Db.Table("authors").Select("id", "first_name", "middle_name", "last_name").Where(whereString).Count(&count)
		
		} else {
			Db.Table("authors").Select("id", "first_name", "middle_name", "last_name").Limit(limit).Offset(pageNumber).Find(&authors)
			Db.Table("authors").Select("id", "first_name", "middle_name", "last_name").Count(&count)
		}


		type PageData struct {
			Keyword string
			Authors []model.Author
			PageNumbers []int64
			CountShownPageNumber int64
			CurrentPage int64
			PreviousPageNumber int64
			NextPageNumber int64
			MaxPageNumber int64
			IsNextEnabled bool
		}
		
		var data PageData
		data.Authors = authors
		
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

		parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/authors/index.html")
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

func UIAddAuthorForm(c *gin.Context) {
	c.File(ExPath + "/templates/authors/add_form.html")

	return
}

func UISubmitAddAuthorForm(c *gin.Context) {
	var authorForm authorForm
	c.ShouldBind(&authorForm)

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("authors").Create(&authorForm).Error; err != nil {

			c.IndentedJSON(http.StatusOK, "NOT SAVED")

			return err
		}

		return nil
	})

	// c.Redirect(http.StatusFound, "/authors-ui/add-form?success=true")
	// c.IndentedJSON(http.StatusOK, "OK")

	// c.JSON(200, gin.H{"name": publisherForm.Name})
	// c.File(ExPath + "/templates/publishers/add_form.html")

	return
}

func UIUpdateAuthorForm(c *gin.Context) {
	ID := c.Param("id")
	
	var author model.Author
	Db.Where("id = ?", ID).First(&author)

	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/authors/update_form.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, author); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func UIViewAuthor(c *gin.Context) {
	ID := c.Param("id")

	var author model.Author

	Db.Where("id = ?", ID).First(&author)

	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/authors/view_one.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, author); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}