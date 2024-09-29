package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/a-h/templ"
	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"go-form/entities"
	templs_auth "go-form/templs/auth"
	templs_event "go-form/templs/event"
	templs "go-form/templs/generic"
	templs_user "go-form/templs/user"
)

// TODO: check what else makes sense as const
const Domain string = "localhost"
const CookieName string = "gohtmxplanner_cookie"
const CookieMaxAge = int(time.Minute * 10)
const CookieSecure = true
const CookieHTTPOnly = true

// TODO: timezones

func simpleRender(tc templ.Component) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", tc)
	}
}

func renderPage(tc templ.Component) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "Home", templs.Page(tc))
	}
}

func ErrorNotification(c *gin.Context, text string) {
	log.Println(text)
	c.HTML(200, "", templs.NotificationWithText(templs.BadReq, text))
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: reuse bad response
		// TODO: Nicer looking bad response
		// TODO: difft bad response, whether user is logged in/auth is expired or he lacks rights
		sessionId, err := c.Cookie(CookieName)
		if err != nil {
			fmt.Println(err)
			c.Redirect(301, "/loginOrRegister")
			c.Abort()
			return
		}

		sessionToFind := entities.Session{}

		res := db.Debug().Joins("User").First(&sessionToFind, entities.Session{ID: sessionId})

		if res.Error != nil {
			fmt.Println(res.Error.Error())
			c.Redirect(301, "/loginOrRegister")
			c.Abort()
			return
		}

		c.Set("auth-context", sessionToFind)

		c.Next()

	}
}

func GetAuthContext(c *gin.Context) *entities.Session {
	sesionPtr := c.MustGet("AuthContext").(*entities.Session)
	return sesionPtr
}

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
	// server.LoadHTMLFiles("./templs/event.html")

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	server.SetTrustedProxies(nil)

	server.GET("/", renderPage(templs.Home()))

	server.GET("/register", simpleRender(templs_auth.RegisterPage()))

	server.GET("/login", simpleRender(templs_auth.LoginPage()))

	server.GET("/loginOrRegister", simpleRender(templs_auth.LoginOrRegister()))

	type NewUserData struct {
		UserName string `form:"username" binding:"required"`
		Password string `form:"username" binding:"required"`
	}
	type LoginData = NewUserData

	server.POST("/htmx/register", func(c *gin.Context) {
		var login LoginData
		err := c.ShouldBind(&login)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}
		user := entities.NewUser(login.UserName, login.Password)

		// TODO: put db ops in transaction

		if err := db.Create(&user).Error; err != nil {
			fmt.Println(err.Error())
			ErrorNotification(c, err.Error())
			return
		}

		fmt.Println(login.UserName)
		sessionToken := xid.New().String()

		session := entities.NewSession(sessionToken, user.ID)

		if err := db.Create(&session).Error; err != nil {
			fmt.Println(err.Error())
			ErrorNotification(c, err.Error())
			return
		}

		c.SetCookie(CookieName, sessionToken, CookieMaxAge, "/", Domain, CookieSecure, CookieHTTPOnly)

		c.Header("HX-Redirect", "/")
		c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
	})

	server.POST("/htmx/login", func(c *gin.Context) {
		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}
		user := entities.NewUser(newUser.UserName, newUser.Password)

		// TODO: put db ops in transaction

		if err := db.Where(&user).Take(&user).Error; err != nil {
			fmt.Println(err.Error())
			ErrorNotification(c, "User does not exist, or wrong password")
			return
		}

		sessionToken := xid.New().String()
		session := entities.NewSession(sessionToken, user.ID)

		if err := db.Create(&session).Error; err != nil {
			fmt.Println(err.Error())
			ErrorNotification(c, err.Error())
			return
		}

		c.SetCookie(CookieName, sessionToken, CookieMaxAge, "/", Domain, CookieSecure, CookieHTTPOnly)

		c.Header("HX-Redirect", "/")
		// TODO: figure the nicest way to combine redirect & notifications (special header?)
		// c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
	})

	server.GET("/createEvent", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create Event", templs.Page(templs_event.CreateEvent()))
	})

	server.GET("/events", func(c *gin.Context) {
		var events []entities.CalendarEvent
		db.Find(&events)
		renderPage(templs_event.EventList(&events))(c)
	})

	type NewEventData struct {
		Title    string `form:"title" binding:"required"`
		DateTime string `form:"date-time" binding:"required"`
	}

	server.POST("/htmx/createEvent", func(c *gin.Context) {
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

	})

	server.GET("/users", func(c *gin.Context) {
		var users []entities.User
		db.Find(&users)
		renderPage(templs_user.UserList(&users))(c)
	})

	server.POST("/htmx/createUser", func(c *gin.Context) {
		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		res := db.Create(entities.NewUser(newUser.UserName, newUser.Password))
		if res.Error != nil {
			println("bad db")
			fmt.Println(res.Error)
			c.HTML(200, "", templs.Notification(templs.BadReq))
		} else {
			println("succ")
			fmt.Println(newUser.UserName)
			c.Header("HX-Redirect", "/users")
			c.HTML(http.StatusCreated, "", templs.Notification(templs.Success))
		}

	})

	server.GET("/event/:id", func(c *gin.Context) {
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
	})

	server.DELETE("htmx/deleteEvent/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Delete(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			simpleRender(templs.Notification(templs.BadReq))(c)
		} else {
			simpleRender(templs.NotificationOob(templs.Success))(c)
		}
	})

	server.DELETE("htmx/deleteUser/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var user entities.User
		res := db.Delete(&user, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			simpleRender(templs.Notification(templs.BadReq))(c)
		} else {
			simpleRender(templs.NotificationOob(templs.Success))(c)
		}
	})

	server.GET("/updateEvent/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ev entities.CalendarEvent
		res := db.Take(&ev, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		renderPage(templs_event.UpdateEvent(&ev))(c)
	})

	server.GET("/updateUser/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var user entities.User
		res := db.Take(&user, id)
		if res.Error != nil {
			log.Println("Nooooooooo")
			log.Println(res.Error)
			c.HTML(http.StatusNotFound, "Not Found", templs.FoOhFo())
			return
		}
		renderPage(templs_user.UpdateUser(&user))(c)
	})

	server.PUT("/htmx/updateEvent/:id", func(c *gin.Context) {
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
	})

	server.PUT("/htmx/updateUser/:id", func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		updatedUser := entities.NewUser(newUser.UserName, newUser.Password)
		updatedUser.ID = uint(id)

		res := db.Save(&updatedUser)
		if res.Error != nil {
			println(res.Error.Error())
			fmt.Println(res.Error)
			c.HTML(200, "", templs.NotificationOobWithText(templs.BadReq, res.Error.Error()))
		} else {
			println("succ")
			c.Header("HX-Redirect", "/users")
			c.HTML(http.StatusCreated, "", templs.NotificationOob(templs.Success))
		}
	})

	server.Use(AuthMiddleware(db))
	{
		server.POST("/htmx/logout", func(c *gin.Context) {

			session := c.MustGet("auth-context").(entities.Session)

			if err := db.Delete(&session).Error; err != nil {
				fmt.Println(err.Error())
				ErrorNotification(c, err.Error())
				return
			}

			c.SetCookie(CookieName, "", -1, "/", Domain, CookieSecure, CookieHTTPOnly)

			c.Header("HX-Redirect", "/login")
			// TODO: figure the nicest way to combine redirect & notifications (special header?)
			c.String(200, "")
		})

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
