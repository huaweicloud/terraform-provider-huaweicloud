package auth

// CreateResp is the structure that represents the response of the API request.
type CreateResp struct {
	// The obtained user token is valid for 12 hours.
	Token string `json:"token"`
}
