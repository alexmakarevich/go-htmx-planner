package routes

import (
	"fmt"
	"go-form/sqlc/db_entities"
	templs "go-form/templs/generic"
	"net/http"
	"time"

	"github.com/rs/xid"

	"github.com/gin-gonic/gin"
)

// TODO: check what else makes sense as const
const Domain string = "localhost"
const CookieName string = "gohtmxplanner_cookie"
const CookieMaxAge = int(time.Minute * 10)
const CookieSecure = true
const CookieHTTPOnly = true

func AuthMiddleware(q *db_entities.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: reuse bad response
		// TODO: Nicer looking bad response
		// TODO: difft bad response, whether user is logged in/auth is expired or he lacks rights
		sessionId, err := c.Cookie(CookieName)
		if err != nil {
			fmt.Println("session cookie not found")
			fmt.Println(err)
			// c.Header("HX-Redirect", "/loginOrRegister")
			// c.String(200, "")
			c.Redirect(307, "/loginOrRegister")
			c.Abort()
			return
		}

		fmt.Println("sessionId")
		fmt.Println(sessionId)

		// TODO: join w/ user

		sessionWithUser, err := q.GetSessionWithUser(c, sessionId)

		if err != nil {
			fmt.Println("ERROR:")
			fmt.Println(err.Error())
			// c.Header("HX-Redirect", "/loginOrRegister")
			c.Redirect(307, "/loginOrRegister")
			c.Abort()
			return
		}

		c.Set("auth-context", sessionWithUser)

		c.Next()

	}
}

func GetAuthContext(c *gin.Context) db_entities.GetSessionWithUserRow {
	sesionPtr := c.MustGet("auth-context").(db_entities.GetSessionWithUserRow)
	return sesionPtr
}

func LogoutHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := c.MustGet("auth-context").(db_entities.GetSessionWithUserRow)

		err := q.DeleteSession(c, session.SessionID)

		if err != nil {
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
	Password string `form:"password" binding:"required"`
}
type LoginData = NewUserData

func LoginHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		var newUser NewUserData
		err := c.ShouldBind(&newUser)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}
		// TODO: put db ops in transaction

		user, err := q.FindUser(c, db_entities.FindUserParams{Password: newUser.Password, UserName: newUser.UserName})

		if err != nil {
			fmt.Println(err.Error())
			ErrorNotification(c, "User does not exist, or wrong password")
			return
		}

		sessionToken := xid.New().String()

		if _, err := q.CreateSession(c, db_entities.CreateSessionParams{ID: sessionToken, UserID: user.ID}); err != nil {
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

func RegisterHandler(q *db_entities.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		var login LoginData
		err := c.ShouldBind(&login)

		if err != nil {
			ErrorNotification(c, err.Error())
			return
		}

		// TODO: put db ops in transaction

		user, err := q.CreateUser(c, db_entities.CreateUserParams{UserName: login.UserName, Password: login.Password})

		if err != nil {
			println("bad db")
			fmt.Println(err.Error())
			c.HTML(200, "", templs.Notification(templs.BadReq))
			return
		}
		fmt.Println(login.UserName)
		sessionToken := xid.New().String()

		if _, err := q.CreateSession(c, db_entities.CreateSessionParams{ID: sessionToken, UserID: user.ID}); err != nil {
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
