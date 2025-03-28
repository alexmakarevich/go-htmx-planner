package routes

import (
	"fmt"
	"go-form/sqlc/db_entities"
	protos "go-form/v2/proto-go"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
)

func MakeUserRoutes(router gin.IRouter, q *db_entities.Queries) {
	router.GET("/list-users", func(c *gin.Context) {
		users, err := q.ListUsers(c)
		fmt.Println(users)
		fmt.Println(len(users))

		if err != nil {
			fmt.Println(err.Error())
			c.Status(500)
			return
		}
		userList := &protos.UserList{Users: []*protos.ShareableUserData{}}

		for _, u := range users {
			userList.Users = append(userList.Users, &protos.ShareableUserData{
				Name: u.UserName,
				Id:   int32(u.ID),
			})
		}

		data, err := proto.Marshal(userList)
		if err != nil {
			fmt.Println("Error marshaling:", err)
			c.AbortWithError(500, err)
			return
		}

		c.Data(200, "application/x-protobuf", data)
	})

	router.POST("/create-user", func(c *gin.Context) {
		newUser := protos.CreateUserParams{}
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(500, err)
		}
		if err := proto.Unmarshal(data, &newUser); err != nil {
			c.AbortWithError(500, err)
		}
		_, err = q.CreateUser(c, db_entities.CreateUserParams{UserName: newUser.Name, Password: newUser.Password})
		// TODO: better handling of specific DB errors (unique constraint etc.)
		if err != nil {
			fmt.Println(err.Error())
			c.String(500, "db op failed")
		} else {
			c.String(http.StatusCreated, "user created")
		}
	})
}
