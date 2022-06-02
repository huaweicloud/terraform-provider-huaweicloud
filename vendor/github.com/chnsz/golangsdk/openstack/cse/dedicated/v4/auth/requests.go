package auth

import (
	"github.com/chnsz/golangsdk"
)

// CreateOpts is the structure required by the Create method to create a token for connecting to the engine.
type CreateOpts struct {
	// Account name.
	Name string `json:"name" required:"true"`
	// Account password.
	Password string `json:"password" required:"true"`
}

// Create is a method to create a token using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(rootURL(c), b, &r, nil)
	return &r, err
}

// BuildMoreHeaderUsingToken is a method to build a specified request header using given token.
func BuildMoreHeaderUsingToken(c *golangsdk.ServiceClient, token string) map[string]string {
	moreHeader := map[string]string{
		"Content-Type": "application/json",
		"X-Language":   "en-us",
	}

	if token != "" {
		moreHeader["Authorization"] = token
	}
	return moreHeader
}
