package providers

import (
	"log"
	"net/http"
	"os"

	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/markbates/goth/providers/lastfm"
)

type LastFmProvider struct {
	Provider
}

func NewLastFm(w http.ResponseWriter, r *http.Request) *LastFmProvider {
	p := lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), utils.GetCallbackURL("lastfm"))
	database := db.New()
	userCookie, err := r.Cookie("access_token")
	if err != nil {
		return nil
	}

	user, err := database.Auth.User(r.Context(), userCookie.Value)

	if err != nil {
		return nil
	}

	return &LastFmProvider{
		Provider{
			loginHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				session, err := p.BeginAuth("")

				if err != nil {
					return nil, err
				}

				url, err := session.GetAuthURL()

				if err != nil {
					return nil, err
				}

				http.Redirect(w, r, url, http.StatusFound)

				return nil, nil
			},
			callbackHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				token := r.URL.Query().Get("token")

				sess, err := p.GetSession(token)

				if err != nil {
					log.Println(err)
					return nil, err
				}

				var userToken string
				var userName string

				for k, v := range sess {
					if k == "token" {
						userToken = v
					}
					if k == "login" {
						userName = v
					}
				}

				var exists []any

				err = database.DB.From("providers").Select("id").Eq("type", "lastfm").Eq("user", user.ID).Execute(&exists)

				if err != nil {
					log.Println(exists, err)
					return nil, err
				}

				var provider []any

				if len(exists) > 0 {
					database.DB.From("providers").Update(map[string]interface{}{
						"type":                  "lastfm",
						"provider_id":           userName,
						"provider_access_token": userToken,
						"provider_data": map[string]interface{}{
							"status": "active",
						},
					}).Eq("user", user.ID).Execute(&provider)
				} else {
					database.DB.From("providers").Insert(map[string]interface{}{
						"user":                  user.ID,
						"type":                  "lastfm",
						"provider_id":           userName,
						"provider_access_token": userToken,
						"provider_data": map[string]interface{}{
							"status": "active",
						},
					}).Execute(&provider)
				}

				http.Redirect(w, r, os.Getenv("AUTH_WEBSITE"), http.StatusFound)

				return nil, nil
			},
			logoutHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				var res any

				err := database.DB.From("providers").Delete().Eq("type", "lastfm").Eq("user", user.ID).Execute(&res)

				if err != nil {
					log.Println(err)
					return nil, err
				}

				http.Redirect(w, r, os.Getenv("AUTH_WEBSITE"), http.StatusFound)

				return nil, nil
			},
		},
	}
}

func (p *LastFmProvider) LoginHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.loginHandler(w, r)
}

func (p *LastFmProvider) CallbackHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.callbackHandler(w, r)
}

func (p *LastFmProvider) LogoutHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.logoutHandler(w, r)
}
