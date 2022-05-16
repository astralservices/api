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
	r.HandleFunc("/logout/{provider}", LogoutHandler).Methods("GET", "OPTIONS", "POST")

	gated := r.PathPrefix("/providers").Subrouter()
	gated.Use(utils.AuthMiddleware)
	gated.Use(utils.ProfileMiddleware)
	gated.HandleFunc("/", ProvidersHandler).Methods("GET", "OPTIONS")
	gated.HandleFunc("/{provider}", ProviderHandler).Methods("GET", "POST", "OPTIONS")

	status := r.PathPrefix("/status").Subrouter()
	status.Use(utils.AuthMiddleware)
	status.HandleFunc("/", StatusHandler)

	return r
}
