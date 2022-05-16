package providers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	db "github.com/astralservices/api/supabase"
	"github.com/astralservices/api/utils"
	"github.com/astralservices/goblox/goblox"
)

type RobloxProvider struct {
	Provider
}

func NewRoblox(w http.ResponseWriter, r *http.Request) *RobloxProvider {
	database := db.New()
	userCookie, err := r.Cookie("access_token")
	if err != nil {
		d, e := json.Marshal(utils.Response[struct {
			Message string `json:"message"`
		}]{
			Code: http.StatusUnauthorized,
			Result: struct {
				Message string `json:"message"`
			}{
				Message: "You must be logged in to access this resource",
			},
			Error: err.Error(),
		})

		if e != nil {
			w.Write([]byte(e.Error()))
			return nil
		}

		w.Write(d)

		return nil
	}

	user, err := database.Auth.User(r.Context(), userCookie.Value)

	if err != nil {
		d, e := json.Marshal(utils.Response[struct {
			Message string `json:"message"`
		}]{
			Code: http.StatusUnauthorized,
			Result: struct {
				Message string `json:"message"`
			}{
				Message: "You must be logged in to access this resource",
			},
			Error: err.Error(),
		})

		if e != nil {
			w.Write([]byte(e.Error()))
			return nil
		}

		w.Write(d)

		return nil
	}

	roblox := goblox.New()

	return &RobloxProvider{
		Provider{
			loginHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				r.ParseForm()

				username := r.Form.Get("username")

				var callbackUrl string
				site := os.Getenv("AUTH_WEBSITE")
				callbackUrl = site + "/providers/roblox"

				if username == "" {
					http.Redirect(w, r, callbackUrl+"/error?error=No+Username+Provided", http.StatusFound)

					return nil, nil
				}

				robloxUser, err := roblox.Users.GetUserByUsername(username)

				if err != nil {
					http.Redirect(w, r, callbackUrl+"/error?error=User+Not+Found", http.StatusFound)
				}

				type ProviderData struct {
					Username string `json:"username"`
					AuthCode string `json:"auth_code"`
					Status   string `json:"status"`
				}

				type UserProvider struct {
					ID           string       `json:"id"`
					User         string       `json:"user"`
					ProviderData ProviderData `json:"provider_data"`
				}

				var userProviders []UserProvider = []UserProvider{}

				err = database.DB.From("providers").Select("id, user, provider_data").Eq("user", user.ID).Eq("type", "roblox").Execute(&userProviders)

				if err != nil {
					log.Println(err.Error())
					http.Redirect(w, r, callbackUrl+"/error?error=Error+checking+for+existing+account", http.StatusFound)
					return nil, nil
				}

				if len(userProviders) > 0 {
					var userProvider UserProvider = userProviders[0]

					if userProvider.ProviderData.Status == "active" {
						http.Redirect(w, r, callbackUrl+"/error?error=Account+Already+Linked", http.StatusFound)
						return nil, nil
					}

					if userProvider.ProviderData.Status == "pending" {
						http.Redirect(w, r, callbackUrl+"/code?code="+userProvider.ProviderData.AuthCode, http.StatusFound)
						return nil, nil
					}
				}

				var provider any
				var words []string

				for i := 0; i < 8; i++ {
					words = append(words, utils.RandomWord())
				}

				authCode := strings.Join(words, " ")
				log.Println(authCode)
				err = database.DB.From("providers").Insert(utils.IProvider{
					User:       user.ID,
					Type:       "roblox",
					ProviderID: strconv.Itoa(int(robloxUser.ID)),
					ProviderData: map[string]interface{}{
						"status":    "pending",
						"auth_code": authCode,
						"username":  robloxUser.Name,
					},
				}).Execute(&provider)

				if err != nil {
					log.Println(err.Error())
					http.Redirect(w, r, callbackUrl+"/error?error=Error+inserting+row", http.StatusFound)
					return nil, nil
				}

				http.Redirect(w, r, callbackUrl+"/code?code="+authCode, http.StatusFound)

				return nil, nil
			},
			callbackHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				r.ParseForm()

				authCode := r.Form.Get("code")

				var callbackUrl string
				site := os.Getenv("AUTH_WEBSITE")
				callbackUrl = site + "/providers/roblox"

				if authCode == "" {
					http.Redirect(w, r, callbackUrl+"/error?error=No+Auth+Code+Provided", http.StatusFound)

					return nil, nil
				}

				var providers []utils.IProvider

				err := database.DB.From("providers").Select("id, type, provider_data, provider_id, user").Eq("type", "roblox").Eq("user", user.ID).Execute(&providers)

				log.Println(providers, authCode)

				var provider utils.IProvider = providers[0]

				if err != nil {
					log.Println(err.Error())
					http.Redirect(w, r, callbackUrl+"/error?error=Error+fetching+auth+code", http.StatusFound)

					return nil, nil
				}

				if provider.ProviderData["auth_code"] != authCode {
					http.Redirect(w, r, callbackUrl+"/error?error=Invalid+Auth+Code", http.StatusFound)

					return nil, nil
				}

				if provider.ProviderData["status"] != "pending" {
					http.Redirect(w, r, callbackUrl+"/error?error=Account+Already+Linked", http.StatusFound)

					return nil, nil
				}

				userId, err := strconv.ParseInt(provider.ProviderID, 10, 64)

				if err != nil {
					log.Println(err.Error())

					http.Redirect(w, r, callbackUrl+"/error?error=Error+parsing+user+id", http.StatusFound)

					return nil, nil
				}

				robloxUser, err := roblox.Users.GetUserById(userId)

				if err != nil {
					log.Println(err.Error())

					http.Redirect(w, r, callbackUrl+"/error?error=Error+fetching+Roblox+user", http.StatusFound)

					return nil, nil
				}

				// check if the robloxUser's description contains the code

				if !strings.Contains(robloxUser.Description, authCode) {
					http.Redirect(w, r, callbackUrl+"/error?error=Auth+Code+Not+Found+in+profile", http.StatusFound)

					return nil, nil
				}

				// update the provider's status to active

				var updatedProviders []utils.IProvider

				provider.ProviderData["status"] = "active"

				err = database.DB.From("providers").Update(map[string]interface{}{
					"provider_data": provider.ProviderData,
				}).Eq("id", *provider.ID).Execute(&updatedProviders)

				if err != nil {
					log.Println(err.Error())

					http.Redirect(w, r, callbackUrl+"/error?error=Error+updating+user+in+database", http.StatusFound)

					return nil, nil
				}

				http.Redirect(w, r, os.Getenv("AUTH_WEBSITE"), http.StatusFound)

				return nil, nil
			},
			logoutHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				var res any
				err := database.DB.From("providers").Delete().Eq("type", "roblox").Eq("user", user.ID).Execute(&res)

				if err != nil {
					return nil, err
				}

				http.Redirect(w, r, os.Getenv("AUTH_WEBSITE"), http.StatusFound)

				return nil, nil
			},
		},
	}
}

func (p *RobloxProvider) LoginHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.loginHandler(w, r)
}

func (p *RobloxProvider) CallbackHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.callbackHandler(w, r)
}

func (p *RobloxProvider) LogoutHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.logoutHandler(w, r)
}
