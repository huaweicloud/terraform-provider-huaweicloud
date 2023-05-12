package acls

import "github.com/chnsz/golangsdk/pagination"

// Policy is the structure represents the ACL policy details.
type Policy struct {
	// The ACL name.
	Name string `json:"acl_name"`
	// The ACL type. The valid values are as follows:
	// + PERMIT
	// + DENY
	Type string `json:"acl_type"`
	// The value of the ACL policy.
	Value string `json:"acl_value"`
	// The entity type. The valid values are as follows:
	// + IP
	// + DOMAIN
	// + DOMAIN_ID
	EntityType string `json:"entity_type"`
	// The ID of the ACL policy.
	ID string `json:"id"`
	// The latest update time.
	UpdatedAt string `json:"update_time"`
}

// BindPage is a single page maximum result representing a query by offset page.
type PolicyPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a PolicyPage struct is empty.
func (b PolicyPage) IsEmpty() (bool, error) {
	arr, err := ExtractPolicies(b)
	return len(arr) == 0, err
}

// ExtractPolicies is a method to extract the list of ACL policies.
func ExtractPolicies(r pagination.Page) ([]Policy, error) {
	var s []Policy
	err := r.(PolicyPage).Result.ExtractIntoSlicePtr(&s, "acls")
	return s, err
}

// BindResp is the structure that represents the API response of the ACL policy binding.
type BindResp struct {
	// The ID of the binding relationship.
	ID string `json:"id"`
	// The API ID.
	ApiId string `json:"api_id"`
	// The environment ID where the API is published.
	EnvId string `json:"env_id"`
	// Throttling policy ID
	PolicyId string `json:"acl_id"`
	// The creation time.
	CreatedAt string `json:"create_time"`
}

// BindPage is a single page maximum result representing a query by offset page.
type BindPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a BindPage struct is empty.
func (b BindPage) IsEmpty() (bool, error) {
	arr, err := ExtractBindInfos(b)
	return len(arr) == 0, err
}

// ExtractBinds is a method to extract the list of binding details for ACL policy.
func ExtractBindInfos(r pagination.Page) ([]AclBindApiInfo, error) {
	var s []AclBindApiInfo
	err := r.(BindPage).Result.ExtractIntoSlicePtr(&s, "apis")
	return s, err
}

// AclBindApiInfo is the structure that represents the binding details.
type AclBindApiInfo struct {
	// The API ID.
	ID string `json:"api_id"`
	// The API Name.
	Name string `json:"api_name"`
	// The API type.
	Type int `json:"api_type"`
	// The API type.
	Description string `json:"api_remark"`
	// The environment ID where the API is published.
	EnvId string `json:"env_id"`
	// The name of the environment published by the API.
	EnvName string `json:"env_name"`
	// The binding ID.
	BindId string `json:"bind_id"`
	// Group name to which the API belongs.
	GroupName string `json:"group_name"`
	// The time when the API and the policy were bound.
	BoundAt string `json:"bind_time"`
	// API publish record ID.
	PublishId string `json:"publish_id"`
}

// BatchUnbindResp is the structure that represents the API response of the ACL policy unbinding.
type BatchUnbindResp struct {
	// Number of API and throttling policy bindings that have been successfully unbound.
	SuccessCount int `json:"success_count"`
	// Unbind failed error code
	Failures []Failure `json:"failure"`
}

// Failure is the structure that represents the failure details.
type Failure struct {
	// API and ACL policy binding relationship ID that failed to unbind.
	BindId string `json:"bind_id"`
	// Unbind failed error code.
	ErrorCode string `json:"error_code"`
	// Unbind failed error message.
	ErrorMsg string `json:"error_msg"`
	// ID of the API that failed to unbind.
	ApiId string `json:"api_id"`
	// The name of the API that failed to unbind.
	ApiName string `json:"api_name"`
}
