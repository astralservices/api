package workspaces

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"sort"
	"strconv"

	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
	"github.com/nqd/flat"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

func GetWorkspaces(ctx *fiber.Ctx) error {
	database := db.New()

	user := ctx.Locals("user").(utils.IProvider)
	// typing just does not like to work here
	var workspace_memberships []any

	err := database.DB.From("workspace_members").Select("workspace(*)").Eq("profile", *user.ID).Execute(&workspace_memberships)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	return ctx.Status(200).JSON(utils.Response[[]any]{
		Result: workspace_memberships,
		Code:   http.StatusOK,
	})
}

func CreateWorkspace(ctx *fiber.Ctx) error {
	database := db.New()

	user := ctx.Locals("user").(utils.IProvider)
	profile := ctx.Locals("profile").(utils.IProfile)

	workspaceData := struct {
		Name        string `json:"name" form:"name"`
		Description string `json:"description" form:"description"`
		Visibility  string `json:"visibility" form:"visibility"`
		Plan        string `json:"plan" form:"plan"`
		Redirect    string `json:"redirect" form:"redirect"`
	}{}

	err := ctx.BodyParser(&workspaceData)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	client := fiber.AcquireClient()

	plan := utils.IPlan{}

	err = database.DB.From("plans").Select("*").Single().Eq("id", workspaceData.Plan).Execute(&plan)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	plans := map[string]int{
		"free":    1,
		"starter": 2,
		"pro":     3,
	}

	// create the stripe subscription

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	stripeParams := &stripe.SubscriptionParams{
		Customer: stripe.String(profile.StripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			&stripe.SubscriptionItemsParams{
				Price: stripe.String(plan.ID),
			},
		},
	}

	subscription, err := sub.New(stripeParams)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	// create the workspace

	workspaces := []utils.IWorkspace{}

	err = database.DB.From("workspaces").Insert(map[string]interface{}{
		"name":       workspaceData.Name,
		"visibility": workspaceData.Visibility,
		"plan":       plans[workspaceData.Plan],
		"owner":      *user.ID,
		"settings": map[string]interface{}{
			"isPaidPlan":  plans[workspaceData.Plan] > 1,
			"description": workspaceData.Description,
			"stripe": map[string]interface{}{
				"subscription": subscription.ID,
			},
		},
	}).Execute(&workspaces)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	workspace := workspaces[0]

	// upload the workspace logo

	path := os.Getenv("SUPABASE_URL") + "/storage/v1/object/workspaces-data/workspaces/" + *workspace.ID + "/logo.png"
	publicPath := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/workspaces-data/workspaces/" + *workspace.ID + "/logo.png"

	agent := client.Post(path)

	agent.Add("Content-Type", "image/png")
	agent.Add("Authorization", "Bearer "+os.Getenv("SUPABASE_KEY"))

	fileHeader, err := ctx.FormFile("icon")

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	file, err := fileHeader.Open()

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	// resize the image

	img, _, err := image.Decode(file)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	resized := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	err = png.Encode(agent.Request().BodyWriter(), resized)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	req, res := agent.Request(), fiber.AcquireResponse()

	err = agent.Do(req, res)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	// update the workspace with the path to the icon

	updatedWorkspaces := []utils.IWorkspace{}

	err = database.DB.From("workspaces").Update(map[string]interface{}{
		"logo": publicPath,
	}).Eq("id", *workspace.ID).Execute(&updatedWorkspaces)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	updatedWorkspace := updatedWorkspaces[0]

	// create the workspace member

	workspace_membership := utils.IWorkspaceMember{}

	database.DB.From("workspace_members").Insert(map[string]interface{}{
		"workspace": workspace.ID,
		"profile":   *user.ID,
		"role":      "owner",
	}).Execute(&workspace_membership)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	redirect := workspaceData.Redirect

	if redirect != "" {
		return ctx.Redirect(redirect + "workspaces/" + *workspace.ID)
	}

	return ctx.Status(200).JSON(utils.Response[utils.IWorkspace]{
		Result: updatedWorkspace,
		Code:   http.StatusOK,
	})
}

func GetWorkspace(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	return ctx.Status(200).JSON(utils.Response[utils.IWorkspace]{
		Result: workspace,
		Code:   http.StatusOK,
	})
}

func UpdateWorkspace(ctx *fiber.Ctx) error {
	database := db.New()

	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	workspaceData := struct {
		Name        string `json:"name" form:"name"`
		Description string `json:"description" form:"description"`
		Visibility  string `json:"visibility" form:"visibility"`
		Plan        string `json:"plan" form:"plan"`
		Redirect    string `json:"redirect" form:"redirect"`
	}{}

	err := ctx.BodyParser(&workspaceData)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	client := fiber.AcquireClient()

	plan := utils.IPlan{}

	err = database.DB.From("plans").Select("*").Single().Eq("id", workspaceData.Plan).Execute(&plan)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	plans := map[string]int{
		"free":    1,
		"starter": 2,
		"pro":     3,
	}

	// create the workspace

	workspaces := []utils.IWorkspace{}

	err = database.DB.From("workspaces").Update(map[string]interface{}{
		"name":       workspaceData.Name,
		"visibility": workspaceData.Visibility,
		"plan":       plans[workspaceData.Plan],
		"settings": map[string]interface{}{
			"isPaidPlan":  plans[workspaceData.Plan] > 1,
			"description": workspaceData.Description,
		},
	}).Eq("id", *workspace.ID).Execute(&workspaces)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	workspace = workspaces[0]

	mf, err := ctx.MultipartForm()

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	iconExists := mf.File["icon"]

	if iconExists != nil {

		icon, err := ctx.FormFile("icon")

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}

		// upload the workspace logo

		path := os.Getenv("SUPABASE_URL") + "/storage/v1/object/workspaces-data/workspaces/" + *workspace.ID + "/logo.png"
		publicPath := os.Getenv("SUPABASE_URL") + "/storage/v1/object/public/workspaces-data/workspaces/" + *workspace.ID + "/logo.png"

		agent := client.Put(path)

		agent.Add("Content-Type", "image/png")
		agent.Add("Authorization", "Bearer "+os.Getenv("SUPABASE_KEY"))

		file, err := icon.Open()

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}

		// resize the image

		img, _, err := image.Decode(file)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}

		resized := resize.Thumbnail(200, 200, img, resize.Lanczos3)

		err = png.Encode(agent.Request().BodyWriter(), resized)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}

		req, res := agent.Request(), fiber.AcquireResponse()

		err = agent.Do(req, res)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}

		// update the workspace with the path to the icon

		updatedWorkspaces := []utils.IWorkspace{}

		err = database.DB.From("workspaces").Update(map[string]interface{}{
			"logo": publicPath,
		}).Eq("id", *workspace.ID).Execute(&updatedWorkspaces)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}

		workspace = updatedWorkspaces[0]

	}

	redirect := workspaceData.Redirect

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[utils.IWorkspace]{
		Result: workspace,
		Code:   http.StatusOK,
	})
}

func GetWorkspaceMembers(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	database := db.New()

	countOnly := ctx.Query("count") == "true"

	var workspace_members []any

	err := database.DB.From("workspace_members").Select("*, profile(*)").Eq("workspace", *workspace.ID).Execute(&workspace_members)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if countOnly {
		return ctx.Status(200).JSON(utils.Response[any]{
			Result: len(workspace_members),
			Code:   http.StatusOK,
		})
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: workspace_members,
		Code:   http.StatusOK,
	})
}

func AddWorkspaceMember(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	self_member := ctx.Locals("workspace_member").(utils.IWorkspaceMember)

	redirect := ctx.FormValue("redirect")

	var workspace_membership []any
	member_profile := utils.IProfile{}

	database := db.New()

	err := database.DB.From("profiles").Select("*").Single().Eq("discord_id", ctx.FormValue("discord")).Execute(&member_profile)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	// check if the user is already a member of the workspace

	err = database.DB.From("workspace_members").Select("*").Eq("workspace", *workspace.ID).Eq("profile", member_profile.ID).Execute(&workspace_membership)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, errors.New("User is already a member of this workspace"), true)
	}

	err = database.DB.From("workspace_members").Insert(map[string]interface{}{
		"workspace":  workspace.ID,
		"profile":    member_profile.ID,
		"role":       ctx.FormValue("role"),
		"invited_by": self_member.ID,
		"pending":    true,
	}).Execute(&workspace_membership)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: workspace_membership[0],
		Code:   http.StatusOK,
	})
}

func GetWorkspaceMember(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	var workspace_members []any

	database := db.New()

	err := database.DB.From("workspace_members").Select("*, profile(*)").Eq("workspace", *workspace.ID).Eq("profile", ctx.Params("member")).Execute(&workspace_members)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if len(workspace_members) == 0 {
		return ctx.Status(404).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  "Workspace member not found",
		})
	}

	workspace_member := workspace_members[0]

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: workspace_member,
		Code:   http.StatusOK,
	})
}

func UpdateWorkspaceMember(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	redirect := ctx.FormValue("redirect")

	var workspace_members []any

	database := db.New()

	err := database.DB.From("workspace_members").Update(map[string]interface{}{
		"role": ctx.FormValue("role"),
	}).Eq("workspace", *workspace.ID).Eq("profile", ctx.Params("member")).Execute(&workspace_members)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: workspace_members,
		Code:   http.StatusOK,
	})
}

func RemoveWorkspaceMember(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	redirect := ctx.FormValue("redirect")

	var workspace_members []any

	database := db.New()

	err := database.DB.From("workspace_members").Delete().Eq("workspace", *workspace.ID).Eq("id", ctx.Params("member")).Execute(&workspace_members)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: workspace_members,
		Code:   http.StatusOK,
	})
}

func GetWorkspaceAnalytics(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	database := db.New()

	var bots []utils.IBot

	err := database.DB.From("bots").Select("id").Eq("workspace", *workspace.ID).Execute(&bots)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if len(bots) == 0 {
		return ctx.Status(200).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusOK,
		})
	}

	bot := bots[0]

	var analytics []utils.IBotAnalytics

	err = database.DB.From("bot_analytics").Select("commands, timestamp, members, messages").Eq("bot", *bot.ID).Execute(&analytics)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	// sort array by timestamp, latest first
	sort.Slice(analytics, func(i, j int) bool {
		return analytics[i].Timestamp.UTC().After(analytics[j].Timestamp.UTC())
	})

	return ctx.Status(200).JSON(utils.Response[[]utils.IBotAnalytics]{
		Result: analytics,
		Code:   http.StatusOK,
	})
}

func GetWorkspaceBot(ctx *fiber.Ctx) error {
	bot := ctx.Locals("bot").(utils.IBot)

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: bot,
		Code:   http.StatusOK,
	})
}

type BotSettings struct {
	Guild               string               `json:"guild,omitempty"`
	Prefix              string               `json:"prefix,omitempty" form:"prefix,omitempty"`
	Status              string               `json:"status,omitempty" form:"status,omitempty"`
	Activities          []utils.IBotActivity `json:"activities,omitempty" form:"activities,omitempty"`
	RandomizeActivities bool                 `json:"randomizeActivities,omitempty" form:"randomizeActivities,omitempty"`
	ActivityInterval    int                  `json:"activityInterval,omitempty" form:"activityInterval,omitempty"`
	CurrentActivity     int                  `json:"currentActivity,omitempty"`
	Modules             utils.IBotModules    `json:"modules" form:"modules"`
}

type BotFormData struct {
	Region      *string                `json:"region,omitempty" form:"region,omitempty"`
	Settings    *BotSettings           `json:"settings,omitempty" form:"settings,omitempty"`
	Permissions *utils.IBotPermissions `json:"permissions,omitempty" form:"permissions,omitempty"`
	Token       *string                `json:"token,omitempty" form:"token,omitempty"`
}

func CreateWorkspaceBot(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	user := ctx.Locals("user").(utils.IProvider)

	redirect := ctx.FormValue("redirect")

	var formData BotFormData = BotFormData{
		Settings: &BotSettings{
			Activities: []utils.IBotActivity{
				{
					Type: "PLAYING",
					Name: "a game!",
				},
			},
			RandomizeActivities: false,
			ActivityInterval:    300,
			CurrentActivity:     0,
			Modules:             utils.IBotModules{},
			Status:              "online",
		},
	}

	err := ctx.BodyParser(&formData)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	// validate the token through Discord's API by fetching the self user

	client := fiber.AcquireClient()

	agent := client.Get("https://discord.com/api/v9/users/@me")

	agent.Add("Authorization", "Bot "+*formData.Token)

	req, res := agent.Request(), fiber.AcquireResponse()

	err = agent.Do(req, res)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if res.StatusCode() != 200 {
		return utils.ErrorResponse(ctx, 422, errors.New("Invalid token"), true)
	}

	var bots []any

	var bot any

	database := db.New()

	err = database.DB.From("bots").Insert(map[string]interface{}{
		"workspace": workspace.ID,
		"region":    formData.Region,
		"settings":  formData.Settings,
		"token":     formData.Token,
		"owner":     user.ID,
	}).Execute(&bots)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	bot = bots[0]

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: bot,
		Code:   http.StatusOK,
	})
}

func UpdateWorkspaceBot(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	bot := ctx.Locals("bot").(utils.IBot)

	redirect := ctx.FormValue("redirect")

	form := BotFormData{
		Settings: &BotSettings{
			Guild:               bot.Settings.Guild,
			Prefix:              bot.Settings.Prefix,
			Status:              bot.Settings.Status,
			Activities:          bot.Settings.Activities,
			RandomizeActivities: bot.Settings.RandomizeActivities,
			ActivityInterval:    bot.Settings.ActivityInterval,
			CurrentActivity:     bot.Settings.CurrentActivity,
			Modules:             bot.Settings.Modules,
		},
		Permissions: &bot.Permissions,
	}

	err := ctx.BodyParser(&form)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	var bots []any

	var updatedBot any

	database := db.New()

	err = database.DB.From("bots").Update(BotFormData{
		Region:      form.Region,
		Settings:    form.Settings,
		Token:       form.Token,
		Permissions: form.Permissions,
	}).Eq("workspace", *workspace.ID).Execute(&bots)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	updatedBot = bots[0]

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: updatedBot,
		Code:   http.StatusOK,
	})
}

func DeleteWorkspaceBot(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	redirect := ctx.FormValue("redirect")

	var bot utils.IBot

	database := db.New()

	err := database.DB.From("bots").Delete().Eq("workspace", *workspace.ID).Eq("id", ctx.Params("bot")).Execute(&bot)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: bot,
		Code:   http.StatusOK,
	})
}

func GetWorkspaceIntegrations(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	var integrations []utils.IWorkspaceIntegration

	database := db.New()

	err := database.DB.From("workspace_integrations").Select("*").Eq("workspace", *workspace.ID).Execute(&integrations)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: integrations,
		Code:   http.StatusOK,
	})
}

func GetWorkspaceIntegration(ctx *fiber.Ctx) error {
	integration := ctx.Locals("integration").(utils.IWorkspaceIntegration)

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: integration,
		Code:   http.StatusOK,
	})
}

func EnableWorkspaceIntegration(ctx *fiber.Ctx) error {
	// check if the workspace integration exists, if not to create it and set enabled to true, and if it does exist to just set enabled to true
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	integrationId := ctx.Params("integrationId")

	redirect := ctx.FormValue("redirect")

	database := db.New()

	var integrations []utils.IWorkspaceIntegration
	var integration utils.IWorkspaceIntegration

	err := database.DB.From("workspace_integrations").Select("*").Eq("workspace", *workspace.ID).Eq("integration", integrationId).Execute(&integrations)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if len(integrations) == 0 {
		err = database.DB.From("workspace_integrations").Insert(map[string]interface{}{
			"workspace":   *workspace.ID,
			"integration": integrationId,
			"enabled":     true,
		}).Execute(&integrations)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}
	} else {
		err = database.DB.From("workspace_integrations").Update(map[string]interface{}{
			"enabled": true,
		}).Eq("workspace", *workspace.ID).Eq("integration", integrationId).Execute(&integrations)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}
	}

	integration = integrations[0]

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: integration,
		Code:   http.StatusOK,
	})
}

func DisableWorkspaceIntegration(ctx *fiber.Ctx) error {
	// check if the workspace integration exists, if not to create it and set enabled to true, and if it does exist to just set enabled to true
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	integrationId := ctx.Params("integrationId")

	redirect := ctx.FormValue("redirect")

	database := db.New()

	var integrations []utils.IWorkspaceIntegration
	var integration utils.IWorkspaceIntegration

	err := database.DB.From("workspace_integrations").Select("*").Eq("workspace", *workspace.ID).Eq("integration", integrationId).Execute(&integrations)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if len(integrations) == 0 {
		err = database.DB.From("workspace_integrations").Insert(map[string]interface{}{
			"workspace":   *workspace.ID,
			"integration": integrationId,
			"enabled":     false,
		}).Execute(&integrations)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}
	} else {
		err = database.DB.From("workspace_integrations").Update(map[string]interface{}{
			"enabled": false,
		}).Eq("workspace", *workspace.ID).Eq("integration", integrationId).Execute(&integrations)

		if err != nil {
			return utils.ErrorResponse(ctx, 500, err, false)
		}
	}

	integration = integrations[0]

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: integration,
		Code:   http.StatusOK,
	})
}

func UpdateWorkspaceIntegration(ctx *fiber.Ctx) error {
	integration := ctx.Locals("integration").(utils.IWorkspaceIntegration)

	redirect := ctx.FormValue("redirect")

	database := db.New()

	form, err := ctx.Request().MultipartForm()

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	data := make(map[string]interface{})

	for key, value := range form.Value {
		data[key] = value[0]
	}

	out, err := flat.Unflatten(data, nil)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	fmt.Printf("%+v\n", out)

	var integrations []utils.IWorkspaceIntegration

	err = database.DB.From("workspace_integrations").Update(map[string]interface{}{
		"settings": out,
	}).Eq("id", strconv.Itoa(integration.ID)).Execute(&integrations)

	if err != nil {
		return utils.ErrorResponse(ctx, 500, err, false)
	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: integration,
		Code:   http.StatusOK,
	})
}
