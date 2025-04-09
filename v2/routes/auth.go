package routes

import (
	"errors"
	"fmt"
	"go-form/sqlc/db_entities"
	"io"
	"net/http"
	"time"

	"github.com/rs/xid"
	"google.golang.org/protobuf/proto"

	protos "go-form/v2/proto-go"
	"go-form/v2/utils"

	"github.com/gin-gonic/gin"
)

// TODO: check what else makes sense as const
// FYI: some of them are "ocnst" just for prettiness
const Domain string = "localhost"
const CookieName string = "gohtmxplanner_cookie"
const CookieMaxAge = int(time.Minute * 10)
const CookieSecure = true
const CookieHTTPOnly = true

func MakeAuthRoutes(router gin.IRouter, q *db_entities.Queries) {
	router.POST("/login", func(c *gin.Context) {
		// loginParams, err := utils.ParseProto[*protos.LoginParams](
		// 	c,
		// )

		// loginParams := protos.LoginParams{}

		// data, err := io.ReadAll(c.Request.Body)
		// if err != nil {
		// 	println("i can't reeeeeeaaad")
		// 	c.Status(401)
		// 	return
		// }
		// if err := proto.Unmarshal(data, &loginParams); err != nil {
		// 	println("couldn't unmarhal")
		// 	c.Status(401)
		// 	return
		// }

		loginParams := protos.LoginParams{}
		err := utils.ParseProto(
			c, &loginParams,
		)

		if err != nil {
			msg := "Couldn't parse params"
			utils.SendProto(c, 500, &protos.ErrorResponse{Message: &msg})
			return
		}

		user, err := q.FindUser(c, db_entities.FindUserParams{Password: *loginParams.Password, UserName: *loginParams.Name})

		if err != nil {
			msg := "User doesn't exist"
			utils.SendProto(c, 401, &protos.ErrorResponse{Message: &msg})
			return
		}

		sessionToken := xid.New().String()

		if _, err := q.CreateSession(c, db_entities.CreateSessionParams{ID: sessionToken, UserID: user.ID}); err != nil {
			fmt.Println(err.Error())
			msg := "Couldn't create session"
			utils.SendProto(c, 500, &protos.ErrorResponse{Message: &msg})
			return
		}

		fmt.Println("New Session:")
		fmt.Println(sessionToken)

		c.SetCookie(CookieName, sessionToken, CookieMaxAge, "/", Domain, CookieSecure, CookieHTTPOnly)
		c.Status(http.StatusOK)

	})

	router.POST("/logout", func(c *gin.Context) {
		session := c.MustGet("auth-context").(db_entities.GetSessionWithUserRow)

		err := q.DeleteSession(c, session.SessionID)

		c.SetCookie(CookieName, "", -1, "/", Domain, CookieSecure, CookieHTTPOnly)
		if err != nil {
			fmt.Println(err.Error())
			c.String(500, "Could not delete session")
			return
		}

		c.String(200, "")
	})
}

func HandleLogout(q *db_entities.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		newUser := protos.CreateUserParams{}
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(500, err)
		}
		if err := proto.Unmarshal(data, &newUser); err != nil {
			c.AbortWithError(500, err)
		}
		user, err := q.CreateUser(c, db_entities.CreateUserParams{UserName: newUser.Name, Password: newUser.Password})
		// TODO: better handling of specific DB errors (unique constraint etc.)
		if err != nil {
			fmt.Println(err.Error())
			c.String(500, "db op failed")
		} else {
			c.String(http.StatusCreated, "user created")
		}

		sessionToken := xid.New().String()

		if _, err := q.CreateSession(c, db_entities.CreateSessionParams{ID: sessionToken, UserID: user.ID}); err != nil {
			fmt.Println(err.Error())
			c.String(500, "Could not create session")
			return
		}

		fmt.Println("New Session:")
		fmt.Println(sessionToken)

		c.SetCookie(CookieName, sessionToken, CookieMaxAge, "/", Domain, CookieSecure, CookieHTTPOnly)
		c.Status(http.StatusOK)
	}
}

func AuthMiddleware(q *db_entities.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: reuse bad response
		// TODO: Nicer looking bad response
		// TODO: difft bad response, whether user is logged in/auth is expired or he lacks rights
		sessionId, err := c.Cookie(CookieName)
		if err != nil {
			fmt.Println("session cookie not found")
			fmt.Println(err)
			c.Redirect(307, "/login")
			c.Abort()
			return
		}

		fmt.Println("sessionId")
		fmt.Println(sessionId)

		sessionWithUser, err := q.GetSessionWithUser(c, sessionId)

		if err != nil {
			fmt.Println("ERROR:")
			fmt.Println(err.Error())
			c.Redirect(307, "login")
			c.Abort()
			return
		}

		c.Set("auth-context", sessionWithUser)
		c.Next()
	}
}

func CheckAuth(q *db_entities.Queries, c *gin.Context) error {
	sessionId, err := c.Cookie(CookieName)
	if err != nil {
		return errors.New("cookie not found in context")
	}

	fmt.Println("sessionId")
	fmt.Println(sessionId)

	sessionWithUser, err := q.GetSessionWithUser(c, sessionId)

	if err != nil {
		return errors.New("cookie/user not found in DB")
	}

	c.Set("auth-context", sessionWithUser)

	return nil
}

func GetAuthContext(c *gin.Context) db_entities.GetSessionWithUserRow {
	sesionPtr := c.MustGet("auth-context").(db_entities.GetSessionWithUserRow)
	return sesionPtr
}
