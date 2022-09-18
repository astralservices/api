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
	workspaceRouter.Put("/", UpdateWorkspace)
	workspaceRouter.Post("/", UpdateWorkspace)
	// workspaceRouter.Delete("/:id", DeleteWorkspace)

	workspaceRouter.Get("/members", GetWorkspaceMembers)
	workspaceRouter.Post("/members", AddWorkspaceMember)
	workspaceRouter.Get("/members/:member", GetWorkspaceMember)
	workspaceRouter.Put("/members/:member", UpdateWorkspaceMember)
	workspaceRouter.Delete("/members/:member", RemoveWorkspaceMember)
	workspaceRouter.Post("/members/:member/remove", RemoveWorkspaceMember) // Fallback for HTML Forms

	// compatablity with HTML forms
	workspaceRouter.Post("/bot/create", CreateWorkspaceBot)

	botRouter := workspaceRouter.Group("/bot").Use(utils.BotMiddleware)
	botRouter.Get("/", GetWorkspaceBot)
	botRouter.Post("/", UpdateWorkspaceBot)

	workspaceRouter.Get("/analytics", GetWorkspaceAnalytics)

	workspaceRouter.Get("/integrations", GetWorkspaceIntegrations)

	workspaceRouter.Post("/integrations/enable/:integrationId", EnableWorkspaceIntegration)
	workspaceRouter.Post("/integrations/disable/:integrationId", DisableWorkspaceIntegration)

	integrationRouter := workspaceRouter.Group("/integrations/:integrationId").Use(utils.WorkspaceIntegrationMiddleware, utils.BotMiddleware)
	integrationRouter.Get("/", GetWorkspaceIntegration)
	integrationRouter.Post("/", UpdateWorkspaceIntegration)
	// integrationRouter.Delete("/", DeleteWorkspaceIntegration)

	integrationRouter.Get("/data", GetIntegrationData)
	// integrationRouter.Post("/data", UpdateIntegrationData)

	integrationRouter.Get("/data/@me", GetIntegrationDataForUser)
	integrationRouter.Post("/data/@me", UpdateIntegrationDataForUser)
}
