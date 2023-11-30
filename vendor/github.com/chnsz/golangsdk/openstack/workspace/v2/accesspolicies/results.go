package accesspolicies

import "github.com/chnsz/golangsdk/pagination"

type policiesResp struct {
	// List of policy.
	Policies []AccessPolicyDetailInfo `json:"policies"`
	// Total number of the policies.
	Total int `json:"total"`
}

// AccessPolicyDetailInfo is the structure that represents the access policy information.
type AccessPolicyDetailInfo struct {
	// Policy name.
	// + PRIVATE_ACCESS: Private line access
	PolicyName string `json:"policy_name"`
	// Blacklist type.
	// + INTERNET: Internet.
	BlacklistType string `json:"blacklist_type"`
	// Policy ID.
	PolicyId string `json:"policy_id"`
	// The creation time of the policy.
	CreateTime string `json:"create_time"`
}

// AccessPolicyObject is the structure that represents the policy access object information.
type AccessPolicyObject struct {
	// Object ID.
	ObjectId string `json:"object_id"`
	// Object name.
	ObjectName string `json:"object_name"`
	// Object type.
	ObjectType string `json:"object_type"`
}

// AccessPolicyPage is a single page maximum result representing a query by offset page.
type AccessPolicyPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether an AccessPolicyPage struct is empty.
func (b AccessPolicyPage) IsEmpty() (bool, error) {
	arr, err := ExtractAccessPolicies(b)
	return len(arr) == 0, err
}

// ExtractAccessPolicies is a method to extract the list of access objects.
func ExtractAccessPolicies(r pagination.Page) ([]AccessPolicyObject, error) {
	var s []AccessPolicyObject
	err := r.(AccessPolicyPage).Result.ExtractIntoSlicePtr(&s, "policy_objects_list")
	return s, err
}
