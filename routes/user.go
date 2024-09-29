package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"go-form/entities"
	templs "go-form/templs/generic"
	templs_user "go-form/templs/user"
)

func ListUsersHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var users []entities.User
		db.Find(&users)
		RenderPage(templs_user.UserList(&users))(c)
	}
}

func UpdateUserHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		updatedUser := entities.NewUser(newUser.UserName, newUser.Password)
		updatedUser.ID = uint(id)

		res := db.Save(&updatedUser)
		if res.Error != nil {
			println(res.Error.Error())
			fmt.Println(res.Error)
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, res.Error.Error()))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/users")
			c.HTML(http.StatusCreated, "", templs.NotificationOob(templs.Success))
		}
	}
}

func UpdateUserPageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var user entities.User
		res := db.Take(&user, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		RenderPage(templs_user.UpdateUser(&user))(c)
	}
}

func DeleteUserHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var user entities.User
		res := db.Delete(&user, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			SimpleRender(templs.Notification(templs.BadReq))(c)
		} else {
			SimpleRender(templs.NotificationOob(templs.Success))(c)
		}
	}
}

// not actively used - for dev/testing only
func CreateUserHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		res := db.Create(entities.NewUser(newUser.UserName, newUser.Password))
		if res.Error != nil {
			println("bad db")
			fmt.Println(res.Error)
			c.HTML(200, "", templs.Notification(templs.BadReq))
		} else {
			println("succ")
			fmt.Println(newUser.UserName)
			c.Header("HX-Redirect", "/users")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}
	}
}

// func BaseHandler(db *gorm.DB) func(c *gin.Context) {
// 	return func(c *gin.Context) {

// 	}
// }
