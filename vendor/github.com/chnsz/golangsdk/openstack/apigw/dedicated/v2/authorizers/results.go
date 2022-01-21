package authorizers

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// UpdateResult represents a result of the Update method.
type UpdateResult struct {
	commonResult
}

// GetResult represents a result of the Update method.
type GetResult struct {
	commonResult
}

// CustomAuthorizer is a struct that represents the result of Create, Update and Get methods.
type CustomAuthorizer struct {
	// Custom authorizer name., which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	Name string `json:"name"`
	// Custom authorizer type, which support 'FRONTEND' and 'BACKEND'.
	Type string `json:"type"`
	// Authorizer type, and the value is 'FUNC'.
	AuthorizerType string `json:"authorizer_type"`
	// Function URN.
	AuthorizerURI string `json:"authorizer_uri"`
	// Identity source.
	Identities []Identity `json:"identities"`
	// Maximum cache age.
	TTL int `json:"ttl"`
	// User data.
	UserData string `json:"user_data"`
	// Indicates whether to send the body.
	IsBodySend bool `json:"need_body"`
	// Custom authorizer ID
	ID string `json:"id"`
	// Creation time.
	CreateTime string `json:"create_time"`
}

// Identity is an object struct that represents the elements of the Identities parameter.
type Identity struct {
	// Parameter name.
	Name string `json:"name"`
	// Parameter location, which support 'HEADER' and 'QUERY'.
	Location string `json:"location"`
	// Parameter verification expression.
	// The default value is null, indicating that no verification is performed.
	Validation string `json:"validation"`
}

func (r commonResult) Extract() (*CustomAuthorizer, error) {
	var s CustomAuthorizer
	err := r.ExtractInto(&s)
	return &s, err
}

// CustomAuthPage represents the response pages of the List method.
type CustomAuthPage struct {
	pagination.SinglePageBase
}

// ExtractCustomAuthorizers is a method which to extract the response to a custom authorizer list.
func ExtractCustomAuthorizers(r pagination.Page) ([]CustomAuthorizer, error) {
	var s []CustomAuthorizer
	err := r.(CustomAuthPage).Result.ExtractIntoSlicePtr(&s, "authorizer_list")
	return s, err
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	golangsdk.ErrResult
}
