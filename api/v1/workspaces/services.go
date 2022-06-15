package workspaces

import (
	"net/http"

	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
)

func GetWorkspaces(ctx *fiber.Ctx) error {
	database := db.New()

	user := ctx.Locals("user").(utils.IProvider)

	workspaces := []utils.IWorkspace{}
	workspace_memberships := []utils.IWorkspaceMember{}

	err := database.DB.From("workspace_members").Select("*").Eq("profile", *user.ID).Execute(&workspace_memberships)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	workspaceIds := []string{}

	for _, workspace_membership := range workspace_memberships {
		workspaceIds = append(workspaceIds, workspace_membership.ID)
	}

	err = database.DB.From("workspaces").Select("*").In("id", workspaceIds).Execute(&workspaces)

	if err != nil {
		return ctx.Status(500).JSON(utils.Response[any]{
			Result: nil,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	return ctx.Status(200).JSON(utils.Response[[]utils.IWorkspace]{
		Result: workspaces,
		Code:   http.StatusOK,
	})
}