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

func RefreshSearchAndSelect(q *db_entities.Queries, c *gin.Context) {
	eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)
	fmt.Println(eventId)

	query := c.Query("search")
	users, err := q.SearchUsersExcludingParticipants(c, db_entities.SearchUsersExcludingParticipantsParams{Query: query, EventID: eventId})

	if err != nil {
		println("could not search")
		fmt.Println(err.Error())
		c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "could not search"))
		return
	}

	selected, err := q.GetParticipantsByEventId(c, db_entities.GetParticipantsByEventIdParams{Status: "selected", EventID: eventId})

	if err != nil {
		println("could not search")
		fmt.Println(err.Error())
		c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "could not search"))
		return
	}

	c.HTML(http.StatusOK, "", templs_event.SearchAndSelectUsers(&eventId, &users, &selected))
}

func SearchParticipantsHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		RefreshSearchAndSelect(q, c)
	}
}

func SelectParticipantHanlder(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)
		userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

		_, err := q.AddParticipant(c, db_entities.AddParticipantParams{EventID: eventId, UserID: userId, Status: "selected"})

		if err != nil {
			println("could not add participation")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationWithText(templs.BadReq, "could not add participation"))
			return
		}

		RefreshSearchAndSelect(q, c)
	}
}

func DeselectParticipantHanlder(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)
		userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

		err := q.DeleteParticipant(c, db_entities.DeleteParticipantParams{EventID: eventId, UserID: userId})

		if err != nil {
			println("could not deselect participation")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationWithText(templs.BadReq, "could not deselect"))
			return
		}

		RefreshSearchAndSelect(q, c)
	}
}

func InviteParticipantsHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)

		updatedEvents, err := q.InviteParticipants(c, eventId)

		if err != nil {
			println("could not add participation")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "could not add participation"))
			return
		}

		if len(updatedEvents) == 0 {
			fmt.Println("No users to invite")
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "No users to invite"))
			return
		}

		println("succ")
		c.Header("HX-Refresh", "true")
		c.HTML(http.StatusCreated, "", templs.NotificationOobWithText(templs.Success, "added participant"))
	}
}

func UpdateParticipantHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)
		userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
		status := c.Param("status")

		fmt.Println(eventId)
		fmt.Println(userId)
		fmt.Println(status)

		err := q.UpdateParticipant(c, db_entities.UpdateParticipantParams{EventID: eventId, UserID: userId, Status: status})

		if err != nil {
			println("could not add participation")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "could not add participation"))
		} else {
			println("succ")
			c.Header("HX-Refresh", "true")
			c.HTML(http.StatusCreated, "", templs.NotificationOobWithText(templs.Success, "added participant"))
		}
	}
}

func AddParticipantHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)
		userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
		status := c.Param("status")

		_, err := q.AddParticipant(c, db_entities.AddParticipantParams{EventID: eventId, UserID: userId, Status: status})

		if err != nil {
			println("could not add participation")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationWithText(templs.BadReq, "could not add participation"))
		} else {
			println("succ")
			c.Header("HX-Refresh", "true")
			c.HTML(http.StatusCreated, "", templs.NotificationOobWithText(templs.Success, "added participant"))
		}
	}
}

func DeleteParticipantHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		eventId, _ := strconv.ParseInt(c.Param("eventId"), 10, 64)
		userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

		err := q.DeleteParticipant(c, db_entities.DeleteParticipantParams{EventID: eventId, UserID: userId})

		if err != nil {
			println("could not remove participation")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, "could not remove participation"))
		} else {
			println("succ")
			c.Header("HX-Refresh", "true")
			c.HTML(http.StatusCreated, "", templs.NotificationOobWithText(templs.Success, "removed participant"))
		}
	}
}

func UpdateEventPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		event, err := q.GetCalendarEventWithOwner(c, id)
		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}

		maybeParicipants, err := q.ListUsersInRelationToThisEvent(c, id)

		log.Println("maybeParicipants")
		log.Println(maybeParicipants)

		if err != nil {
			log.Println("Nooooooooo")
			log.Println(err.Error())
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}

		participants := []db_entities.User{}
		selected := []db_entities.User{}
		nonParticipants := []db_entities.User{}

		for _, maybe := range maybeParicipants {
			if !maybe.EventID.Valid {
				nonParticipants = append(nonParticipants, db_entities.User{ID: maybe.ID, UserName: maybe.UserName})
			} else {
				if maybe.Status.String == "selected" {
					selected = append(selected, db_entities.User{ID: maybe.ID, UserName: maybe.UserName})
				} else {
					participants = append(participants, db_entities.User{ID: maybe.ID, UserName: maybe.UserName})
				}
			}
		}

		RenderPage(templs_event.UpdateEvent(&event, &participants, &nonParticipants, &selected))(c)
	}
}

// TODO: is this used?
func GetEventPageHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		ev, err := q.GetCalendarEvent(c, id)
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
