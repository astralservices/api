package providers

import "net/http"

type RobloxProvider struct {
	Provider
}

func NewRoblox() *RobloxProvider {
	return &RobloxProvider{
		Provider{
			loginHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
				return nil, nil
			},
			callbackHandler: func(w http.ResponseWriter, r *http.Request) ([]byte, error) {
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
