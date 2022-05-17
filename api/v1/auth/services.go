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

// Provider Callback
// @Summary Callback for provider
// @Description 
// @ID provider-callback
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param provider path string true "Provider"
// @Success 301
// @Failure 500 {object} utils.DocsAPIError "Internal Server Error"
// @Router /auth/callback/{provider} [get]
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

// Provider Login
// @Summary Login to provider
// @Description 
// @ID provider-login
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param provider path string true "Provider"
// @Success 301
// @Failure 500 {object} utils.DocsAPIError "Internal Server Error"
// @Router /auth/login/{provider} [post]
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

// Provider Logout
// @Summary Logout of provider
// @Description 
// @ID provider-logout
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param provider path string true "Provider"
// @Success 301
// @Failure 500 {object} utils.DocsAPIError "Internal Server Error"
// @Router /auth/logout/{provider} [post]
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

// Provider Information
// @Summary Get provider information
// @Description 
// @ID provider-info
// @Tags Authentication
// @Accept  json
// @Produce  json
// @securityDefinitions.apikey ApiKeyAuth
// @Param provider path string true "Provider"
// @Success 200 {object} utils.IProvider "OK"
// @Failure 500 {object} utils.DocsAPIError "Internal Server Error"
// @Router /auth/providers/{provider} [get]
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

// Providers Information
// @Summary Get all provider information
// @Description 
// @ID providers-info
// @Tags Authentication
// @Accept  json
// @Produce  json
// @securityDefinitions.apikey ApiKeyAuth
// @Success 200 {array} utils.IProvider "OK"
// @Failure 500 {object} utils.DocsAPIError "Internal Server Error"
// @Router /auth/providers [get]
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

// User Status
// @Summary Get the authenticated user's status
// @Description 
// @ID user-status
// @Tags User
// @Accept  json
// @Produce  json
// @securityDefinitions.apikey ApiKeyAuth
// @Success 200 {object} StatusResponse "OK"
// @Failure 500 {object} utils.DocsAPIError "Internal Server Error"
// @Router /auth/status [get]
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
