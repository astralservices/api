package utils

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"

	db "github.com/astralservices/api/supabase"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/nedpals/supabase-go"
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
		database := db.New()
		userCookie, err := r.Cookie("access_token")

		var res []byte

		if err != nil {
			res, err = json.Marshal(Response[struct {
				Message string `json:"message"`
			}]{
				Result: struct {
					Message string "json:\"message\""
				}{Message: "You must be logged in to access this page!"},
				Code: http.StatusUnauthorized,
			})

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(res)
		}

		user, err := database.Auth.User(r.Context(), userCookie.Value)

		if err != nil {
			res, err = json.Marshal(Response[struct {
				Message string `json:"message"`
			}]{
				Result: struct {
					Message string "json:\"message\""
				}{Message: "You must be logged in to access this page!"},
				Code: http.StatusUnauthorized,
			})

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(res)
		}

		context.Set(r, "user", user)

		h.ServeHTTP(w, r)
	})
}

func ProfileMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		database := db.New()

		user := context.Get(r, "user").(*supabase.User)

		var profile IProfile

		err := database.DB.From("profiles").Select("*").Eq("id", user.ID).Execute(&profile)

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			res, err := json.Marshal(Response[struct {
				Message string `json:"message"`
			}]{
				Result: struct {
					Message string "json:\"message\""
				}{Message: "Error fetching profile: " + err.Error()},
				Code: http.StatusNotFound,
			})

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(res)
		}

		context.Set(r, "profile", profile)

		h.ServeHTTP(w, r)
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
