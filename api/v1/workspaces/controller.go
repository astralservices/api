package workspaces

import (
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
)

func WorkspacesHandler(router fiber.Router) {
	authed := router.Use(utils.AuthMiddleware, utils.ProfileMiddleware)
	authed.Get("/", GetWorkspaces)
	authed.Post("/", CreateWorkspace)

	workspaceRouter := authed.Group("/:workspace_id").Use(utils.WorkspaceMiddleware, utils.WorkspaceMemberMiddleware)

	workspaceRouter.Get("/", GetWorkspace)
	// workspaceRouter.Put("/:id", UpdateWorkspace)
	// workspaceRouter.Delete("/:id", DeleteWorkspace)

	workspaceRouter.Get("/members", GetWorkspaceMembers)
	workspaceRouter.Post("/members", AddWorkspaceMember)
	workspaceRouter.Get("/members/:member", GetWorkspaceMember)
	workspaceRouter.Put("/members/:member", UpdateWorkspaceMember)
	workspaceRouter.Delete("/members/:member", RemoveWorkspaceMember)
	workspaceRouter.Post("/members/:member/remove", RemoveWorkspaceMember) // Fallback for HTML Forms

	workspaceRouter.Get("/bot", GetWorkspaceBot)

	workspaceRouter.Get("/analytics", GetWorkspaceAnalytics)
}
