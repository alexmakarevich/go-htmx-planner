package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"

	"go-form/entities"
	"go-form/templs"
)

type Participation struct {
	gorm.Model
	EventId uint `gorm:"not null"`
	UserId  uint `gorm:"not null"`
}

type User struct {
	gorm.Model
	Name string `gorm:"not null"`
}

func newUser(Name string) *User {
	return &User{Name: Name}
}

func newParticipation(EventId uint, UserId uint) *Participation {
	return &Participation{EventId: EventId, UserId: UserId}
}

func main() {
	fmt.Println("kek")

	db, err := gorm.Open(sqlite.Open("db/dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// db.AutoMigrate(&entities.CalendarEvent{})
	// db.Where("1 = 1").Delete(&entities.CalendarEvent{})

	// res := db.Create(entities.NewCalendarEvent("Pizza", time.Now()))
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

	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", templs.Page(templs.Home()))
	})

	server.GET("/createEvent", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create Event", templs.Page(templs.CreateEvent()))
	})

	type NewEventData struct {
		Title string `form:"title"`
	}

	server.POST("/createEvent", func(c *gin.Context) {
		var newEvent NewEventData
		err := c.ShouldBind(&newEvent)
		if err != nil {
			fmt.Println(err)
			// if i want to send a negative status, I'll need to intercept in in the browser
			c.HTML(200, "", templs.Notification(templs.BadReq))
			return
		} else {
			println("succ")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}

	})

	// thots
	// server.POST("/api/v1/createEvent", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "Create Event", templs.Page(templs.CreateEvent()))
	// })

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
		c.HTML(http.StatusOK, "event", templs.Page(templs.Event(ev)))
	})

	server.Run("localhost:9999")

}
