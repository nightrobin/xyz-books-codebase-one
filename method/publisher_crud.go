package method

import(
	// "fmt"
	// "encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	
	"xyz-books/model"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type publisherForm struct {
    Name string `form:"name"`
}

func UIPublisherIndex(c *gin.Context) {
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
		
		var publishers []model.Publisher
		var count int64
		
		if (len(keyword) != 0){
		
			whereString := " name LIKE '%" + keyword + "%' " 

			Db.Table("publishers").Select("id", "name").Where(whereString).Limit(limit).Offset(pageNumber).Find(&publishers)
			Db.Table("publishers").Select("id", "name").Where(whereString).Count(&count)
		
		} else {
			Db.Table("publishers").Select("id", "name").Limit(limit).Offset(pageNumber).Find(&publishers)
			Db.Table("publishers").Select("id", "name").Count(&count)
		}


		type PageData struct {
			Keyword string
			Publishers []model.Publisher
			PageNumbers []int64
			CountShownPageNumber int64
			CurrentPage int64
			PreviousPageNumber int64
			NextPageNumber int64
			MaxPageNumber int64
			IsNextEnabled bool
		}
		
		var data PageData
		data.Publishers = publishers
		
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

		parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/publishers/index.html")
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

func UIAddPublisherForm(c *gin.Context) {
	c.File(ExPath + "/templates/publishers/add_form.html")

	return
}

func UISubmitAddPublisherForm(c *gin.Context) {
	var publisherForm publisherForm
	c.ShouldBind(&publisherForm)

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("publishers").Create(&publisherForm).Error; err != nil {

			c.IndentedJSON(http.StatusOK, "NOT SAVED")

			return err
		}

		return nil
	})

	c.IndentedJSON(http.StatusOK, "OK")

	// c.JSON(200, gin.H{"name": publisherForm.Name})
	// c.File(ExPath + "/templates/publishers/add_form.html")

	return
}

func UIUpdatePublisherForm(c *gin.Context) {
	ID := c.Param("id")
	
	var publisher model.Publisher
	Db.Where("id = ?", ID).First(&publisher)

	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/publishers/update_form.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, publisher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func UIViewPublisher(c *gin.Context) {
	ID := c.Param("id")

	var publisher model.Publisher

	Db.Where("id = ?", ID).First(&publisher)

	w := c.Writer

	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/publishers/view_one.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, publisher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}