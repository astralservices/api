package v1

import (
	"net/http"

	"github.com/astralservices/api/api/v1/auth"
	"github.com/astralservices/api/api/v1/workspaces"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gofiber/fiber/v2"
)

func V1Handler(router fiber.Router) {
	router.Get("/stats", StatsHandler)
	router.Get("/regions", RegionsHandler)
	router.Get("/team", TeamHandler)

	auth.AuthHandler(router.Group("/auth").Use(utils.AuthInjectorMiddleware))
	workspaces.WorkspacesHandler(router.Group("/workspaces"))
}


func StatsHandler(c *fiber.Ctx) error {
	var stats []utils.IStatistic

	database := db.New()

	err := database.DB.From("stats").Select("*").Execute(&stats)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	return c.JSON(utils.Response[[]utils.IStatistic]{
		Result: stats,
		Code:   http.StatusOK,
	})
}

func RegionsHandler(c *fiber.Ctx) error {
	var regions []utils.IRegion

	database := db.New()

	err := database.DB.From("regions").Select("id, flag, city, region, country, prettyName, lat, long, maxBots").Execute(&regions)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	// remove the region "localhost"
	for i, region := range regions {
		if region.ID == "localhost" {
			regions = append(regions[:i], regions[i+1:]...)
			break
		}
	}

	return c.JSON(utils.Response[[]utils.IRegion]{
		Result: regions,
		Code:   http.StatusOK,
	})
}

func TeamHandler(c *fiber.Ctx) error {
	var teams []any

	database := db.New()

	err := database.DB.From("teamMembers").Select("*, user(identity_data)").Execute(&teams)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	return c.JSON(utils.Response[any]{
		Result: teams,
		Code:   http.StatusOK,
	})
}
