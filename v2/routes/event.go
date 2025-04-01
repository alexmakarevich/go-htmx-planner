package routes

import (
	"context"
	"database/sql"
	"fmt"
	"go-form/sqlc/db_entities"
	protos "go-form/v2/proto-go"
	"go-form/v2/utils"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gin-gonic/gin"
)

func MakeEventRoutes(router gin.IRouter, q *db_entities.Queries, db *sql.DB) {
	router.GET("/list-events", func(c *gin.Context) {
		events, err := q.ListCalendarEventsWithOwner(c)
		fmt.Println(events)
		fmt.Println(len(events))

		if err != nil {
			fmt.Println(err.Error())
			c.Status(500)
			return
		}
		eventList := &protos.ListEventsWithOwnerResponse{Events: []*protos.EventAndOwner{}}
		for _, e := range events {
			eventList.Events = append(eventList.Events, &protos.EventAndOwner{
				Event: &protos.Event{Title: &e.Title, DateTime: timestamppb.New(e.DateTime), Id: &e.ID},
				Owner: &protos.Owner{Id: &e.OwnerID, Name: &e.OwnerName},
			})
		}

		data, err := proto.Marshal(eventList)
		if err != nil {
			fmt.Println("Error marshaling:", err)
			c.AbortWithError(500, err)
			return
		}

		c.Data(200, "application/x-protobuf", data)
	})

	router.POST("/create-event", func(c *gin.Context) {
		response := &protos.CreateEventResponse{}

		// TODO: use response for errors too
		newEvent := protos.CreateEventParams{}
		if err := utils.ParseProto(c, &newEvent); err != nil {
			c.AbortWithError(500, err)
			return
		}

		newId, err := CreateEvent(db, q, c, *newEvent.Title, newEvent.DateTime.AsTime(), newEvent.InvitedUserIds)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		response.IdOrError = &protos.CreateEventResponse_Id{Id: newId}
	})
}

func CreateEvent(db *sql.DB, q *db_entities.Queries, c context.Context, title string, dateTime time.Time, participantIds []int64) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	res, err := qtx.CreateCalendaEvent(c, db_entities.CreateCalendaEventParams{Title: title, DateTime: dateTime})
	if err != nil {
		return 0, err
	}

	for _, pId := range participantIds {
		// TODO: obviously better as a bulk op, but SQLC doesn't support them :(
		_, err := qtx.AddParticipant(c, db_entities.AddParticipantParams{EventID: res.ID, UserID: pId})
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return res.ID, nil
}
