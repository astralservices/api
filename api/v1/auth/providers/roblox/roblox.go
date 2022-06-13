package roblox

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/astralservices/api/utils"
	"github.com/astralservices/goblox/goblox"
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
)

type RobloxProvider struct {
	ctx      *fiber.Ctx
	database *supabase.Client
	redirect string
	domain   string
	roblox *goblox.Client
}

func New(c *fiber.Ctx, database *supabase.Client, redirect string) RobloxProvider {
	roblox := goblox.New()
	provider := RobloxProvider{
		ctx:      c,
		database: database,
		redirect: redirect,
		roblox: roblox,
	}
	return provider
}

func (p RobloxProvider) GenerateCodeForUser() error {
	ctx, database, redirect, roblox := p.ctx, p.database, p.redirect, p.roblox

	userName := ctx.Query("username")

	user, userErr := roblox.Users.GetUserByUsername(userName)

	if userErr != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  userErr.Error(),
		})
	}

	if user.ID == 0 {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  "User not found",
		})
	}

	discordUser := ctx.Locals("user").(utils.IProvider)

	var out []utils.IProvider

	// generate 5 word code using a for statement separated by a space
	var codes []string
	for i := 0; i < 5; i++ {
		codes = append(codes, utils.RandomWord())
	}

	code := strings.Join(codes, " ")

	insertErr := database.DB.From("providers").Insert(map[string]interface{}{
		"type":                   "roblox",
		"user":                   discordUser.ID,
		"provider_id":            user.ID,
		"provider_data": 		  map[string]interface{}{
			"status": "pending",
			"code": code,
			"username": user.Name,
		},
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
		return ctx.Redirect(redirect+"?code="+code)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: out,
		Code:   http.StatusOK,
	})
}

func (p RobloxProvider) VerifyUser() error {
	ctx, database, redirect, roblox := p.ctx, p.database, p.redirect, p.roblox

	code := ctx.Query("code")

	var out []utils.IProvider

	err := database.DB.From("providers").Select("*").Eq("provider_data->>code", code).Execute(&out)

	if err != nil {
		if redirect != "" {
			ctx.ClearCookie("redirect")
			return ctx.Redirect(redirect+"?error="+err.Error())
		}
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if len(out) == 0 {
		if redirect != "" {
			ctx.ClearCookie("redirect")
			return ctx.Redirect(redirect+"?error=Code not found")
		}
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  "Code not found",
		})
	}

	if out[0].ProviderData["status"] == "pending" {
		id, err := strconv.ParseInt(out[0].ProviderID, 10, 64)
		if err != nil {
			if redirect != "" {
				ctx.ClearCookie("redirect")
				return ctx.Redirect(redirect+"/error?error="+err.Error())
			}
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  err.Error(),
			})
		}
		user, err := roblox.Users.GetUserById(id)

		if err != nil {
			if redirect != "" {
				ctx.ClearCookie("redirect")
				return ctx.Redirect(redirect+"/error?error="+err.Error())
			}
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  err.Error(),
			})
		}

		if user.ID == 0 {
			if redirect != "" {
				ctx.ClearCookie("redirect")
				return ctx.Redirect(redirect+"/error?error=User not found")
			}
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  "User not found",
			})
		}

		authCode := ctx.Query("code")

		if !strings.Contains(user.Description, authCode) {
			if redirect != "" {
				ctx.ClearCookie("redirect")
				return ctx.Redirect(redirect+"/error?error=Invalid Code")
			}
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  "Invalid code",
			})
		}

		out[0].ProviderData["status"] = "verified"
		err = database.DB.From("providers").Update(map[string]interface{}{
			"provider_data": out[0].ProviderData,
		}).Eq("id", *out[0].ID).Execute(&out)

		if err != nil {
			return ctx.Status(500).JSON(utils.Response[any]{
				Result: nil,
				Code:   http.StatusInternalServerError,
				Error:  err.Error(),
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

	if redirect != "" {
		ctx.ClearCookie("redirect")
		return ctx.Redirect(redirect+"/error?error=Code already verified")
	}

	return ctx.Status(500).JSON(utils.Response[any]{
		Result: nil,
		Code:   http.StatusInternalServerError,
		Error:  "Code already verified",
	})
}