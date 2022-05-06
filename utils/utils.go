package utils

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"

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

type String string

func (s String) Format(data map[string]interface{}) (out string, err error) {
	t := template.Must(template.New("").Parse(string(s)))
	builder := &strings.Builder{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.String()
	return
}

func GetCallbackURL(provider string) string {
	callbackUrl := os.Getenv("CALLBACK_URL")

	s, err := String(callbackUrl).Format(map[string]interface{}{
		"Provider": provider,
	})

	if err != nil {
		return ""
	}

	return s
}
