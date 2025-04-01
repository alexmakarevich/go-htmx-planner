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

// TODO: embed FE files in prod, but serve from filesystem in dev (for easier hot reloads)
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

		println("ROUTE", ctx.Request.URL.Path)

		isAuthNeeded := IsAuthNeeded(ctx)
		if isAuthNeeded {
			err := routesV2.CheckAuth(queries, ctx)
			// TODO: query param to forward to given path after login
			if err != nil {
				println("AUTH CHECK FAILED")
				fmt.Println("session cookie not found")
				fmt.Println(err)
				ctx.Redirect(307, "/v2/login")
				ctx.Abort()
				return
			}
			println("AUTH CHECK SUCCEEDED")
		} else {
			println("AUTH CHECK SKIPPED")
		}

		RedirectOrPassthrough(ctx)
	})

	v2 := server.Group("/v2")
	routesV2.MakeV2Routes(v2, queries)

	v1 := server.Group("/v1")
	routesV1.MakeV1Routes(v1, queries)

	server.Run("localhost:19999")

}

func IsAuthNeeded(ctx *gin.Context) bool {
	// v1 auth done separately
	if len(ctx.Request.URL.Path) >= 4 && ctx.Request.URL.Path[0:4] == "/v1/" {
		return false
	}

	if ctx.Request.URL.Path == "/v2/login" || ctx.Request.URL.Path == "/v2/api/login" {
		// TODO: FE-redirect after api login succeeds
		return false
	}

	if ctx.Request.URL.Path == "/v2" || ctx.Request.URL.Path == "/v2/" {
		if ctx.Request.URL.Query().Get("fe-route") == "/v2/login" {
			return false
		}
	}

	// always allow to download public site files
	if len(ctx.Request.URL.Path) >= 11 && ctx.Request.URL.Path[0:11] == "/v2/assets/" {
		println("NO redirect:", ctx.Request.URL.Path)
		ctx.Next()
		return false
	}
	// otherwise auth is necessary
	return true
}

func RedirectOrPassthrough(ctx *gin.Context) {

	if len(ctx.Request.URL.Path) > 11 && ctx.Request.URL.Path[0:11] == "/v2/assets/" {
		ctx.Next()
		return
	}

	if ctx.Request.URL.Path == "/" {
		ctx.Redirect(301, "/v2/")
		return
	}

	if len(ctx.Request.URL.Path) >= 4 && ctx.Request.URL.Path[0:4] == "/v1/" {
		ctx.Next()
		return
	}

	// we are somewhere in v2
	if len(ctx.Request.URL.Path) >= 4 && ctx.Request.URL.Path[0:4] == "/v2/" {
		// we are aat the root of v2 - should get index.html via static handler later
		if len(ctx.Request.URL.Path) == 4 {
			ctx.Next()
			return
		}
		// all API reqs are checked
		if len(ctx.Request.URL.Path) >= 8 && ctx.Request.URL.Path[4:8] == "api/" {
			ctx.Next()
			return
		}
		// // TODO: check query/body
		ctx.Redirect(307, "/v2/?fe-route="+ctx.Request.URL.Path)
		return
	}

	println("ERROR: unexpected route - not covered by initial global middleware")
	ctx.Status(404)
	return
}
