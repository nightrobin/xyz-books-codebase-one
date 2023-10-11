package method

import(
	// "fmt"
	"encoding/json"
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

	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully added an Author" 
	pageData.HasError = false

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("authors").Create(&authorForm).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot add an Author." 
		pageData.HasError = true
	}

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/authors/result.html")
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

func UISubmitUpdateAuthorForm(c *gin.Context) {
	ID := c.Param("id")

	var author model.Author
	c.ShouldBind(&author)

	var existingAuthor model.Author
	Db.Where("id = ?", ID).First(&existingAuthor)

	existingAuthor.FirstName = author.FirstName
	existingAuthor.MiddleName = author.MiddleName
	existingAuthor.LastName = author.LastName

	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully updated Author"
	pageData.HasError = false

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Save(&existingAuthor).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot update Author." 
		pageData.HasError = true
	}

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/authors/result.html")
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

func UIDeleteAuthor(c *gin.Context) {
	ID := c.Param("id")

	var author model.Author
	result := Db.Where("id = ?", ID).First(&author)
	
	type PageData struct {
		Message string
		HasError bool
	}

	var pageData PageData
	pageData.Message = "Successfully deleted Author" 
	pageData.HasError = false

	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + "/templates/authors/result.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(parsedIndexTemplate, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if result.Error == gorm.ErrRecordNotFound {
		pageData.Message = "Author not found." 
		pageData.HasError = true
		if err := tmpl.Execute(w, pageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
	
		if err := tx.Table("authors").Where("id = ?", author.ID).Unscoped().Delete(&model.Author{}).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot delete this Author because it is currently used in a book." 
		pageData.HasError = true
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func GetAuthors(c *gin.Context) {
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
		
		var authors []model.Author
		var result *gorm.DB

		if (len(keyword) != 0){
		
			whereString := " first_name LIKE '%" + keyword + "%' " 
			whereString = whereString + " or middle_name LIKE '%" + keyword + "%' " 
			whereString = whereString + " or last_name LIKE '%" + keyword + "%' " 
			
			result = Db.Table("authors").Select("id", "first_name", "middle_name", "last_name").Where(whereString).Limit(limit).Offset((page - 1 )* recordLimitPerPage).Find(&authors)
		
		} else {
			result = Db.Table("authors").Select("id", "first_name", "middle_name", "last_name").Limit(limit).Offset((page - 1 )* recordLimitPerPage).Find(&authors)
		}

		if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
			response := model.Response[map[string]string]{
				Message: "No Authors yet / Authors not found",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		
		authorDataJson, _ := json.Marshal(authors)
		authorDataJsonStr := string(authorDataJson)

		data := make(map[string]string)
		data["authors"] = authorDataJsonStr
		response := model.Response[map[string]string]{
			Message: "Successfully retrieved authors",
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

func GetAuthor(c *gin.Context) {
	ID := c.Param("id")

	var author model.Author

	result := Db.Where("id = ?", ID).First(&author)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Author not found",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	authorDataJson, _ := json.Marshal(author)
	authorDataJsonStr := string(authorDataJson)

	data := make(map[string]string)
	data["author"] = authorDataJsonStr
	response := model.Response[map[string]string]{
		Message: "Successfully retrieved author",
		Count: result.RowsAffected,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)
	return

}