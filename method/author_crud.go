package method

import(
	// "fmt"
	// "encoding/json"
	// "html/template"
	"net/http"
	// "log"
	// "sync"
	
	// "xyz-books/model"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authorForm struct {
    FirstName string `form:"first-name"`
    MiddleName string `form:"middle-name"`
    LastName string `form:"last-name"`
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

	c.IndentedJSON(http.StatusOK, "OK")

	// c.JSON(200, gin.H{"name": publisherForm.Name})
	// c.File(ExPath + "/templates/publishers/add_form.html")

	return
}