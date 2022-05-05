package utils

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "production" {
			handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}), handlers.AllowedOrigins([]string{"*", "http://localhost:3000", "http://localhost:8000", "https://*.astralapp.io", "https://astralapp.io"}), handlers.AllowCredentials())(h)
		}

		h.ServeHTTP(w, r)
	})
}

func JSONMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enabled := false

		if enabled {
			h.ServeHTTP(w, r)
		} else {
			data, err := json.Marshal(Response[any]{
				Error: "Authentication is disabled",
				Code:  http.StatusForbidden,
			})

			if err != nil {
				w.Write([]byte("Error"))
				return
			}

			w.WriteHeader(http.StatusUnauthorized)

			w.Write(data)
		}
	})
}
