package v1

import (
	"encoding/json"
	"net/http"
	"sort"

	auth "github.com/astralservices/api/api/v1/auth"
	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(utils.Response[struct {
		Message string `json:"message"`
	}]{
		Result: struct {
			Message string "json:\"message\""
		}{Message: "API v1 is running!"},
		Code: http.StatusOK,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

// swagger:route GET /api/v1/stats Stats get-stats
//
// Gets the various statistics of Astral Services.
//
// responses:
//   200: APIResponse
//   404: ErrorResponse
//   500: ErrorResponse
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats []utils.IStatistic

	database := db.New()

	err := database.DB.From("stats").Select("*").Execute(&stats)

	if err != nil {
		errorData, err := json.Marshal(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(errorData)

		return
	}

	data, err := json.Marshal(utils.Response[[]utils.IStatistic]{
		Result: stats,
		Code:   http.StatusOK,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

func RegionsHandler(w http.ResponseWriter, r *http.Request) {
	var regions []utils.IRegion

	database := db.New()

	err := database.DB.From("regions").Select("id, flag, city, region, country, prettyName, lat, long, maxBots").Execute(&regions)

	if err != nil {
		errorData, err := json.Marshal(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(errorData)

		return
	}

	// remove the region "localhost"
	for i, region := range regions {
		if region.ID == "localhost" {
			regions = append(regions[:i], regions[i+1:]...)
			break
		}
	}

	data, err := json.Marshal(utils.Response[[]utils.IRegion]{
		Result: regions,
		Code:   http.StatusOK,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

func TeamHandler(w http.ResponseWriter, r *http.Request) {
	var teams []utils.ITeamMember

	database := db.New()

	err := database.DB.From("teamMembers").Select("*, user(identity_data)").Execute(&teams)

	if err != nil {
		errorData, err := json.Marshal(utils.Response[any]{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(errorData)

		return
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].ID < teams[j].ID
	})

	data, err := json.Marshal(utils.Response[[]utils.ITeamMember]{
		Result: teams,
		Code:   http.StatusOK,
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
}

func New(ref *mux.Router) *mux.Router {
	r := ref.PathPrefix("/api/v1").Subrouter()

	r.StrictSlash(true)

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/stats", StatsHandler)
	r.HandleFunc("/regions", RegionsHandler)
	r.HandleFunc("/team", TeamHandler)

	r.Handle("/auth", auth.New(r))

	return r
}
