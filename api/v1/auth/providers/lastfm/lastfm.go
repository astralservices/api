package lastfm

import (
	"net/http"
	"time"

	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/nedpals/supabase-go"
)

type LastfmProvider struct {
	ctx      *fiber.Ctx
	database *supabase.Client
	user     goth.User
	redirect string
	domain   string
}

func New(c *fiber.Ctx, database *supabase.Client, user goth.User, redirect string, domain string) LastfmProvider {
	provider := LastfmProvider{
		ctx:      c,
		database: database,
		user:     user,
		redirect: redirect,
		domain:   domain,
	}
	return provider
}

func (p LastfmProvider) CreateUser() error {
	ctx, database, user, redirect := p.ctx, p.database, p.user, p.redirect

	discordUser := ctx.Locals("user").(utils.IProvider)

	var out []utils.IProvider

	insertErr := database.DB.From("providers").Insert(map[string]interface{}{
		"type":                   user.Provider,
		"user":                   discordUser.ID,
		"provider_id":            user.Name,
		"provider_access_token":  user.AccessToken,
		"provider_refresh_token": user.RefreshToken,
		"provider_expires_at":    user.ExpiresAt.UTC(),
		"provider_data":          user.RawData,
		"provider_avatar_url":    &user.AvatarURL,
		"provider_email":         &user.Email,
		"updated_at":             time.Now().UTC(),
	}).Execute(&out)

	if insertErr != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  insertErr.Error(),
		})
	}

	if redirect != "" {
		ctx.ClearCookie("redirect")
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: out,
		Code:   http.StatusOK,
	})
}

func (p LastfmProvider) UpdateUser() error {
	ctx, database, user, redirect := p.ctx, p.database, p.user, p.redirect
	var out []utils.IProvider

	discordUser := ctx.Locals("user").(utils.IProvider)

	insertErr := database.DB.From("providers").Update(map[string]interface{}{
		"type":                   user.Provider,
		"provider_id":            user.UserID,
		"provider_access_token":  user.AccessToken,
		"provider_refresh_token": user.RefreshToken,
		"provider_expires_at":    user.ExpiresAt.UTC(),
		"provider_data":          user.RawData,
		"provider_avatar_url":    &user.AvatarURL,
		"provider_email":         &user.Email,
		"updated_at":             time.Now().UTC(),
	}).Eq("provider_id", user.UserID).Eq("type", user.Provider).Eq("user", *discordUser.ID).Execute(&out)

	if insertErr != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  insertErr.Error(),
		})
	}

	if redirect != "" {
		ctx.ClearCookie("redirect")
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: out,
		Code:   http.StatusOK,
	})
}
