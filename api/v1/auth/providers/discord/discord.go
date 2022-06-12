package discord

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/nedpals/supabase-go"
)

type DiscordProvider struct {
	ctx      *fiber.Ctx
	database *supabase.Client
	user     goth.User
	redirect string
	domain   string
}

func New(c *fiber.Ctx, database *supabase.Client, user goth.User, redirect string, domain string) DiscordProvider {
	provider := DiscordProvider{
		ctx:      c,
		database: database,
		user:     user,
		redirect: redirect,
		domain:   domain,
	}
	return provider
}

func (p DiscordProvider) CreateUser() error {
	ctx, database, user, redirect, domain := p.ctx, p.database, p.user, p.redirect, p.domain

	var out []utils.IProvider

	insertErr := database.DB.From("providers").Insert(map[string]interface{}{
		"type":                   user.Provider,
		"provider_id":            user.UserID,
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

	log.Println("making profile", out)

	var profile []utils.IProfile

	profileErr := database.DB.From("profiles").Insert(map[string]interface{}{
		"id": out[0].ID,
		"email": &user.Email,
		"preferred_name": out[0].ProviderData["username"],
		"identity_data":  out[0].ProviderData,
		"discord_id":    out[0].ProviderID,
	}).Execute(&profile)

	if profileErr != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  profileErr.Error(),
		})
	}

	if redirect != "" {
		TokenString, _ := utils.CreateToken(user.UserID, out[0])

		ctx.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    TokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			Domain:   domain,
			HTTPOnly: true,
			Secure:   os.Getenv("ENV") == "production",
		})

		ctx.ClearCookie("redirect")
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: out,
		Code:   http.StatusOK,
	})
}

func (p DiscordProvider) UpdateUser() error {
	ctx, database, user, redirect, domain := p.ctx, p.database, p.user, p.redirect, p.domain
	var out []utils.IProvider

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
	}).Eq("provider_id", user.UserID).Eq("type", user.Provider).Execute(&out)

	if insertErr != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  insertErr.Error(),
		})
	}

	var profile []utils.IProfile

	log.Println("making profile", out)

	profileErr := database.DB.From("profiles").Update(map[string]interface{}{
		"preferred_name": out[0].ProviderData["username"],
		"identity_data":  out[0].ProviderData,
	}).Eq("id", *out[0].ID).Execute(&profile)

	if profileErr != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  profileErr.Error(),
		})
	}

	if redirect != "" {
		TokenString, _ := utils.CreateToken(user.UserID, out[0])

		ctx.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    TokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			Domain:   domain,
			HTTPOnly: true,
			Secure:   os.Getenv("ENV") == "production",
		})

		ctx.ClearCookie("redirect")
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: out,
		Code:   http.StatusOK,
	})
}
