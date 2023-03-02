package main

import (
	"context"
	"errors"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/arikkfir/greenstar/backend/util/bootutil"
	"github.com/arikkfir/greenstar/backend/util/ginutil"
	"github.com/gin-contrib/cors"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

func main() {
	config := Config{}
	bootutil.Boot(&config)

	// Service context
	ctx := context.Background()

	// Setup Google OAuth
	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.Auth.Google.CallbackURL,
		ClientID:     config.Auth.Google.ClientID,
		ClientSecret: config.Auth.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Setup routes
	router := ginutil.NewGin(config.DevMode)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:        false,
		AllowOrigins:           config.HTTP.CORS.AllowedOrigins,
		AllowMethods:           config.HTTP.CORS.AllowMethods,
		AllowHeaders:           config.HTTP.CORS.AllowHeaders,
		AllowCredentials:       !config.HTTP.CORS.DisableCredentials,
		ExposeHeaders:          config.HTTP.CORS.ExposeHeaders,
		MaxAge:                 config.HTTP.CORS.MaxAge,
		AllowBrowserExtensions: false,
		AllowFiles:             false,
		AllowWebSockets:        true,
		AllowWildcard:          true,
	}))
	router.GET("/google/login", CreateAuthGoogleLoginHandler(googleOauthConfig, config.HTTP.StateCookieName, config.HTTP.AppURL))
	router.GET("/google/callback", CreateAuthGoogleCallbackHandler(googleOauthConfig, config.HTTP.StateCookieName, config.HTTP.ClaimsCookieName, config.HTTP.ShouldUseSecureCookies(), "greenstar.auth"))
	router.GET("/google/logout", CreateAuthGoogleLogoutHandler(
		config.HTTP.StateCookieName,
		config.HTTP.ClaimsCookieName,
		config.HTTP.ShouldUseSecureCookies(),
		config.HTTP.AppURL),
	)
	router.GET(
		"/user",
		auth.CreateVerifyTokenMiddleware(
			config.HTTP.ClaimsCookieName,
			googleOauthConfig.ClientSecret,
			"greenstar.auth",
			"greenstar.auth"),
		CreateAuthGoogleUserInfoHandler(),
	)

	// Setup HTTP server
	httpServer := &http.Server{Addr: ":" + strconv.Itoa(config.HTTP.Port), Handler: router}

	// Start server
	log.Ctx(ctx).Info().Str("addr", httpServer.Addr).Msg("Starting HTTP server")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("HTTP server failed")
	}
}
