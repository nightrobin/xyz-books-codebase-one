package method

import(
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"sync"
	
	"xyz-books/model"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


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

		RenderPage(c, "/templates/authors/index.html", data)

		return
	}()
	
	wg.Wait()

	return
}

func UIAddAuthorForm(c *gin.Context) {
	c.File(ExPath + "/templates/authors/add_form.html")
	return
}

func UISubmitAddAuthorForm(c *gin.Context) {
	var author model.Author
	c.ShouldBind(&author)
	
	var pageData model.PageData
	
	pageData.Errors = FieldValidator(author)

	if len(pageData.Errors) > 0 {
		pageData.Message = "Cannot add the author."
		RenderPage(c, "/templates/authors/result.html", pageData)
		return
	}

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("authors").Create(&author).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Cannot add the author." 
		RenderPage(c, "/templates/authors/result.html", pageData)
		return
	}


	pageData.Message = "Successfully added the Author" 

	RenderPage(c, "/templates/authors/result.html", pageData)
	
	return
}

func UIUpdateAuthorForm(c *gin.Context) {
	ID := c.Param("id")
	
	var author model.Author
	Db.Where("id = ?", ID).First(&author)

	RenderPage(c, "/templates/authors/update_form.html", author)

	return
}

func UISubmitUpdateAuthorForm(c *gin.Context) {
	ID := c.Param("id")

	var pageData model.PageData

	var author model.Author
	c.ShouldBind(&author)

	pageData.Errors = FieldValidator(author)
	if len(pageData.Errors) > 0 {
		pageData.Message = "Cannot update the author."

		RenderPage(c, "/templates/authors/result.html", pageData)

		return
	}

	var existingAuthor model.Author
	Db.Where("id = ?", ID).First(&existingAuthor)

	if existingAuthor.FirstName != author.FirstName {
		existingAuthor.FirstName = author.FirstName
	}
	
	if existingAuthor.MiddleName != author.MiddleName {
		existingAuthor.MiddleName = author.MiddleName
	}
	
	if existingAuthor.LastName != author.LastName {
		existingAuthor.LastName = author.LastName
	}

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existingAuthor).Error; err != nil {
			return err
		}

		return nil
	})

	pageData.Message = "Successfully updated the Author"

	RenderPage(c, "/templates/authors/result.html", pageData)

	return
}

func UIViewAuthor(c *gin.Context) {
	ID := c.Param("id")

	var pageData model.PageData

	var author model.Author
	c.ShouldBind(&author)
	
	result := Db.Where("id = ?", ID).First(&author)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		pageData.Message = "This Author does not exist."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ID", Message: "Invalid ID given."}}

		RenderPage(c, "/templates/authors/result.html", pageData)

		return
	}

	RenderPage(c, "/templates/authors/view_one.html", author)

	return
}

func UIDeleteAuthor(c *gin.Context) {
	ID := c.Param("id")

	var pageData model.PageData

	var author model.Author
	c.ShouldBind(&author)
	
	result := Db.Where("id = ?", ID).First(&author)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		pageData.Message = "This Author does not exist."
		pageData.Errors = []model.ApiError{model.ApiError{Param: "ID", Message: "Invalid ID given."}}

		RenderPage(c, "/templates/authors/result.html", pageData)

		return
	}

	transactionErr := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("authors").Where("id = ?", author.ID).Unscoped().Delete(&model.Author{}).Error; err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		pageData.Message = "Deletion failed"
		pageData.Errors =  []model.ApiError{model.ApiError{Param: "Author", Message: "Cannot delete this Author because it is currently used in a book."}}
		RenderPage(c, "/templates/authors/result.html", pageData)

		return
	}

	pageData.Message = "Successfully deleted author."

	RenderPage(c, "/templates/authors/result.html", pageData)

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
				Message: "No Authors yet / Authors not found with the given keyword.",
			}

			c.IndentedJSON(http.StatusBadRequest, response)
			return
		}
		
		authorDataJson, _ := json.Marshal(authors)
		authorDataJsonStr := string(authorDataJson)

		data := make(map[string]string)
		data["authors"] = authorDataJsonStr

		response := model.Response[map[string]string]{
			Message: "Successfully retrieved the authors.",
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
			Message: "Author not found with the given ID.",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	
	authorDataJson, _ := json.Marshal(author)
	authorDataJsonStr := string(authorDataJson)

	data := make(map[string]string)
	data["author"] = authorDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully retrieved the author.",
		Count: result.RowsAffected,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)
	
	return
}

func AddAuthor(c *gin.Context) {
	var author model.Author
	c.ShouldBind(&author)

	var errors []model.ApiError
	errors = FieldValidator(author)
	if errors != nil {
		response := model.Response[map[string]string]{
			Message: "Field Errors",
		}
		response.Errors = errors

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("authors").Create(&author).Error; err != nil {
			return err
		}

		return nil
	})

	authorDataJson, _ := json.Marshal(author)
	authorDataJsonStr := string(authorDataJson)

	data := make(map[string]string)
	data["author"] = authorDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully added the Author",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}

	c.IndentedJSON(http.StatusOK, response)

	return
}

func UpdateAuthor(c *gin.Context) {
	ID := c.Param("id")

	var newAuthor model.Author
	c.ShouldBind(&newAuthor)

	var existingAuthor model.Author
	result := Db.Where("id = ?", ID).First(&existingAuthor)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Author not found",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	if existingAuthor.FirstName != newAuthor.FirstName {
		existingAuthor.FirstName = newAuthor.FirstName
	}
	
	if existingAuthor.MiddleName != newAuthor.MiddleName {
		existingAuthor.MiddleName = newAuthor.MiddleName
	}
	
	if existingAuthor.LastName != newAuthor.LastName {
		existingAuthor.LastName = newAuthor.LastName
	}

	var errors []model.ApiError
	errors = FieldValidator(existingAuthor)
	if errors != nil {
		response := model.Response[map[string]string]{
			Message: "Field Errors",
		}
		response.Errors = errors

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existingAuthor).Error; err != nil {
			return err
		}

		return nil
	})

	authorDataJson, _ := json.Marshal(existingAuthor)
	authorDataJsonStr := string(authorDataJson)

	data := make(map[string]string)
	data["author"] = authorDataJsonStr

	response := model.Response[map[string]string]{
		Message: "Successfully updated the author.",
		Count: 1,
		Page: int64(1),
		Data:	data,
	}
	
	c.IndentedJSON(http.StatusOK, response)

	return
}

func DeleteAuthor(c *gin.Context) {
	ID := c.Param("id")

	var author model.Author

	result := Db.Where("id = ?", ID).First(&author)

	if result.Error == gorm.ErrRecordNotFound || result.RowsAffected == 0 {
		response := model.Response[map[string]string]{
			Message: "Author not found with the given ID.",
		}

		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Table("authors").Where("id = ?", author.ID).Unscoped().Delete(&model.Author{}).Error; err != nil {
			response := model.Response[map[string]string]{
				Message: "Cannot delete this Author because it is currently used in a book.",
			}
	
			c.IndentedJSON(http.StatusBadRequest, response)
			return err
		}

		response := model.Response[map[string]string]{
			Message: "Successfully deleted the author",
			Count: result.RowsAffected,
			Page: int64(1),
		}
		
		c.IndentedJSON(http.StatusOK, response)
	
		return nil
	})
	
	return
}