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
