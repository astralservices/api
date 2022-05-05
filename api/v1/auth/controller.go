package auth

import (
	"github.com/astralservices/api/utils"
	"github.com/gorilla/mux"
)

func New(ref *mux.Router) *mux.Router {
	r := ref.PathPrefix("/auth").Subrouter()

	r.StrictSlash(true)

	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/callback/{provider}", CallbackHandler).Methods("GET", "OPTIONS", "POST")
	r.HandleFunc("/login/{provider}", LoginHandler).Methods("GET", "OPTIONS", "POST")

	gated := r.PathPrefix("/providers").Subrouter()
	gated.Use(utils.AuthMiddleware)
	gated.HandleFunc("/{provider}", ProviderHandler).Methods("GET", "POST", "OPTIONS")

	return r
}
