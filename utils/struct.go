package utils

import (
	"time"

	"github.com/aybabtme/orderedjson"
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

type IWorkspaceMemberWithoutProfile struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Workspace any    `json:"workspace"`
	Role      string `json:"role"`
	Pending   bool   `json:"pending"`
	InvitedBy string `json:"invited_by"`
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

	Bots int `json:"bots"`
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
	ID        *string       `json:"id,omitempty"`
	CreatedAt *time.Time    `json:"created_at,omitempty"`
	Region    string        `json:"region" form:"region"`
	Owner     *string       `json:"owner,omitempty"`
	Workspace *IWorkspace   `json:"workspace,omitempty"`
	Settings  IBotSettings  `json:"settings" form:"settings"`
	Token     string        `json:"token" form:"token"`
	Commands  []IBotCommand `json:"commands" form:"commands"`
}

type IBotCommand struct {
	ID      string      `json:"id"`
	Options interface{} `json:"options"`
	Enabled bool        `json:"enabled"`
}

type IBotSettings struct {
	Guild               string         `json:"guild" form:"guild"`
	Prefix              string         `json:"prefix" form:"prefix"`
	Status              string         `json:"status" form:"status"`
	Activities          []IBotActivity `json:"activities" form:"activities"`
	RandomizeActivities bool           `json:"randomizeActivities" form:"randomizeActivities"`
	ActivityInterval    int            `json:"activityInterval" form:"activityInterval"`
	CurrentActivity     int            `json:"currentActivity"`
	Modules             IBotModules    `json:"modules" form:"modules"`
}

type IBotActivity struct {
	Name string `json:"name" form:"name"`
	Type string `json:"type" form:"type"`
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

type IWorkspaceIntegration struct {
	ID          int         `json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	Integration string      `json:"integration"`
	Settings    interface{} `json:"settings"`
	Workspace   string      `json:"workspace"`
	Enabled     bool        `json:"enabled"`
}

type IIntegration struct {
	ID               string             `json:"id"`
	CreatedAt        time.Time          `json:"created_at"`
	Name             string             `json:"name"`
	PrettyName       string             `json:"prettyName"`
	Icon             string             `json:"icon"`
	IsIconSimpleIcon bool               `json:"isIconSimpleIcon"`
	Website          string             `json:"website"`
	Enabled          bool               `json:"enabled"`
	Description      string             `json:"description"`
	Schema           IIntegrationSchema `json:"schema"`
}

type IIntegrationSchema struct {
	Type       string     `json:"type"`
	Title      string     `json:"title"`
	Properties OrderedMap `json:"properties"`
}

type IIntegrationSchemaInput struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

type IRobloxIntegration struct {
	GroupId           string                              `json:"groupId" form:"groupId"`
	Token             string                              `json:"token" form:"token"`
	MemberCounter     IRobloxIntegrationMemberCounter     `json:"memberCounter" form:"memberCounter"`
	ShoutProxy        IRobloxIntegrationShoutProxy        `json:"shoutProxy" form:"shoutProxy"`
	BadActorDetection IRobloxIntegrationBadActorDetection `json:"badActorDetection" form:"badActorDetection"`
	WallFilter        IRobloxIntegrationWallFilter        `json:"wallFilter" form:"wallFilter"`
}

type IRobloxIntegrationMemberCounter struct {
	Enabled bool   `json:"enabled" form:"enabled"`
	Message string `json:"message" form:"message"`
	Webhook string `json:"webhook" form:"webhook"`
	GroupId string `json:"groupId" form:"groupId"`
}

type IRobloxIntegrationShoutProxy struct {
	Enabled bool   `json:"enabled" form:"enabled"`
	Webhook string `json:"webhook" form:"webhook"`
	GroupId string `json:"groupId" form:"groupId"`
}

type IRobloxIntegrationBadActorDetection struct {
	Enabled bool `json:"enabled" form:"enabled"`
	Factors struct {
		BannedGroups    string `json:"bannedGroups" form:"bannedGroups"`
		SketchyUsername bool   `json:"sketchyUsername" form:"sketchyUsername"`
		NoDescription   bool   `json:"noDescription" form:"noDescription"`
	}
}

type IRobloxIntegrationWallFilter struct {
	Enabled       bool   `json:"enabled" form:"enabled"`
	BannedPhrases string `json:"bannedPhrases" form:"bannedPhrases"`
}

type IRobloxSchema struct {
	GroupId           orderedjson.Map `json:"groupId"`
	Token             orderedjson.Map `json:"token"`
	MemberCounter     orderedjson.Map `json:"member_counter"`
	ShoutProxy        orderedjson.Map `json:"shout_proxy"`
	BadActorDetection orderedjson.Map `json:"bad_actor_detection"`
	WallFilter        orderedjson.Map `json:"wall_filter"`
	Submit            orderedjson.Map `json:"submit"`
}

type IRobloxSchemaMemberCounter struct {
	Enabled orderedjson.Map `json:"enabled"`
	Message orderedjson.Map `json:"message"`
	Webhook orderedjson.Map `json:"webhook"`
	GroupId orderedjson.Map `json:"groupId"`
}

type IRobloxSchemaShoutProxy struct {
	Enabled orderedjson.Map `json:"enabled"`
	Webhook orderedjson.Map `json:"webhook"`
	GroupId orderedjson.Map `json:"groupId"`
}

type IRobloxSchemaBadActorDetection struct {
	Enabled orderedjson.Map `json:"enabled"`
	Factors orderedjson.Map `json:"factors"`
}

type IRobloxSchemaFactors struct {
	BannedGroups    orderedjson.Map `json:"banned_groups"`
	SketchyUsername orderedjson.Map `json:"sketchyUsername"`
	NoDescription   orderedjson.Map `json:"noDescription"`
}

type IRobloxSchemaWallFilter struct {
	BannedPhrases orderedjson.Map `json:"bannedPhrases"`
}

type IBotModerationAction struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Bot       string    `json:"bot"`
	Guild     string    `json:"guild"`
	Action    string    `json:"action"`
	Moderator string    `json:"moderator"`
	Reason    string    `json:"reason"`
	Expires   bool      `json:"expires"`
	Expiry    time.Time `json:"expiry"`
	User      string    `json:"user"`
}

type IDiscordApiUser struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	Avatar        *string `json:"avatar,omitempty"`
	Bot           *bool   `json:"bot,omitempty"`
	System        *bool   `json:"system,omitempty"`
	MFAEnabled    *bool   `json:"mfa_enabled,omitempty"`
	Banner        *string `json:"banner,omitempty"`
	AccentColor   *int    `json:"accent_color,omitempty"`
	Locale        *string `json:"locale,omitempty"`
	Verified      *bool   `json:"verified,omitempty"`
	Email         *string `json:"email,omitempty"`
	Flags         *int    `json:"flags,omitempty"`
	PremiumType   *int    `json:"premium_type,omitempty"`
	PublicFlags   *int    `json:"public_flags,omitempty"`
}

type IIntegrationData struct {
	ID                   int     `json:"id"`
	CreatedAt            *string `json:"created_at,omitempty"`
	Integration          string  `json:"integration"`
	WorkspaceIntegration int     `json:"workspaceIntegration"`
	User                 string  `json:"user"`
	Data                 any     `json:"data"`
}
