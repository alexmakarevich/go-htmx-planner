package routes

import (
	"go-form/sqlc/db_entities"

	templs_auth "go-form/v1/templs/auth"
	templs "go-form/v1/templs/generic"

	"github.com/gin-gonic/gin"
)

func MakeV1Routes(router gin.IRouter, queries *db_entities.Queries) {
	router.Static("/public", "./v1/public")
	router.POST("/htmx/register", RegisterHandler(queries))
	router.POST("/htmx/login", LoginHandler(queries))

	router.GET("/register", SimpleRender(templs_auth.RegisterPage()))
	router.GET("/login", SimpleRender(templs_auth.LoginPage()))
	router.GET("/loginOrRegister", SimpleRender(templs_auth.LoginOrRegister()))

	router.Use(AuthMiddleware(queries))
	{
		router.GET("/", RenderPage(templs.Home()))

		router.POST("/htmx/logout", LogoutHandler(queries))
		router.GET("/settings", SettingsPageHandler(queries))
		router.DELETE("/htmx/delete-self", DeleteSelfHandler(queries))

		// USER
		router.GET("/createUser", CreateUserPageHandler(queries))
		router.POST("/htmx/createUser", CreateUserHandler(queries))

		router.GET("/users", ListUsersPageHandler(queries))

		router.GET("/updateUser/:id", UpdateUserPageHandler(queries))
		router.PUT("/htmx/updateUser/:id", UpdateUserHandler(queries))

		router.DELETE("/htmx/deleteUser/:id", DeleteUserHandler(queries))

		// EVENT
		router.GET("/createEvent", CreateEventPageHandler(queries))
		router.POST("/htmx/createEvent", CreateEventHandler(queries))

		router.GET("/events", ListEventsPageHandler(queries))
		router.GET("/event/:id", ViewOrUpdateEventPageHandler(queries, false))
		router.GET("/event/:id/invite", ViewOrUpdateEventPageHandler(queries, false))

		router.GET("/myInvites", ListInvitesPagehandler(queries))

		router.GET("/updateEvent/:id", ViewOrUpdateEventPageHandler(queries, true))

		router.PUT("/htmx/updateEvent/:id", UpdateEventHandler(queries))

		router.DELETE("/htmx/deleteEvent/:id", DeleteEventHandler(queries))

		// PARTICIPATION
		router.GET("/htmx/searchParticipants/:eventId", SearchParticipantsHandler(queries))
		router.POST("/htmx/selectParticipant/:eventId/:userId", SelectParticipantHanlder(queries))
		router.DELETE("/htmx/deselectParticipant/:eventId/:userId", DeselectParticipantHanlder(queries))

		router.POST("/htmx/addParticipant/:eventId/:userId/:status", AddParticipantHandler(queries))
		router.PUT("/htmx/inviteParticipants/:eventId", InviteParticipantsHandler(queries))
		router.PUT("/htmx/updateParticipant/:eventId/:userId/:status", UpdateParticipantHandler(queries, UpdateParticipantResponse(Notification)))
		router.PUT("/htmx/updateParticipant/:eventId/:userId/:status/newState", UpdateParticipantHandler(queries, UpdateParticipantResponse(NewStatusAndButttons)))

		// router.PUT("/htmx/updateParticipant/:eventId/:userId/:status", UpdateParticipantHandler(queries))

		router.DELETE("/htmx/removeParticipant/:eventId/:userId", DeleteParticipantHandler(queries))

	}

}
