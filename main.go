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

	routesV1 "go-form/v1/routes"
	routesV2 "go-form/v2/routes"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlc/schema.sql
var ddl string

// TODO: embed FE files
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

	ginHtmlRenderer := server.HTMLRender
	server.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	server.SetTrustedProxies(nil)

	// forwarding to correct FE-route of SPA
	server.Use(func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}

		fmt.Printf("%+v\n", ctx.Request.URL)

		if ctx.Request.URL.Path == "/v2/" ||
			(len(ctx.Request.URL.Path) > 6 && ctx.Request.URL.Path[0:7] == "/v2/api") ||
			(len(ctx.Request.URL.Path) > 9 && ctx.Request.URL.Path[0:10] == "/v2/assets") {
			ctx.Next()
			println("NO redirect:", ctx.Request.URL.Path)
			return
		}

		if len(ctx.Request.URL.Path) > 3 && ctx.Request.URL.Path[0:4] == "/v2/" {
			println("REDIRECT:", ctx.Request.URL.Path)
			// // TODO: check query/body
			ctx.Redirect(307, "/v2/?fe-route="+ctx.Request.URL.Path)
			return
		}

		ctx.Next()
		return
	})

	v2 := server.Group("/v2")
	routesV2.MakeV2Routes(v2, queries)

	v1 := server.Group("/v1")
	routesV1.MakeV1Routes(v1, queries)

	server.Run("localhost:19999")

}
