package auth

import (
	"github.com/gorilla/mux"
)

func New(ref *mux.Router) *mux.Router {
	r := ref.PathPrefix("/auth").Subrouter()

	r.StrictSlash(true)

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/callback", CallbackHandler).Methods("GET")
	r.HandleFunc("/login/{provider}", LoginHandler).Methods("GET")

	return r
}
