package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"go-form/entities"
	templs_event "go-form/templs/event"
	templs "go-form/templs/generic"
)

type NewEventData struct {
	Title    string `form:"title" binding:"required"`
	DateTime string `form:"date-time" binding:"required"`
}

func ListEventsPageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var events []entities.CalendarEvent
		db.Find(&events)
		RenderPage(templs_event.EventList(&events))(c)
	}
}

func UpdateEventHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func UpdateEventPageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Take(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		RenderPage(templs_event.UpdateEvent(&ev))(c)
	}
}

func GetEventPageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func GetEventHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func DeleteEventHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Delete(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			SimpleRender(templs.Notification(templs.BadReq))(c)
		} else {
			SimpleRender(templs.NotificationOob(templs.Success))(c)
		}
	}
}

func CreateEventPageHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create Event", templs.Page(templs_event.CreateEvent()))
	}
}

// not actively used - for dev/testing only
func CreateEventHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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

	}
}

// func BaseHandler(db *gorm.DB) func(c *gin.Context) {
// 	return func(c *gin.Context) {

// 	}
// }
