package auth

import (
	"net/http"

	db "github.com/astralservices/api/supabase"
	sb "github.com/nedpals/supabase-go"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!\n"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	supabase := db.New()

	authDetails, err := supabase.Auth.SignInWithProvider(sb.ProviderSignInOptions{
		Provider:   "discord",
		RedirectTo: "http://localhost:3000/auth/callback/discord",
		Scopes:     []string{"identify", "email", "guilds"},
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, authDetails.URL, http.StatusPermanentRedirect)
}
