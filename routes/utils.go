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

func ErrorNotificationOob(c *gin.Context, text string) {
	log.Println(text)
	c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, text))
}

func GenericDbErrorHandler(c *gin.Context, err error) {
	log.Println("DB error", err)
	c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "Something went wrong in the DB"))
}

func GenericInvalidStateHandler(c *gin.Context, explanation string, errMaybe *error) {
	log.Println("Unexpected error", explanation, &errMaybe)
	c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "An unexpected error occured"))
}
