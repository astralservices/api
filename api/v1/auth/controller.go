package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/astralservices/api/api/v1/auth/providers/roblox"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/shareed2k/goth_fiber"
)

func AuthHandler(router fiber.Router) {
	router.Get("/callback/:provider", CallbackHandler)
	router.Post("/login/:provider", goth_fiber.BeginAuthHandler)
	router.Get("/login/:provider", func(c *fiber.Ctx) error {
		redirect := c.Query("redirect")
		provider := c.Params("provider")

		c.Cookie(&fiber.Cookie{
			Name:  "redirect",
			Value: redirect,
		})

		roblox := roblox.New(c, db.New(), redirect)

		if provider == "roblox" {
			return roblox.GenerateCodeForUser()
		}

		return goth_fiber.BeginAuthHandler(c)
	})
	router.Get("/logout/:provider", LogoutHandler)
	router.Get("/session", SessionHandler)

	authed := router.Use(utils.AuthMiddleware, utils.ProfileMiddleware)
	authed.Get("/providers", ProvidersHandler)
	authed.Get("/providers/:provider", ProviderHandler)
	authed.Post("/providers/:provider", UpdateProviderHandler)
	authed.Get("/status", StatusHandler)
	authed.Get("/gdpr", DataHandler)
	authed.Post("/delete", DeleteAccountHandler)
}

func InitGoth() {
	pgStore := postgres.New(postgres.Config{
		Host:       "db.aoeinliucinfkgwibqgr.supabase.co",
		Port:       5432,
		Database:   "postgres",
		Table:      "fiber_storage",
		Reset:      false,
		GCInterval: 10 * time.Second,
		SslMode:    "disable",
		Username:   "postgres",
		Password:   os.Getenv("POSTGRES_PASSWORD"),
	})

	sessions := session.New(session.Config{
		Storage:        pgStore,
		Expiration:     24 * time.Hour,
		KeyLookup:      fmt.Sprintf("cookie:%s", gothic.SessionName),
		CookieHTTPOnly: true,
	})

	goth_fiber.SessionStore = sessions

	goth.UseProviders(
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), utils.GetCallbackURL("discord"), discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeGuilds),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), utils.GetCallbackURL("lastfm")),
	)
}
