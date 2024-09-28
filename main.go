package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/a-h/templ"
	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"

	"go-form/entities"
	templs_event "go-form/templs/event"
	templs "go-form/templs/generic"
	templs_user "go-form/templs/user"
)

// TODO: timezones

func simpleRender(tc templ.Component) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", tc)
	}
}

func renderPage(tc templ.Component) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", templs.Page(tc))
	}
}

func ErrorNotification(c *gin.Context, text string) {
	log.Println(text)
	c.HTML(200, "", templs.NotificationWithText(templs.BadReq, text))
}

func main() {
	db, err := gorm.Open(sqlite.Open("db/dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&entities.CalendarEvent{})
	db.AutoMigrate(&entities.User{})

	// db.Where("1 = 1").Delete(&entities.CalendarEvent{})

	// res := db.Create(entities.NewCalendarEvent("Kepchuck", time.Now()))
	// res := db.Create(&CalendarEvent{})

	// if res.Error != nil {
	// 	log.Fatal(res.Error)
	// }

	var ev entities.CalendarEvent

	result := db.Model(&entities.CalendarEvent{Title: "Burger"}).First(&ev)

	if result.Error != nil {
		fmt.Println("ERRORINO")
		fmt.Println(result.Error)
	}

	fmt.Print(&ev)

	fmt.Println("kek+")

	server := gin.Default()
	server.Static("/public", "./public")
	// server.LoadHTMLFiles("./templs/event.html")

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	server.SetTrustedProxies(nil)

	server.GET("/", renderPage(templs.Home()))

	server.GET("/createEvent", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create Event", templs.Page(templs_event.CreateEvent()))
	})

	server.GET("/createUser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create User", templs.Page(templs_user.CreateUser()))
	})

	var events []entities.CalendarEvent
	var users []entities.User

	server.GET("/events", func(c *gin.Context) {
		db.Find(&events)
		renderPage(templs_event.EventList(&events))(c)
	})

	type NewEventData struct {
		Title    string `form:"title" binding:"required"`
		DateTime string `form:"date-time" binding:"required"`
	}

	server.POST("/htmx/createEvent", func(c *gin.Context) {
		var newEvent NewEventData
		err := c.ShouldBind(&newEvent)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		newEventTime, err := time.Parse("2006-01-02T15:04", newEvent.DateTime)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		res := db.Create(entities.NewCalendarEvent(newEvent.Title, newEventTime))
		if res.Error != nil {
			println("bad db")
			fmt.Println(res.Error)
			c.HTML(200, "", templs.Notification(templs.BadReq))
		} else {
			println("succ")
			fmt.Println(newEvent.Title)
			c.Header("HX-Redirect", "/events")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}

	})

	server.GET("/users", func(c *gin.Context) {
		db.Find(&users)
		renderPage(templs_user.UserList(&users))(c)
	})

	type NewUserData struct {
		UserName string `form:"username" binding:"required"`
	}

	server.POST("/htmx/createUser", func(c *gin.Context) {
		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		res := db.Create(entities.NewUser(newUser.UserName))
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

	})

	server.GET("/event/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Take(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		c.HTML(http.StatusOK, "event", templs.Page(templs_event.Event(ev)))
	})

	server.DELETE("htmx/deleteEvent/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Delete(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			simpleRender(templs.Notification(templs.BadReq))(c)
		} else {
			// c.Data(200, gin.MIMEHTML, nil)
			simpleRender(templs.NotificationOob(templs.Success))(c)
		}
	})

	server.GET("/updateEvent/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Take(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		renderPage(templs_event.UpdateEvent(&ev))(c)
	})

	server.GET("/updateUser/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var user entities.User
		res := db.Take(&user, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		renderPage(templs_user.UpdateUser(&user))(c)
	})

	server.PUT("/htmx/updateEvent/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		var newEvent NewEventData
		err := c.ShouldBind(&newEvent)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		newEventTime, err := time.Parse("2006-01-02T15:04", newEvent.DateTime)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		updatedEvent := entities.NewCalendarEvent(newEvent.Title, newEventTime)
		updatedEvent.ID = uint(id)

		res := db.Save(&updatedEvent)
		if res.Error != nil {
			println(res.Error.Error())
			fmt.Println(res.Error)
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, res.Error.Error()))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/events")
			c.HTML(http.StatusCreated, "", templs.NotificationOob(templs.Success))
		}
	})

	server.Run("localhost:19999")

}
