package auth

import (
	"os"

	"github.com/astralservices/api/api/v1/auth/providers/roblox"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
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
	authed.Get("/status", StatusHandler)
}

func InitGoth() {
	goth.UseProviders(
		discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), utils.GetCallbackURL("discord"), discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeGuilds),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), utils.GetCallbackURL("lastfm")),
	)
}