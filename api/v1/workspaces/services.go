package workspaces

import (
	"image"
	"image/png"
	"net/http"
	"os"
	"sort"

	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	client := fiber.AcquireClient()

	plan := utils.IPlan{}

	err = database.DB.From("plans").Select("*").Single().Eq("id", workspaceData.Plan).Execute(&plan)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	file, err := fileHeader.Open()

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	// resize the image

	img, _, err := image.Decode(file)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	resized := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	err = png.Encode(agent.Request().BodyWriter(), resized)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	req, res := agent.Request(), fiber.AcquireResponse()

	err = agent.Do(req, res)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	// update the workspace with the path to the icon

	updatedWorkspaces := []utils.IWorkspace{}

	err = database.DB.From("workspaces").Update(map[string]interface{}{
		"logo": publicPath,
	}).Eq("id", *workspace.ID).Execute(&updatedWorkspaces)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
	// workspaceID := ctx.Params("id")

	// return not implemented
	return ctx.Status(501).JSON(utils.Response[any]{
		Result: nil,
		Code:   http.StatusNotImplemented,
	})

	// workspaces := []utils.IWorkspace{}

	// database := db.New()
}

func GetWorkspaceMembers(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)
	database := db.New()

	countOnly := ctx.Query("count") == "true"

	var workspace_members []any

	err := database.DB.From("workspace_members").Select("*, profile(*)").Eq("workspace", *workspace.ID).Execute(&workspace_members)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		if redirect != "" {
			return ctx.Redirect(redirect + "?error=" + err.Error())
		}

		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	// check if the user is already a member of the workspace

	err = database.DB.From("workspace_members").Select("*").Eq("workspace", *workspace.ID).Eq("profile", member_profile.ID).Execute(&workspace_membership)

	if err != nil {
		if redirect != "" {
			return ctx.Redirect(redirect + "?error=User is already a member of this workspace")
		}

		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	err = database.DB.From("workspace_members").Insert(map[string]interface{}{
		"workspace":  workspace.ID,
		"profile":    member_profile.ID,
		"role":       ctx.FormValue("role"),
		"invited_by": self_member.ID,
		"pending":    true,
	}).Execute(&workspace_membership)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		if redirect != "" {
			return ctx.Redirect(redirect + "?error=" + err.Error())
		}

		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if len(bots) == 0 {
		return ctx.Status(200).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusOK,
		})
	}

	bot := bots[0]

	var analytics []utils.IBotAnalytics

	err = database.DB.From("bot_analytics").Select("commands, timestamp, members, messages").Eq("bot", bot.ID).Execute(&analytics)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
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
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	var bots []utils.IBot

	database := db.New()

	err := database.DB.From("bots").Select("id, created_at, owner, region, settings, token, commands").Eq("workspace", *workspace.ID).Execute(&bots)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if len(bots) == 0 {
		return ctx.Status(404).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  "Bot not found",
		})
	}

	bot := bots[0]

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: bot,
		Code:   http.StatusOK,
	})
}

func CreateWorkspaceBot(ctx *fiber.Ctx) error {
	workspace := ctx.Locals("workspace").(utils.IWorkspace)

	redirect := ctx.FormValue("redirect")

	var bot utils.IBot

	database := db.New()

	err := database.DB.From("bots").Insert(map[string]interface{}{
		"workspace": workspace.ID,
		"name":      ctx.FormValue("name"),
		"token":     ctx.FormValue("token"),
	}).Execute(&bot)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	if redirect != "" {
		return ctx.Redirect(redirect)
	}

	return ctx.Status(200).JSON(utils.Response[any]{
		Result: bot,
		Code:   http.StatusOK,
	})
}
