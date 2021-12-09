package oidcconfig

type OpenIDConnectConfig struct {
	AccessMode            string `json:"access_mode"`
	IdpURL                string `json:"idp_url"`
	ClientID              string `json:"client_id"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	Scope                 string `json:"scope"`
	ResponseType          string `json:"response_type"`
	ResponseMode          string `json:"response_mode"`
	SigningKey            string `json:"signing_key"`
}
