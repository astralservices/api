package v1

import (
	"encoding/json"
	"net/http"

	auth "github.com/astralservices/api/api/v1/auth"
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

func New(ref *mux.Router) *mux.Router {
	r := ref.PathPrefix("/api/v1").Subrouter()

	r.StrictSlash(true)

	r.HandleFunc("/", IndexHandler)
	r.Handle("/auth", auth.New(r))

	return r
}
