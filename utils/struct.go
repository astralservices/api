package utils

import "time"

type Response[T any] struct {
	Result T      `json:"result"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

type IProfile struct {
	ID               string        `json:"id"`
	Email            string        `json:"email"`
	PreferredName    string        `json:"preferred_name"`
	IdentityData     IIdentityData `json:"identity_data"`
	Access           string        `json:"access"`
	DiscordID        string        `json:"discord_id"`
	RobloxID         interface{}   `json:"roblox_id"`
	StripeCustomerID string        `json:"stripe_customer_id"`
	CreatedAt        string        `json:"created_at"`
	Location         string        `json:"location"`
	Language         string        `json:"language"`
	Pronouns         []string      `json:"pronouns"`
	Hireable         bool          `json:"hireable"`
	About            string        `json:"about"`
	Strengths        []string      `json:"strengths"`
	Weaknesses       []string      `json:"weaknesses"`
	Banner           string        `json:"banner"`
	Verified         bool          `json:"verified"`
	Public           bool          `json:"public"`
	Workspaces       []IWorkspace  `json:"workspaces"`
}

type IIdentityData struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	FullName      string `json:"full_name"`
	AvatarURL     string `json:"avatar_url"`
	ProviderID    string `json:"provider_id"`
	EmailVerified bool   `json:"email_verified"`
}

type IWorkspace struct {
	ID           string      `json:"id"`
	CreatedAt    string      `json:"created_at"`
	Owner        string      `json:"owner"`
	Members      []string    `json:"members"`
	GroupID      string      `json:"group_id"`
	Name         string      `json:"name"`
	Logo         string      `json:"logo"`
	Settings     interface{} `json:"settings"`
	Plan         int64       `json:"plan"`
	Visibility   string      `json:"visibility"`
	Integrations interface{} `json:"integrations"`
	Pending      bool        `json:"pending"`
}

type IProvider struct {
	ID                  string                 `json:"id"`
	CreatedAt           time.Time              `json:"created_at"`
	User                string                 `json:"user"`
	Type                string                 `json:"type"`
	ProviderID          string                 `json:"provider_id"`
	ProviderAccessToken string                 `json:"provider_access_token"`
	ProviderData        map[string]interface{} `json:"provider_data"`
	ProviderExpiresAt   time.Time              `json:"provider_expires_at"`
	DiscordID           string                 `json:"discord_id"`
}
