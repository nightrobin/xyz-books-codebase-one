package method

import(
	// "fmt"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	
	"xyz-books/model"
	
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully added a Publisher"
	pageData.HasError = false

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("publishers").Create(&publisherForm).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot add a Publisher." 
		pageData.HasError = true
	}

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/publishers/result.html")
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

func UISubmitUpdatePublisherForm(c *gin.Context) {
	ID := c.Param("id")

	var publisher model.Publisher
	c.ShouldBind(&publisher)

	var existingPublisher model.Publisher
	Db.Where("id = ?", ID).First(&existingPublisher)

	existingPublisher.Name = publisher.Name

	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully updated Publisher"
	pageData.HasError = false

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Save(&existingPublisher).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot update Publisher." 
		pageData.HasError = true
	}

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/publishers/result.html")
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

func UIDeletePublisher(c *gin.Context) {
	ID := c.Param("id")

	var publisher model.Publisher
	result := Db.Where("id = ?", ID).First(&publisher)
	
	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully deleted Publisher" 
	pageData.HasError = false

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/publishers/result.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if result.Error == gorm.ErrRecordNotFound {
		pageData.Message = "Publisher not found." 
		pageData.HasError = true
		if err := tmpl.Execute(w, pageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Table("publishers").Where("id = ?", publisher.ID).Unscoped().Delete(&model.Publisher{}).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot delete this Publisher because it is currently used in a book." 
		pageData.HasError = true
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func GetPublishers(c *gin.Context) {
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
		
		var publishers []model.Publisher
		var result *gorm.DB

		if (len(keyword) != 0){
		
			whereString := " name LIKE '%" + keyword + "%' " 
			
			result = Db.Table("publishers").Select("id", "name").Where(whereString).Limit(limit).Offset((page - 1 )* recordLimitPerPage).Find(&publishers)
		
		} else {
			result = Db.Table("publishers").Select("id", "name").Limit(limit).Offset((page - 1 )* recordLimitPerPage).Find(&publishers)
		}

		if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
			response := model.Response[map[string]string]{
				Message: "No Publishers yet / Publishers not found",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		
		publisherDataJson, _ := json.Marshal(publishers)
		publisherDataJsonStr := string(publisherDataJson)

		data := make(map[string]string)
		data["publishers"] = publisherDataJsonStr

		response := model.Response[map[string]string]{
			Message: "Successfully retrieved publishers",
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

func GetPublisher(c *gin.Context) {
	ID := c.Param("id")

	var publisher model.Publisher

	result := Db.Where("id = ?", ID).First(&publisher)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Publisher not found",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	publisherDataJson, _ := json.Marshal(publisher)
	publisherDataJsonStr := string(publisherDataJson)

	data := make(map[string]string)
	data["publisher"] = publisherDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully retrieved publisher",
		Count: result.RowsAffected,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)
	
	return
}

func AddPublisher(c *gin.Context) {
	var publisher model.Publisher

	if err := c.BindJSON(&publisher); err != nil {
		return
	}

	err := Validate.Struct(publisher)
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

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("publishers").Create(&publisher).Error; err != nil {
			return err
		}

		return nil
	})

	publisherDataJson, _ := json.Marshal(publisher)
	publisherDataJsonStr := string(publisherDataJson)

	data := make(map[string]string)
	data["publisher"] = publisherDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully added an Publisher",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)

	return
}

func DeletePublisher(c *gin.Context) {
	ID := c.Param("id")

	var publisher model.Publisher

	result := Db.Where("id = ?", ID).First(&publisher)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Publisher not found",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Table("publishers").Where("id = ?", publisher.ID).Unscoped().Delete(&model.Publisher{}).Error; err != nil {
			response := model.Response[map[string]string]{
				Message: "Cannot delete this Publisher because it is currently used in a book.",
			}
	
			c.IndentedJSON(http.StatusBadRequest, response)
			return err
		}

		response := model.Response[map[string]string]{
			Message: "Successfully deleted publisher",
			Count: result.RowsAffected,
			Page: int64(1),
		}
		
		c.IndentedJSON(http.StatusOK, response)
	
		return nil
	})
	
	
	return
}