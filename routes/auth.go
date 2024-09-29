package routes

import (
	"fmt"
	"go-form/entities"
	templs "go-form/templs/generic"
	"net/http"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// TODO: check what else makes sense as const
const Domain string = "localhost"
const CookieName string = "gohtmxplanner_cookie"
const CookieMaxAge = int(time.Minute * 10)
const CookieSecure = true
const CookieHTTPOnly = true

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: reuse bad response
		// TODO: Nicer looking bad response
		// TODO: difft bad response, whether user is logged in/auth is expired or he lacks rights
		sessionId, err := c.Cookie(CookieName)
		if err != nil {
			fmt.Println(err)
			// c.Header("HX-Redirect", "/loginOrRegister")
			c.Redirect(307, "/loginOrRegister")
			c.Abort()
			return
		}

		fmt.Println("sessionId")
		fmt.Println(sessionId)

		sessionToFind := entities.Session{}

		res := db.Debug().Joins("User").First(&sessionToFind, entities.Session{ID: sessionId})

		if res.Error != nil {
			fmt.Println("ERROR:")
			fmt.Println(res.Error.Error())
			// c.Header("HX-Redirect", "/loginOrRegister")
			c.Redirect(307, "/loginOrRegister")
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

func LogoutHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

type NewUserData struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"username" binding:"required"`
}
type LoginData = NewUserData

func LoginHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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

		fmt.Println("New Session:")
		fmt.Println(sessionToken)

		c.SetCookie(CookieName, sessionToken, CookieMaxAge, "/", Domain, CookieSecure, CookieHTTPOnly)
		c.Header("HX-Redirect", "/")
		// c.String(200, "")
		// c.Redirect(300, "/")
		// TODO: figure the nicest way to combine redirect & notifications (special header?)
		// TODO: i have to have HTMX in the browser for this to work
		c.HTML(http.StatusOK, "", templs.Notification(templs.Success))
	}
}

func RegisterHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}
