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

type publisherForm struct {
    Name string `form:"name"`
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