package auth

import (
	"encoding/json"
	"net/http"

	"github.com/astralservices/api/api/v1/auth/providers"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/nedpals/supabase-go"
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
		providers.NewRoblox(w, r).CallbackHandler(w, r)

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
		providers.NewRoblox(w, r).LoginHandler(w, r)

	case "lastfm":
		providers.NewLastFm(w, r).LoginHandler(w, r)

	default:
		w.Write([]byte("Unknown provider"))
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	switch provider {
	case "discord":
		providers.NewDiscord().LogoutHandler(w, r)
	case "roblox":
		providers.NewRoblox(w, r).LogoutHandler(w, r)
	case "lastfm":
		providers.NewLastFm(w, r).LogoutHandler(w, r)
	default:
		w.Write([]byte("Unknown provider"))
	}
}

func ProviderHandler(w http.ResponseWriter, r *http.Request) {
	profile := context.Get(r, "profile").(utils.IProfile)

	vars := mux.Vars(r)
	providerId := vars["provider"]

	var providers []utils.IProvider

	database := db.New()

	err := database.DB.From("providers").Select("*").Eq("user", profile.ID).Eq("type", providerId).Execute(&providers)

	if len(providers) == 0 {
		data, err := json.Marshal(utils.Response[any]{
			Result: nil,
			Code:   http.StatusOK,
		})

		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write(data)

		return
	}

	var provider utils.IProvider = providers[0]

	if err != nil {
		data, err := json.Marshal(utils.Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  err.Error(),
		})

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(data)

		return
	}

	data, err := json.Marshal(utils.Response[utils.IProvider]{
		Result: provider,
		Code:   http.StatusOK,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

func ProvidersHandler(w http.ResponseWriter, r *http.Request) {
	profile := context.Get(r, "profile").(utils.IProfile)

	var providers []utils.IProvider

	database := db.New()

	err := database.DB.From("providers").Select("*").Eq("user", profile.ID).Execute(&providers)

	if err != nil {
		data, err := json.Marshal(utils.Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  err.Error(),
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

	if err != nil {
		w.Write([]byte("Error"))
		return
	}

	w.Write(res)
}

type StatusResponse struct {
	Authenticated bool              `json:"authenticated"`
	Blacklist     *utils.IBlacklist `json:"blacklist,omitempty"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*supabase.User)
	var blacklist []utils.IBlacklist

	database := db.New()

	err := database.DB.From("blacklist").Select("*").Eq("user", user.ID).Execute(&blacklist)

	if err != nil {
		data, err := json.Marshal(utils.Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  err.Error(),
		})

		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write(data)

		return
	}

	if len(blacklist) == 0 {
		data, err := json.Marshal(utils.Response[StatusResponse]{
			Result: StatusResponse{
				Authenticated: true,
			},
			Code: http.StatusOK,
		})

		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write(data)

		return
	}

	data, err := json.Marshal(utils.Response[StatusResponse]{
		Result: StatusResponse{
			Authenticated: true,
			Blacklist:     &blacklist[0],
		},
		Code: http.StatusForbidden,
	})

	if err != nil {
		r, e := json.Marshal(utils.Response[any]{
			Result: nil,
			Code:   http.StatusNotFound,
			Error:  err.Error(),
		})

		if e != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write(r)

		return
	}

	w.Write(data)
}
