package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	v1 "github.com/astralservices/api/api/v1"
	_ "github.com/astralservices/api/docs"
	"github.com/astralservices/api/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(utils.Response[struct {
		Message string `json:"message"`
	}]{
		Result: struct {
			Message string "json:\"message\""
		}{Message: "API is running!"},
		Code: http.StatusOK,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

// @title Astral API
// @version 1.0.0
// @description The official API for Astral Services.
// @termsOfService https://docs.astralapp.io/legal/terms

// @contact.name DevOps Team
// @contact.url https://astralapp.io
// @contact.email devops@astralapp.io

// @license.name MPL-2.0
// @license.url https://opensource.org/licenses/MPL-2.0

// @host api.astralapp.io
// @BasePath /api/v1
func main() {
	godotenv.Load(".env.local")
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.Use(utils.LoggingMiddleware)
	r.Use(utils.CORSMiddleware)
	r.Use(utils.JSONMiddleware)
	r.StrictSlash(true)

	r.HandleFunc("/", IndexHandler)
	
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	
	r.Handle("/api/v1", v1.New(r))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Infoln("Starting server on port " + port)
		if err := srv.ListenAndServe(); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
