package providers

type ProviderLinks struct {
	Self      string `json:"self"`
	Protocols string `json:"protocols"`
	Next      string `json:"next"`
	Previous  string `json:"previous"`
}

type Provider struct {
	SsoType     string        `json:"sso_type"`
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Enabled     bool          `json:"enabled"`
	RemoteIDs   []string      `json:"remote_ids"`
	Links       ProviderLinks `json:"links"`
}

type ProviderList struct {
	Links             ProviderLinks `json:"links"`
	IdentityProviders []Provider    `json:"identity_providers"`
}
