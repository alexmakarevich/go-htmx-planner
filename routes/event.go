package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go-form/sqlc/db_entities"
	templs_event "go-form/templs/event"
	templs "go-form/templs/generic"
)

func CreateEventPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create Event", templs.Page(templs_event.CreateEvent()))
	}
}

type NewEventData struct {
	Title    string `form:"title" binding:"required"`
	DateTime string `form:"date-time" binding:"required"`
	// OwnerId  int64  `form:"owner-id" binding:"required"`
}

// not actively used - for dev/testing only
func CreateEventHandler(q *db_entities.Queries) func(c *gin.Context) {
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

		// session := GetAuthContext(c)

		_, err = q.CreateCalendaEvent(c, db_entities.CreateCalendaEventParams{Title: newEvent.Title, DateTime: newEventTime, OwnerID: int64(1)})

		if err != nil {
			println("bad db")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.Notification(templs.BadReq))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/events")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}
	}
}

func ListEventsPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		events, err := q.ListCalendaEvents(c)
		if err != nil {
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, err.Error()))
			return
		}
		RenderPage(templs_event.EventList(&events))(c)
	}
}

func UpdateEventHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

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

		err = q.UpdateCalendaEvent(c, db_entities.UpdateCalendaEventParams{ID: id, Title: newEvent.Title, DateTime: newEventTime})

		if err != nil {
			println("bad db")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.Notification(templs.BadReq))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/events")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}
	}
}

func UpdateEventPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		event, err := q.GetCalendaEvent(c, id)
		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		RenderPage(templs_event.UpdateEvent(&event))(c)
	}
}

// TODO: is this used?
func GetEventPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		ev, err := q.GetCalendaEvent(c, id)
		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		c.HTML(http.StatusOK, "event", templs.Page(templs_event.Event(ev)))
	}
}

func GetEventHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

func DeleteEventHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		err := q.DeleteCalendaEvent(c, id)
		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		} else {
			SimpleRender(templs.NotificationOob(templs.Success))(c)
		}
	}
}

// func BaseHandler(q *db_entities.Queries) func(c *gin.Context) {
// 	return func(c *gin.Context) {

// 	}
// }
