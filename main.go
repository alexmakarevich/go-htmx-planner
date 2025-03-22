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

	"go-form/routes"
	templs_auth "go-form/templs/auth"
	templs "go-form/templs/generic"

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
	v1.Static("/public", "./public")
	v1.POST("/htmx/register", routes.RegisterHandler(queries))
	v1.POST("/htmx/login", routes.LoginHandler(queries))

	v1.GET("/register", routes.SimpleRender(templs_auth.RegisterPage()))
	v1.GET("/login", routes.SimpleRender(templs_auth.LoginPage()))
	v1.GET("/loginOrRegister", routes.SimpleRender(templs_auth.LoginOrRegister()))

	v1.Use(routes.AuthMiddleware(queries))
	{
		v1.GET("/", routes.RenderPage(templs.Home()))

		v1.POST("/htmx/logout", routes.LogoutHandler(queries))
		v1.GET("/settings", routes.SettingsPageHandler(queries))
		v1.DELETE("/htmx/delete-self", routes.DeleteSelfHandler(queries))

		// USER
		v1.GET("/createUser", routes.CreateUserPageHandler(queries))
		v1.POST("/htmx/createUser", routes.CreateUserHandler(queries))

		v1.GET("/users", routes.ListUsersPageHandler(queries))

		v1.GET("/updateUser/:id", routes.UpdateUserPageHandler(queries))
		v1.PUT("/htmx/updateUser/:id", routes.UpdateUserHandler(queries))

		v1.DELETE("/htmx/deleteUser/:id", routes.DeleteUserHandler(queries))

		// EVENT
		v1.GET("/createEvent", routes.CreateEventPageHandler(queries))
		v1.POST("/htmx/createEvent", routes.CreateEventHandler(queries))

		v1.GET("/events", routes.ListEventsPageHandler(queries))
		v1.GET("/event/:id", routes.ViewOrUpdateEventPageHandler(queries, false))
		v1.GET("/event/:id/invite", routes.ViewOrUpdateEventPageHandler(queries, false))

		v1.GET("/myInvites", routes.ListInvitesPagehandler(queries))

		v1.GET("/updateEvent/:id", routes.ViewOrUpdateEventPageHandler(queries, true))

		v1.PUT("/htmx/updateEvent/:id", routes.UpdateEventHandler(queries))

		v1.DELETE("/htmx/deleteEvent/:id", routes.DeleteEventHandler(queries))

		// PARTICIPATION
		v1.GET("/htmx/searchParticipants/:eventId", routes.SearchParticipantsHandler(queries))
		v1.POST("/htmx/selectParticipant/:eventId/:userId", routes.SelectParticipantHanlder(queries))
		v1.DELETE("/htmx/deselectParticipant/:eventId/:userId", routes.DeselectParticipantHanlder(queries))

		v1.POST("/htmx/addParticipant/:eventId/:userId/:status", routes.AddParticipantHandler(queries))
		v1.PUT("/htmx/inviteParticipants/:eventId", routes.InviteParticipantsHandler(queries))
		v1.PUT("/htmx/updateParticipant/:eventId/:userId/:status", routes.UpdateParticipantHandler(queries, routes.UpdateParticipantResponse(routes.Notification)))
		v1.PUT("/htmx/updateParticipant/:eventId/:userId/:status/newState", routes.UpdateParticipantHandler(queries, routes.UpdateParticipantResponse(routes.NewStatusAndButttons)))

		// v1.PUT("/htmx/updateParticipant/:eventId/:userId/:status", routes.UpdateParticipantHandler(queries))

		v1.DELETE("/htmx/removeParticipant/:eventId/:userId", routes.DeleteParticipantHandler(queries))

	}

	server.Run("localhost:19999")

}
