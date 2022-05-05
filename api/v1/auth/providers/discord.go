package providers

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	db "github.com/astralservices/api/supabase"
	sb "github.com/nedpals/supabase-go"
)

type DiscordProvider struct {
	Provider
}

func NewDiscord() *DiscordProvider {
	return &DiscordProvider{
		Provider{
			loginHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				supabase := db.New()
				var redirectTo string
				var redirect string
				if r.TLS != nil {
					redirect = "https://" + r.Host + "/api/v1/auth/callback/discord"
				} else {
					redirect = "http://" + r.Host + "/api/v1/auth/callback/discord"
				}

				// get the redirect url from the form data
				redirectTo = r.FormValue("redirect")

				authDetails, err := supabase.Auth.SignInWithProvider(sb.ProviderSignInOptions{
					Provider:   "discord",
					RedirectTo: redirect,
					Scopes:     []string{"identify", "email", "guilds"},
				})

				if err != nil {
					w.Write([]byte(err.Error()))
					return nil, err
				}

				r.AddCookie(&http.Cookie{
					Name:     "redirectTo",
					Value:    redirectTo,
					Path:     "/",
					Domain:   r.Host,
					Expires:  time.Now().Add(time.Hour * 24 * 7),
					HttpOnly: true,
					Secure:   os.Getenv("ENV") == "production",
				})

				http.Redirect(w, r, authDetails.URL, http.StatusFound)

				return nil, nil
			},
			callbackHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				var accessToken string
				var providerToken string

				v, err := url.ParseQuery(r.URL.Fragment)

				if err != nil {
					log.Println(err)
					w.Write([]byte("Error parsing query"))
					return nil, err
				}

				accessToken = v.Get("access_token")
				providerToken = v.Get("provider_token")

				var domain string

				domain = r.Host
				if pos := strings.Index(domain, ":"); pos != -1 {
					domain = domain[:pos]
				}

				log.Println(accessToken, providerToken, domain)

				// set the access token and provider token in the session cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "access_token",
					Value:    accessToken,
					Path:     "/",
					Domain:   domain,
					Expires:  time.Now().Add(time.Hour * 24 * 7),
					HttpOnly: true,
					Secure:   os.Getenv("ENV") == "production",
				})
				http.SetCookie(w, &http.Cookie{
					Name:     "provider_token",
					Value:    providerToken,
					Path:     "/",
					Domain:   domain,
					Expires:  time.Now().Add(time.Hour * 24 * 7),
					HttpOnly: true,
					Secure:   os.Getenv("ENV") == "production",
				})

				// get the redirect url from the redirectTo cookie
				redirectTo := os.Getenv("AUTH_WEBSITE")

				for _, cookie := range r.Cookies() {
					if cookie.Name == "redirectTo" {
						redirectTo = cookie.Value
					}
				}

				http.Redirect(w, r, redirectTo, http.StatusFound)

				return nil, nil
			},
		},
	}
}

func (p *DiscordProvider) LoginHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.loginHandler(w, r)
}

func (p *DiscordProvider) CallbackHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.callbackHandler(w, r)
}
