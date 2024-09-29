package routes

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"

	templs "go-form/templs/generic"
)

func SimpleRender(tc templ.Component) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", tc)
	}
}

func RenderPage(tc templ.Component) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", templs.Page(tc))
	}
}

func ErrorNotification(c *gin.Context, text string) {
	log.Println(text)
	c.HTML(200, "", templs.NotificationWithText(templs.BadReq, text))
}
