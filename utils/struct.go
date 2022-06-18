package utils

import (
	"time"
)

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
	ID           *string     `json:"id,omitempty"`
	CreatedAt    *string     `json:"created_at,omitempty"`
	Owner        *string     `json:"owner,omitempty"`
	Members      *[]string   `json:"members,omitempty"`
	GroupID      *string     `json:"group_id,omitempty"`
	Name         string      `json:"name" form:"name"`
	Logo         string      `json:"logo" form:"logo"`
	Settings     interface{} `json:"settings"`
	Plan         int64       `json:"plan" form:"plan"`
	Visibility   string      `json:"visibility" form:"visibility"`
	Integrations interface{} `json:"integrations"`
	Pending      bool        `json:"pending"`
}

type IWorkspaceMember struct {
	ID        string     `json:"id"`
	CreatedAt string     `json:"created_at"`
	Profile   IProfile   `json:"profile"`
	Workspace IWorkspace `json:"workspace"`
	Role      string     `json:"role"`
	Pending   bool       `json:"pending"`
	InvitedBy string     `json:"invited_by"`
}

type IProvider struct {
	ID                   *string                `json:"id,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	User                 string                 `json:"user"`
	Type                 string                 `json:"type"`
	ProviderID           string                 `json:"provider_id"`
	ProviderAccessToken  string                 `json:"provider_access_token"`
	ProviderRefreshToken string                 `json:"provider_refresh_token"`
	ProviderData         map[string]interface{} `json:"provider_data"`
	ProviderExpiresAt    *time.Time             `json:"provider_expires_at,omitempty"`
	ProviderAvatarUrl    *string                `json:"provider_avatar_url,omitempty"`
	ProviderEmail        *string                `json:"provider_email,omitempty"`
}

type IBlacklist struct {
	ID             int8        `json:"id"`
	CreatedAt      time.Time   `json:"created_at"`
	Moderator      string      `json:"moderator"`
	User           string      `json:"user"`
	DiscordID      string      `json:"discord_id"`
	Reason         string      `json:"reason"`
	Expires        bool        `json:"expires"`
	Expiry         time.Time   `json:"expiry"`
	Flags          interface{} `json:"flags"`
	FactorMatching []string    `json:"factor_matching"`
	Notes          string      `json:"notes"`
}

type IStatistic struct {
	ID        int     `json:"id"`
	Key       string  `json:"key"`
	Value     float32 `json:"value"`
	UpdatedAt string  `json:"updated_at"`
}

type IRegion struct {
	ID         string  `json:"id"`
	Flag       string  `json:"flag"`
	IP         string  `json:"ip"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	PrettyName string  `json:"prettyName"`
	Lat        float64 `json:"lat"`
	Long       float64 `json:"long"`
	MaxBots    int     `json:"maxBots"`
	Status     string  `json:"status"`
}

type ITeamMember struct {
	ID        int             `json:"id"`
	CreatedAt string          `json:"created_at"`
	User      ITeamMemberUser `json:"user"`
	Name      string          `json:"name"`
	Pronouns  string          `json:"pronouns"`
	Location  string          `json:"location"`
	About     string          `json:"about"`
	Role      string          `json:"role"`
}

type ITeamMemberUser struct {
	IdentityData IIdentityData `json:"identity_data"`
}

type IPlan struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PriceMonthly string `json:"priceMonthly"`
	PriceYearly  string `json:"priceYearly"`
	Limit        string `json:"limit"`
	Enabled      bool   `json:"enabled"`
}

type IBot struct {
	ID        string        `json:"id"`
	CreatedAt string        `json:"created_at"`
	Region    string        `json:"region"`
	Owner     string        `json:"owner"`
	Workspace IWorkspace    `json:"workspace"`
	Settings  IBotSettings  `json:"settings"`
	Token     string        `json:"token"`
	Commands  []IBotCommand `json:"commands"`
}

type IBotCommand struct {
	ID      string      `json:"id"`
	Options interface{} `json:"options"`
	Enabled bool        `json:"enabled"`
}

type IBotSettings struct {
	Guild              string      `json:"guild"`
	Prefix             string      `json:"prefix"`
	Status             string      `json:"status"`
	Activities         interface{} `json:"activities"`
	RandomizeActivites bool        `json:"randomizeActivites"`
	ActivityInterval   int         `json:"activityInterval"`
	CurrentActivity    int         `json:"currentActivity"`
	modules            IBotModules `json:"modules"`
}

type IBotModules struct {
	Fun        IBotModule[any] `json:"fun"`
	Moderation IBotModule[struct {
		Logging struct {
			Enabled bool   `json:"enabled"`
			Channel string `json:"channel"`
		} `json:"logging"`
	}] `json:"moderation"`
}

type IBotModule[T any] struct {
	Enabled bool `json:"enabled"`
	Options T    `json:"options"`
}

type IBotAnalytics struct {
	ID        *int        `json:"id,omitempty"`
	Commands  interface{} `json:"commands"`
	Timestamp time.Time   `json:"timestamp"`
	Members   int         `json:"members"`
	Messages  int         `json:"messages"`
	Bot       *IBot       `json:"bot,omitempty"`
}
