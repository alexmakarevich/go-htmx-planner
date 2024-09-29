package main

import (
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"

	"go-form/entities"
	"go-form/routes"
	templs_auth "go-form/templs/auth"
	templs "go-form/templs/generic"
)

// TODO: timezones

func main() {
	// TODO: FUTURE: prod.db, migrations, backups, etc.
	db, err := gorm.Open(sqlite.Open("db/dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// TODO: hanlde more errors
	db.AutoMigrate(&entities.CalendarEvent{}, &entities.User{}, &entities.Session{})

	server := gin.Default()
	server.Static("/public", "./public")

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	server.SetTrustedProxies(nil)

	server.POST("/htmx/register", routes.RegisterHandler(db))
	server.POST("/htmx/login", routes.LoginHandler(db))

	server.GET("/register", routes.SimpleRender(templs_auth.RegisterPage()))
	server.GET("/login", routes.SimpleRender(templs_auth.LoginPage()))
	server.GET("/loginOrRegister", routes.SimpleRender(templs_auth.LoginOrRegister()))

	server.Use(routes.AuthMiddleware(db))
	{

		server.GET("/", routes.RenderPage(templs.Home()))

		server.POST("/htmx/logout", routes.LogoutHandler(db))

		// USER
		server.POST("/htmx/createUser", routes.CreateUserHandler(db))

		server.GET("/users", routes.ListUsersHandler(db))

		server.GET("/updateUser/:id", routes.UpdateUserPageHandler(db))
		server.PUT("/htmx/updateUser/:id", routes.UpdateUserHandler(db))

		server.DELETE("htmx/deleteUser/:id", routes.DeleteUserHandler(db))

		// EVENT
		server.GET("/createEvent", routes.CreateEventPageHandler(db))
		server.POST("/htmx/createEvent", routes.CreateEventHandler(db))

		server.GET("/events", routes.ListEventsPageHandler(db))
		server.GET("/event/:id", routes.GetEventPageHandler(db))

		server.GET("/updateEvent/:id", routes.UpdateEventPageHandler(db))
		server.PUT("/htmx/updateEvent/:id", routes.UpdateEventHandler(db))

		server.DELETE("htmx/deleteEvent/:id", routes.DeleteEventHandler(db))

		// TEMP TEST ROUTE
		server.GET("/protected", func(c *gin.Context) {
			session := c.MustGet("auth-context").(entities.Session)

			// it would print: "12345"
			log.Println("USER FROM CONTEXT")
			log.Println(session.User)
			c.HTML(http.StatusCreated, "", templs.NotificationWithText(templs.Success, "protected content"))

		})
	}

	server.Run("localhost:19999")

}
