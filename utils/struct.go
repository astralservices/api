package utils

import (
	"time"
)

// swagger:model APIResponse
type Response[T any] struct {
	// The result of the operation
	// Example: {"message": "API is running!"}
	Result T      `json:"result"`
	// An error, if applicable
	// Example: {"message": "Internal server error"}
	Error  string `json:"error"`

	// The HTTP status code of the response
	// Example: 200
	Code   int    `json:"code"`
}

// swagger:model Profile
type IProfile struct {
	// The UUID of the user / profile
	// Example: f2d179b0-505b-4a4e-b095-64f76619c177
	ID               string        `json:"id"`
	// The email of the user / profile
	// Example: user@example.com
	Email            string        `json:"email"`
	// The preferred name of the profile
	// Example: AmusedGrape
	PreferredName    string        `json:"preferred_name"`
	// The miscellaneous identity data of the profile
	IdentityData     IIdentityData `json:"identity_data"`
	// The access type, if applicable, of the profile
	// Example: beta
	Access           string        `json:"access"`
	// The Discord ID of the profile
	// Example: 401792058970603539
	DiscordID        string        `json:"discord_id"`
	// The Roblox ID of the profile
	// Example: 59692622
	RobloxID         interface{}   `json:"roblox_id"`
	// The Stripe Customer ID of the profile
	// Example: cus_H0I0Z0Z0Z0Z0Z0
	StripeCustomerID string        `json:"stripe_customer_id"`
	// The time the profile was created
	// Example: 2020-01-01T00:00:00Z
	CreatedAt        string        `json:"created_at"`
	// The location of the user / profile
	// Example: Chicago, IL
	Location         string        `json:"location"`
	// The language of the user / profile
	// Example: English
	Language         string        `json:"language"`
	// The prounouns of the user / profile
	// Example: He/Him
	Pronouns         []string      `json:"pronouns"`
	// If the user is hireable
	// Example: true
	Hireable         bool          `json:"hireable"`
	// The about section of the profile
	// Example: I am a developer
	About            string        `json:"about"`
	// The strengths of the profile
	// Example: Great at coding, great at debugging
	Strengths        []string      `json:"strengths"`
	// The weaknesses of the profile
	// Example: I am a bad developer, I am a bad debugger
	Weaknesses       []string      `json:"weaknesses"`
	// The user's banner, taken from Discord
	// Example: https://cdn.discordapp.com/avatars/401792058970603539/...
	Banner           string        `json:"banner"`
	// If the user is verified to be who they are
	// Example: true
	Verified         bool          `json:"verified"`
	// If the user wants others to see their profile
	// Example: true
	Public           bool          `json:"public"`
	// The user's workspaces
	// Example: [...]
	Workspaces       []IWorkspace  `json:"workspaces"`
}

// swagger:model IdentityData
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

// swagger:model Workspace
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

// swagger:model Provider
type IProvider struct {
	ID                  *string                `json:"id,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	User                string                 `json:"user"`
	Type                string                 `json:"type"`
	ProviderID          string                 `json:"provider_id"`
	ProviderAccessToken string                 `json:"provider_access_token"`
	ProviderData        map[string]interface{} `json:"provider_data"`
	ProviderExpiresAt   *time.Time             `json:"provider_expires_at,omitempty"`
	DiscordID           *string                `json:"discord_id,omitempty"`
}

// swagger:model Blacklist
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

// swagger:model Statistic
type IStatistic struct {
	ID        int     `json:"id"`
	Key       string  `json:"key"`
	Value     float32 `json:"value"`
	UpdatedAt string  `json:"updated_at"`
}

// swagger:model Region
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

// swagger:model TeamMember
type ITeamMember struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	User      any    `json:"user"`
	Name      string `json:"name"`
	Pronouns  string `json:"pronouns"`
	Location  string `json:"location"`
	About     string `json:"about"`
	Role      string `json:"role"`
}
