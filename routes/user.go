package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-form/sqlc/db_entities"
	templs "go-form/templs/generic"
	templs_user "go-form/templs/user"
)

func ListUsersPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		users, err := q.ListUsers(c)
		if err != nil {
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, err.Error()))
			return
		}
		RenderPage(templs_user.UserList(&users))(c)
	}
}

func UpdateUserHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		err = q.UpdateUser(c, db_entities.UpdateUserParams{ID: id, UserName: newUser.UserName, Password: newUser.Password})

		if err != nil {
			println(err.Error())
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, err.Error()))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/users")
			c.HTML(http.StatusCreated, "", templs.NotificationOob(templs.Success))
		}
	}
}

func UpdateUserPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		user, err := q.GetUser(c, id)
		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		RenderPage(templs_user.UpdateUser(&user))(c)
	}
}

func CreateUserPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		RenderPage(templs_user.CreateUser())(c)
	}
}

func DeleteUserHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		err := q.DeleteUser(c, id)
		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		SimpleRender(templs.NotificationOob(templs.Success))(c)
	}
}

// not actively used - for dev/testing only
func CreateUserHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		_, err = q.CreateUser(c, db_entities.CreateUserParams{UserName: newUser.UserName, Password: newUser.Password})

		if err != nil {
			println("bad db")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.Notification(templs.BadReq))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/users")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}
	}
}

// func BaseHandler(q *db_entities.Queries) func(c *gin.Context) {
// 	return func(c *gin.Context) {

// 	}
// }
