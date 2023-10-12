package method

import (
	"html/template"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func RenderPage(c *gin.Context, templatePath string, pageData interface{}) {
	w := c.Writer
	parsedIndexTemplate, err := template.ParseFiles(ExPath + templatePath)

	tmpl := template.Must(parsedIndexTemplate, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}