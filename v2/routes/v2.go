package routes

import (
	"go-form/sqlc/db_entities"

	"github.com/gin-gonic/gin"
)

func MakeV2Routes(router gin.IRouter, q *db_entities.Queries) {
	router.StaticFile("/", "./v2/svelte/dist/index.html")
	router.Static("/assets", "./v2/svelte/dist/assets")
	router.StaticFile("/vite.svg", "./v2/svelte/dist/vite.svg")

	MakeUserRoutes(router, q)
}
