package providers

import "net/http"

type LastFmProvider struct {
	Provider
}

func NewLastFm() *LastFmProvider {
	return &LastFmProvider{
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

func (p *LastFmProvider) LoginHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.loginHandler(w, r)
}

func (p *LastFmProvider) CallbackHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return p.callbackHandler(w, r)
}
