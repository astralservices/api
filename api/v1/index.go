package v1

import (
	"errors"
	"net/http"
	"os"
	"sort"

	"github.com/astralservices/api/api/v1/auth"
	"github.com/astralservices/api/api/v1/workspaces"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/aybabtme/orderedjson"
	"github.com/gofiber/fiber/v2"
)

func V1Handler(router fiber.Router) {
	router.Get("/stats", StatsHandler)
	router.Get("/regions", RegionsHandler)
	router.Get("/team", TeamHandler)
	router.Get("/plans", PlansHandler)
	router.Get("/integrations", IntegrationsHandler)
	router.Get("/integrations/:id", IntegrationHandler)

	auth.AuthHandler(router.Group("/auth").Use(utils.AuthInjectorMiddleware))
	workspaces.WorkspacesHandler(router.Group("/workspaces"))
}

func PlansHandler(c *fiber.Ctx) error {
	var plans []utils.IPlan

	database := db.New()

	err := database.DB.From("plans").Select("*").Execute(&plans)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	sort.Slice(plans, func(i, j int) bool {
		return plans[i].PriceMonthly < plans[j].PriceMonthly
	})

	return c.JSON(utils.Response[[]utils.IPlan]{
		Result: plans,
		Code:   http.StatusOK,
	})
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
	var regions []*utils.IRegion

	database := db.New()

	err := database.DB.From("regions").Select("id, flag, city, region, country, prettyName, lat, long, maxBots").Execute(&regions)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	if os.Getenv("ENV") != "development" {
		// remove the region "localhost"
		for i, region := range regions {
			if region.ID == "localhost" {
				regions = append(regions[:i], regions[i+1:]...)
				break
			}
		}

	}

	var bots []struct {
		Region string `json:"region"`
	}

	err = database.DB.From("bots").Select("region").Execute(&bots)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	// attach the number of bots to each region
	for _, region := range regions {
		region.Bots = 0
		for _, bot := range bots {
			if bot.Region == region.ID {
				region.Bots++
			}
		}
	}

	return c.JSON(utils.Response[[]*utils.IRegion]{
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

func IntegrationsHandler(c *fiber.Ctx) error {
	var integrations []any

	database := db.New()

	err := database.DB.From("integrations").Select("*").Execute(&integrations)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	return c.JSON(utils.Response[any]{
		Result: integrations,
		Code:   http.StatusOK,
	})
}

func IntegrationHandler(c *fiber.Ctx) error {
	var integrations []orderedjson.Map
	var integration any
	id := c.Params("id")

	database := db.New()

	err := database.DB.From("integrations").Select("*").Eq("id", id).Execute(&integrations)

	if err != nil {
		return c.JSON(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
	}

	if len(integrations) == 0 {
		return utils.ErrorResponse(c, http.StatusNotFound, errors.New("Integration not found"), true)
	}

	integration = integrations[0]

	return c.JSON(utils.Response[any]{
		Result: integration,
		Code:   http.StatusOK,
	})
}
