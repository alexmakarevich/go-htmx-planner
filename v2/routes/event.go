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

func MakeEventRoutes(router gin.IRouter, q *db_entities.Queries, db *sql.DB, globalCtx context.Context) {
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
			response.IdOrError = &protos.CreateEventResponse_ErrorMessage{ErrorMessage: err.Error()}
			utils.SendProto(c, 400, response)
			return
		}

		newId, err := CreateEvent(db, q, globalCtx, *c, *newEvent.Title, newEvent.DateTime.AsTime(), newEvent.InvitedUserIds)
		if err != nil {
			response.IdOrError = &protos.CreateEventResponse_ErrorMessage{ErrorMessage: err.Error()}
			utils.SendProto(c, 500, response)
			return
		}

		response.IdOrError = &protos.CreateEventResponse_Id{Id: newId}
		utils.SendProto(c, 200, response)
		return
	})
}

func CreateEvent(db *sql.DB, q *db_entities.Queries, globalCtx context.Context, reqCtx gin.Context, title string, dateTime time.Time, participantIds []int64) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	qtx := q.WithTx(tx)
	session := reqCtx.MustGet("auth-context").(db_entities.GetSessionWithUserRow)

	res, err := qtx.CreateCalendaEvent(globalCtx, db_entities.CreateCalendaEventParams{Title: title, DateTime: dateTime, OwnerID: session.UserID})
	if err != nil {
		return 0, err
	}

	for _, pId := range participantIds {
		// TODO: obviously better as a bulk op, but SQLC doesn't support them :(
		_, err := qtx.AddParticipant(globalCtx, db_entities.AddParticipantParams{EventID: res.ID, UserID: pId})
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
