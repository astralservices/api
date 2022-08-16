package main

import (
	"context"
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	v1 "github.com/astralservices/api/api/v1"
	"github.com/astralservices/api/api/v1/auth"
	_ "github.com/astralservices/api/docs"
	"github.com/astralservices/api/utils"
	"github.com/getsentry/sentry-go"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func IndexHandler(c *fiber.Ctx) error {
	return c.JSON(utils.Response[struct {
		Message string `json:"message"`
	}]{
		Result: struct {
			Message string "json:\"message\""
		}{Message: "API is running!"},
		Code: http.StatusOK,
	})
}

func main() {
	godotenv.Load(".env.local")
	rand.Seed(time.Now().UnixNano())
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	app := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Astral Services API",
		AppName:       "Astral Services API",
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-type", "application/json; charset=utf-8")

		return c.Next()
	})

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/monitor", monitor.New(monitor.Config{Title: "Astral API Metrics Page", Refresh: time.Second * 5}))

	app.Get("/", IndexHandler)

	api := app.Group("/api")

	auth.InitGoth()

	v1.V1Handler(api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	sentry.Init(sentry.ClientOptions{
		Dsn: "https://6fe272ef81454aec990e3f526f51dd7f@gt.astralapp.io/1",
	})

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Infoln("Starting server on port " + port)
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	_, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	app.Shutdown()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	sentry.Flush(time.Second * 5)
	log.Println("Shutting down")
	os.Exit(0)
}
