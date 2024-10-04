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
	server.Static("/public", "./public")

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	server.SetTrustedProxies(nil)

	server.POST("/htmx/register", routes.RegisterHandler(queries))
	server.POST("/htmx/login", routes.LoginHandler(queries))

	server.GET("/register", routes.SimpleRender(templs_auth.RegisterPage()))
	server.GET("/login", routes.SimpleRender(templs_auth.LoginPage()))
	server.GET("/loginOrRegister", routes.SimpleRender(templs_auth.LoginOrRegister()))

	server.Use(routes.AuthMiddleware(queries))
	{
		server.GET("/", routes.RenderPage(templs.Home()))

		server.POST("/htmx/logout", routes.LogoutHandler(queries))
		server.GET("/settings", routes.SettingsPageHandler(queries))
		server.DELETE("/htmx/delete-self", routes.DeleteSelfHandler(queries))

		// USER
		server.GET("/createUser", routes.CreateUserPageHandler(queries))
		server.POST("/htmx/createUser", routes.CreateUserHandler(queries))

		server.GET("/users", routes.ListUsersPageHandler(queries))

		server.GET("/updateUser/:id", routes.UpdateUserPageHandler(queries))
		server.PUT("/htmx/updateUser/:id", routes.UpdateUserHandler(queries))

		server.DELETE("htmx/deleteUser/:id", routes.DeleteUserHandler(queries))

		// EVENT
		server.GET("/createEvent", routes.CreateEventPageHandler(queries))
		server.POST("/htmx/createEvent", routes.CreateEventHandler(queries))

		server.GET("/events", routes.ListEventsPageHandler(queries))
		server.GET("/event/:id", routes.GetEventPageHandler(queries))

		server.GET("/updateEvent/:id", routes.UpdateEventPageHandler(queries))
		server.PUT("/htmx/updateEvent/:id", routes.UpdateEventHandler(queries))

		server.DELETE("htmx/deleteEvent/:id", routes.DeleteEventHandler(queries))
	}

	server.Run("localhost:19999")

}
