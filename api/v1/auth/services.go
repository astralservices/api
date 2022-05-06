package auth

import (
	"encoding/json"
	"net/http"

	"github.com/astralservices/api/api/v1/auth/providers"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	switch provider {
	case "discord":
		providers.NewDiscord().CallbackHandler(w, r)

	case "roblox":
		providers.NewRoblox().CallbackHandler(w, r)

	case "lastfm":
		providers.NewLastFm(w, r).CallbackHandler(w, r)

	default:
		w.Write([]byte("Unknown provider"))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	switch provider {
	case "discord":
		providers.NewDiscord().LoginHandler(w, r)

	case "roblox":
		providers.NewRoblox().LoginHandler(w, r)

	case "lastfm":
		providers.NewLastFm(w, r).LoginHandler(w, r)

	default:
		w.Write([]byte("Unknown provider"))
	}
}

func ProviderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	switch provider {
	default:
		data, err := json.Marshal(map[string]string{
			"error": "Unknown provider",
		})

		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write(data)
	}
}

func ProvidersHandler(w http.ResponseWriter, r *http.Request) {
	profile := context.Get(r, "profile").(utils.IProfile)

	var providers []utils.IProvider

	database := db.New()

	err := database.DB.From("providers").Select("*").Eq("id", profile.ID).Execute(&providers)

	if err != nil {
		data, err := json.Marshal(map[string]string{
			"error": "Profile providers not found!",
		})

		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write(data)
		return
	}

	res, err := json.Marshal(utils.Response[[]utils.IProvider]{
		Result: providers,
		Code:   http.StatusOK,
	})

	w.Write(res)

	return
}
