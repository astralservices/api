package providers

import "net/http"

type Provider struct {
	loginHandler    func(w http.ResponseWriter, r *http.Request) ([]byte, error)
	callbackHandler func(w http.ResponseWriter, r *http.Request) ([]byte, error)
	logoutHandler   func(w http.ResponseWriter, r *http.Request) ([]byte, error)
}

type Providers struct {
	roblox  *RobloxProvider
	lastfm  *LastFmProvider
	discord *DiscordProvider
}

func New() *Providers {
	p := &Providers{}

	return p
}
