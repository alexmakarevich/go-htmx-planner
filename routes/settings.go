package routes

import (
	"fmt"
	"go-form/sqlc/db_entities"
	templs_settings "go-form/templs/settings"

	"github.com/gin-gonic/gin"
)

func SettingsPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return RenderPage(templs_settings.Settings())
}

func DeleteSelfHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {

		session := c.MustGet("auth-context").(db_entities.GetSessionWithUserRow)

		err := q.DeleteUser(c, session.UserID)

		if err != nil {
			fmt.Println(err.Error())
			ErrorNotificationOob(c, err.Error())
			return
		}
		c.SetCookie(CookieName, "", -1, "/", Domain, CookieSecure, CookieHTTPOnly)

		c.Header("HX-Redirect", "/loginOrRegister")
		// c.Redirect(307, "/login")
		// TODO: figure the nicest way to combine redirect & notifications (special header?)
		c.String(200, "")
	}
}
