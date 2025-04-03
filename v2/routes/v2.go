package routes

import (
	"context"
	"database/sql"
	"go-form/sqlc/db_entities"

	"github.com/gin-gonic/gin"
)

func MakeV2Routes(router gin.IRouter, q *db_entities.Queries, db *sql.DB, globalCtx context.Context) {

	router.StaticFile("/", "./v2/svelte/dist/index.html")
	router.Static("/assets", "./v2/svelte/dist/assets")
	// router.StaticFile("/vite.svg", "./v2/svelte/dist/vite.svg")

	api := router.Group("/api")

	MakeAuthRoutes(api, q)

	MakeUserRoutes(api, q)
	MakeEventRoutes(api, q, db, globalCtx)

}
