package providers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
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
				
				redirect := utils.GetCallbackURL("discord")

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
				var accessToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxNjUyODMwMTE4LCJzdWIiOiI3NTQ4NmU4NS1mZmFlLTQwNzAtODhhYi1kNjFiYzMyNWUyMmUiLCJlbWFpbCI6Im1lQGFtdXNlZGdyYXBlLnh5eiIsInBob25lIjoiIiwiYXBwX21ldGFkYXRhIjp7InByb3ZpZGVyIjoiZGlzY29yZCIsInByb3ZpZGVycyI6WyJkaXNjb3JkIl19LCJ1c2VyX21ldGFkYXRhIjp7ImF2YXRhcl91cmwiOiJodHRwczovL2Nkbi5kaXNjb3JkYXBwLmNvbS9hdmF0YXJzLzQwMTc5MjA1ODk3MDYwMzUzOS9hX2E0NzY4YTBlMmZmNDQ2YzgzNzZiZTQxYjBiMTQyYWM3LmdpZiIsImVtYWlsIjoibWVAYW11c2VkZ3JhcGUueHl6IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZ1bGxfbmFtZSI6IkFtdXNlZEdyYXBlIiwiaXNzIjoiaHR0cHM6Ly9kaXNjb3JkLmNvbS9hcGkiLCJuYW1lIjoiQW11c2VkR3JhcGUjMDAwMSIsInBpY3R1cmUiOiJodHRwczovL2Nkbi5kaXNjb3JkYXBwLmNvbS9hdmF0YXJzLzQwMTc5MjA1ODk3MDYwMzUzOS9hX2E0NzY4YTBlMmZmNDQ2YzgzNzZiZTQxYjBiMTQyYWM3LmdpZiIsInByb3ZpZGVyX2lkIjoiNDAxNzkyMDU4OTcwNjAzNTM5Iiwic3ViIjoiNDAxNzkyMDU4OTcwNjAzNTM5In0sInJvbGUiOiJhdXRoZW50aWNhdGVkIn0.PtSsScW5CjvWt-OqV67KhehJlh756zDl78Xij8sbAgc"
				var providerToken string = "4qZ4uZWGRWvGQJ4DvRkLP2iXUgQhS9"

				// v, err := url.ParseQuery(r.URL.Fragment)

				// if err != nil {
				// 	log.Println(err)
				// 	w.Write([]byte("Error parsing query"))
				// 	return nil, err
				// }

				// accessToken = v.Get("access_token")
				// providerToken = v.Get("provider_token")

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
			logoutHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				http.SetCookie(w, &http.Cookie{
					Name:     "access_token",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})

				http.SetCookie(w, &http.Cookie{
					Name:     "provider_token",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					HttpOnly: true,
				})

				http.Redirect(w, r, os.Getenv("AUTH_WEBSITE"), http.StatusFound)

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

func (p *DiscordProvider) LogoutHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.logoutHandler(w, r)
}
