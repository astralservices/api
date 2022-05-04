package v1

import (
	"net/http"

	auth "github.com/astralservices/api/api/v1/auth"
	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func New(ref *mux.Router) *mux.Router {
	r := ref.PathPrefix("/api/v1").Subrouter()

	r.StrictSlash(true)

	r.HandleFunc("/", IndexHandler)
	r.Handle("/auth", auth.New(r))

	return r
}
