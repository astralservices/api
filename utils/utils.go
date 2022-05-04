package utils

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	return handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}), handlers.AllowedOrigins([]string{"*", "http://localhost:3000", "http://localhost:8000", "https://*.astralapp.io", "https://astralapp.io"}))
}
