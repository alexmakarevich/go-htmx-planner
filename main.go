package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"

	db_entities "go-form/sqlc/db_entities"

	routesV1 "go-form/v1/routes"
	templs_auth "go-form/v1/templs/auth"
	templs "go-form/v1/templs/generic"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlc/schema.sql
var ddl string

// TODO: timezones

func main() {
	// TODO: FUTURE: prod.db, migrations, backups, etc.
	db, err := sql.Open("sqlite3", "db/dev.db")

	if err != nil {
		panic("failed to connect to database")
	}

	// create tables
	// FYI: doesn't auto-migrate existing tables
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println(err.Error())
		panic("failed to exec context")
	}

	queries := db_entities.New(db)

	server := gin.Default()

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	server.SetTrustedProxies(nil)

	v2 := server.Group("/v2")

	// v2.StaticFile("/", "./svelte/index.html")
	v2.Static("/", "./svelte/dist")

	v1 := server.Group("/v1")
	v1.Static("/public", "./v1/public")
	v1.POST("/htmx/register", routesV1.RegisterHandler(queries))
	v1.POST("/htmx/login", routesV1.LoginHandler(queries))

	v1.GET("/register", routesV1.SimpleRender(templs_auth.RegisterPage()))
	v1.GET("/login", routesV1.SimpleRender(templs_auth.LoginPage()))
	v1.GET("/loginOrRegister", routesV1.SimpleRender(templs_auth.LoginOrRegister()))

	v1.Use(routesV1.AuthMiddleware(queries))
	{
		v1.GET("/", routesV1.RenderPage(templs.Home()))

		v1.POST("/htmx/logout", routesV1.LogoutHandler(queries))
		v1.GET("/settings", routesV1.SettingsPageHandler(queries))
		v1.DELETE("/htmx/delete-self", routesV1.DeleteSelfHandler(queries))

		// USER
		v1.GET("/createUser", routesV1.CreateUserPageHandler(queries))
		v1.POST("/htmx/createUser", routesV1.CreateUserHandler(queries))

		v1.GET("/users", routesV1.ListUsersPageHandler(queries))

		v1.GET("/updateUser/:id", routesV1.UpdateUserPageHandler(queries))
		v1.PUT("/htmx/updateUser/:id", routesV1.UpdateUserHandler(queries))

		v1.DELETE("/htmx/deleteUser/:id", routesV1.DeleteUserHandler(queries))

		// EVENT
		v1.GET("/createEvent", routesV1.CreateEventPageHandler(queries))
		v1.POST("/htmx/createEvent", routesV1.CreateEventHandler(queries))

		v1.GET("/events", routesV1.ListEventsPageHandler(queries))
		v1.GET("/event/:id", routesV1.ViewOrUpdateEventPageHandler(queries, false))
		v1.GET("/event/:id/invite", routesV1.ViewOrUpdateEventPageHandler(queries, false))

		v1.GET("/myInvites", routesV1.ListInvitesPagehandler(queries))

		v1.GET("/updateEvent/:id", routesV1.ViewOrUpdateEventPageHandler(queries, true))

		v1.PUT("/htmx/updateEvent/:id", routesV1.UpdateEventHandler(queries))

		v1.DELETE("/htmx/deleteEvent/:id", routesV1.DeleteEventHandler(queries))

		// PARTICIPATION
		v1.GET("/htmx/searchParticipants/:eventId", routesV1.SearchParticipantsHandler(queries))
		v1.POST("/htmx/selectParticipant/:eventId/:userId", routesV1.SelectParticipantHanlder(queries))
		v1.DELETE("/htmx/deselectParticipant/:eventId/:userId", routesV1.DeselectParticipantHanlder(queries))

		v1.POST("/htmx/addParticipant/:eventId/:userId/:status", routesV1.AddParticipantHandler(queries))
		v1.PUT("/htmx/inviteParticipants/:eventId", routesV1.InviteParticipantsHandler(queries))
		v1.PUT("/htmx/updateParticipant/:eventId/:userId/:status", routesV1.UpdateParticipantHandler(queries, routesV1.UpdateParticipantResponse(routesV1.Notification)))
		v1.PUT("/htmx/updateParticipant/:eventId/:userId/:status/newState", routesV1.UpdateParticipantHandler(queries, routesV1.UpdateParticipantResponse(routesV1.NewStatusAndButttons)))

		// v1.PUT("/htmx/updateParticipant/:eventId/:userId/:status", routes.UpdateParticipantHandler(queries))

		v1.DELETE("/htmx/removeParticipant/:eventId/:userId", routesV1.DeleteParticipantHandler(queries))

	}

	server.Run("localhost:19999")

}
