package responses

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

// GetResult represents a result of the Get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents a result of the Update operation.
type UpdateResult struct {
	commonResult
}

type Response struct {
	// Response ID.
	Id string `json:"id"`
	// Response name.
	Name string `json:"name"`
	// Response type definition, which includes a key and value. Options of the key:
	//     AUTH_FAILURE: Authentication failed.
	//     AUTH_HEADER_MISSING: The identity source is missing.
	//     AUTHORIZER_FAILURE: Custom authentication failed.
	//     AUTHORIZER_CONF_FAILURE: There has been a custom authorizer error.
	//     AUTHORIZER_IDENTITIES_FAILURE: The identity source of the custom authorizer is invalid.
	//     BACKEND_UNAVAILABLE: The backend service is unavailable.
	//     BACKEND_TIMEOUT: Communication with the backend service timed out.
	//     THROTTLED: The request was rejected due to request throttling.
	//     UNAUTHORIZED: The app you are using has not been authorized to call the API.
	//     ACCESS_DENIED: Access denied.
	//     NOT_FOUND: No API is found.
	//     REQUEST_PARAMETERS_FAILURE: The request parameters are incorrect.
	//     DEFAULT_4XX: Another 4XX error occurred.
	//     DEFAULT_5XX: Another 5XX error occurred.
	// Each error type is in JSON format.
	Responses map[string]ResponseInfo `json:"responses"`
	// Indicates whether the group response is the default response.
	IsDefault bool `json:"default"`
	// Creation time.
	CreateTime string `json:"create_time"`
	// Update time.
	UpdateTime string `json:"update_time"`
}

// Extract is a method to extract an response struct.
func (r commonResult) Extract() (*Response, error) {
	var s Response
	err := r.ExtractInto(&s)
	return &s, err
}

// ResponsePage represents the response pages of the List operation.
type ResponsePage struct {
	pagination.SinglePageBase
}

// ExtractResponses is a method to extract an response struct list.
func ExtractResponses(r pagination.Page) ([]Response, error) {
	var s []Response
	err := r.(ResponsePage).Result.ExtractIntoSlicePtr(&s, "responses")
	return s, err
}

// DeleteResult represents a result of the Delete and DeleteSpecResp method.
type DeleteResult struct {
	golangsdk.ErrResult
}

type SpecRespResult struct {
	commonResult
}

// GetSpecRespResult represents a result of the GetSpecResp method.
type GetSpecRespResult struct {
	SpecRespResult
}

// UpdateSpecRespResult represents a result of the UpdateSpecResp method.
type UpdateSpecRespResult struct {
	SpecRespResult
}

// ExtractSpecResp is a method to extract an response struct using a specifies key.
func (r SpecRespResult) ExtractSpecResp(key string) (*ResponseInfo, error) {
	var s ResponseInfo
	err := r.ExtractIntoStructPtr(&s, key)
	return &s, err
}
