package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	mw "github.com/mattwhip/icenine-services/middleware"
	"github.com/mattwhip/icenine-service-daily_bonus/models"
	"github.com/gobuffalo/x/sessions"
)

// ENV is the go environment
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

var jwtSigningSecret = os.Getenv("JWT_SIGNING_SECRET")

// App is the buffalo app that drives routing, handling, middleware, etc
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_daily_bonus_session",
		})
		// Automatically redirect to SSL
		app.Use(forcessl.Middleware(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		app.Use(popmw.Transaction(models.DB))

		app.GET("/", HomeHandler)

		api := app.Group("/api")
		api.Use(mw.JWTVerification(jwtSigningSecret))
		api.POST("/status", StatusHandler)
		api.POST("/play", PlayHandler)
	}

	return app
}
