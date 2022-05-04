package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/astralservices/api/supabase"
	"github.com/gorilla/mux"
	sb "github.com/nedpals/supabase-go"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	switch provider {
	case "discord":
		DiscordCallbackHandler(w, r)

	default:
		w.Write([]byte("Unknown provider"))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	supabase := db.New()
	var redirect string
	if r.TLS != nil {
		redirect = "https://" + r.Host + "/api/v1/auth/callback/discord"
	} else {
		redirect = "http://" + r.Host + "/api/v1/auth/callback/discord"
	}
	w.Header().Add("redirecturl", redirect)
	authDetails, err := supabase.Auth.SignInWithProvider(sb.ProviderSignInOptions{
		Provider:   "discord",
		RedirectTo: redirect,
		Scopes:     []string{"identify", "email", "guilds"},
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, authDetails.URL, http.StatusPermanentRedirect)
}

func DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var accessToken string
	var providerToken string

	accessToken = r.URL.Query().Get("access_token")
	providerToken = r.URL.Query().Get("provider_token")

	var domain string

	domain = r.Host
	if pos := strings.Index(domain, ":"); pos != -1 {
		domain = domain[:pos]
	}

	// set the access token and provider token in the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
		Path:  "/",
		// Domain:   domain,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "provider_token",
		Value: providerToken,
		Path:  "/",
		// Domain:   domain,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HttpOnly: true,
	})

	http.Redirect(w, r, os.Getenv("AUTH_WEBSITE"), http.StatusPermanentRedirect)
}
