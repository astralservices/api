package workspaces

import (
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
)

func WorkspacesHandler(router fiber.Router) {
	authed := router.Use(utils.AuthMiddleware, utils.ProfileMiddleware)
	authed.Get("/", GetWorkspaces)
}
