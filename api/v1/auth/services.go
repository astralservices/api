package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/astralservices/api/api/v1/auth/providers/discord"
	"github.com/astralservices/api/api/v1/auth/providers/lastfm"
	"github.com/astralservices/api/api/v1/auth/providers/roblox"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
)

func CallbackHandler(ctx *fiber.Ctx) error {
	authErr := ctx.Query("error")

	if authErr != "" {
		return ctx.Status(400).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusBadRequest,
			Error:  authErr,
		})
	}

	database := db.New()

	provider := ctx.Params("provider")

	redirect := ctx.Cookies("redirect")

	if provider == "roblox" {
		rbx := roblox.New(ctx, database, redirect)
		return rbx.VerifyUser()
	}

	user, err := goth_fiber.CompleteUserAuth(ctx)

	if err != nil {
		log.Fatal(err)
	}

	var providers []utils.IProvider

	discordUser := ctx.Locals("user")

	log.Print(discordUser)

	err = database.DB.From("providers").Select("*").Eq("provider_id", user.UserID).Eq("type", user.Provider).Execute(&providers)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	var domain string

	if os.Getenv("ENV") == "development" {
		domain = "localhost"
	} else {
		domain = ctx.BaseURL()
	}

	if user.Provider != "discord" && discordUser == nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	discordProvider := discord.New(ctx, database, user, redirect, domain)
	lastfmProvider := lastfm.New(ctx, database, user, redirect, domain)

	if len(providers) == 0 {
		switch user.Provider {
		case "discord":
			return discordProvider.CreateUser()

		case "lastfm":
			return lastfmProvider.CreateUser()

		default:
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  "Provider not supported",
			})
		}
	} else {
		switch user.Provider {
		case "discord":
			return discordProvider.UpdateUser()

		case "lastfm":
			return lastfmProvider.UpdateUser()

		default:
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  "Provider not supported",
			})
		}
	}

}

func LoginHandler(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")

	if provider == "roblox" {
		// TODO: add custom login handler for roblox
	}

	if gothUser, err := goth_fiber.CompleteUserAuth(ctx); err == nil {
		return ctx.SendString(gothUser.Email)
	} else {
		return goth_fiber.BeginAuthHandler(ctx)
	}
}

func LogoutHandler(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")

	redirect := ctx.Query("redirect")

	if provider == "discord" {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		// clear cookie didnt work for some reason
		ctx.Cookie(&fiber.Cookie{
			Name:  "token",
			Value: "",
		})

		if redirect != "" {
			return ctx.Redirect(redirect)
		} else {
			return ctx.Status(200).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusOK,
			})
		}
	}

	discordUser := ctx.Locals("user").(utils.IProvider)

	database := db.New()

	var deleted any

	err := database.DB.From("providers").Delete().Eq("user", *discordUser.ID).Eq("type", provider).Execute(&deleted)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	} else {
		return ctx.Status(200).JSON(utils.Response[any]{
			Result: deleted,
			Code:   http.StatusOK,
		})
	}
}

func SessionHandler(ctx *fiber.Ctx) error {
	token := ctx.Cookies("token")
	claims, err := utils.GetClaimsFromToken(token)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	return ctx.Status(200).JSON(utils.Response[utils.IProvider]{
		Result: claims.UserInfo,
		Code:   http.StatusOK,
	})
}

func ProviderHandler(ctx *fiber.Ctx) error {
	providerId := ctx.Params("provider")

	profile := ctx.Locals("profile").(utils.IProfile)

	var providers []utils.IProvider

	database := db.New()

	err := database.DB.From("providers").Select("*").Eq("user", profile.ID).Eq("type", providerId).Execute(&providers)

	if len(providers) == 0 {
		return ctx.JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusOK,
		})
	}

	var provider utils.IProvider = providers[0]

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	return ctx.JSON(utils.Response[utils.IProvider]{
		Result: provider,
		Code:   http.StatusOK,
	})
}

func ProvidersHandler(ctx *fiber.Ctx) error {
	profile := ctx.Locals("profile").(utils.IProfile)

	var providers []utils.IProvider

	database := db.New()

	err := database.DB.From("providers").Select("*").Eq("user", profile.ID).Execute(&providers)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	return ctx.Status(200).JSON(utils.Response[[]utils.IProvider]{
		Result: providers,
		Code:   http.StatusOK,
	})
}

type StatusResponse struct {
	Authenticated bool              `json:"authenticated"`
	Blacklist     *utils.IBlacklist `json:"blacklist,omitempty"`
}

func StatusHandler(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(utils.IProvider)
	var blacklist []utils.IBlacklist

	database := db.New()

	err := database.DB.From("blacklist").Select("*").Eq("user", *user.ID).Execute(&blacklist)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())

	}

	if len(blacklist) == 0 {
		return ctx.Status(200).JSON(utils.Response[StatusResponse]{
			Result: StatusResponse{
				Authenticated: true,
			},
			Code: http.StatusOK,
		})
	}

	return ctx.Status(200).JSON(utils.Response[StatusResponse]{
		Result: StatusResponse{
			Authenticated: true,
			Blacklist:     &blacklist[0],
		},
		Code: http.StatusForbidden,
	})
}

// gets all the user's data and returns it as an actual JSON file
func DataHandler(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(utils.IProvider)

	var providers []utils.IProvider

	database := db.New()

	err := database.DB.From("providers").Select("*").Eq("user", *user.ID).Execute(&providers)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	var blacklist []utils.IBlacklist

	err = database.DB.From("blacklist").Select("*").Eq("user", *user.ID).Execute(&blacklist)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	var profile utils.IProfile

	err = database.DB.From("profiles").Select("*").Single().Eq("id", *user.ID).Execute(&profile)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	var moderationActions []utils.IBotModerationAction

	err = database.DB.From("moderation_actions").Select("*").Eq("user", *user.ID).Execute(&moderationActions)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	var workspaceMemberships []utils.IWorkspaceMemberWithoutProfile

	err = database.DB.From("workspace_members").Select("*, workspace(*)").Eq("profile", *user.ID).Execute(&workspaceMemberships)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	var workspaces []any

	for _, workspaceMember := range workspaceMemberships {
		workspaces = append(workspaces, workspaceMember.Workspace)
	}

	var bots []utils.IBot

	err = database.DB.From("bots").Select("id, created_at, owner, region, settings, token, commands").Eq("owner", *user.ID).Execute(&bots)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	type FinalData struct {
		AuthProviders        []utils.IProvider                      `json:"auth_providers"`
		Blacklist            []utils.IBlacklist                     `json:"blacklist"`
		Profile              utils.IProfile                         `json:"profile"`
		BotModerationActions []utils.IBotModerationAction           `json:"bot_moderation_actions"`
		WorkspaceMemberships []utils.IWorkspaceMemberWithoutProfile `json:"workspace_memberships"`
		Workspaces           []any                                  `json:"workspaces"`
		Bots                 []utils.IBot                           `json:"bots"`
	}

	var finalData FinalData = FinalData{
		AuthProviders:        providers,
		Blacklist:            blacklist,
		Profile:              profile,
		BotModerationActions: moderationActions,
		WorkspaceMemberships: workspaceMemberships,
		Workspaces:           workspaces,
		Bots:                 bots,
	}

	// convert finalData to byte and return it
	data, err := json.Marshal(finalData)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err.Error())
	}

	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Response().Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", profile.PreferredName))
	return ctx.Status(200).Send(data)
}
