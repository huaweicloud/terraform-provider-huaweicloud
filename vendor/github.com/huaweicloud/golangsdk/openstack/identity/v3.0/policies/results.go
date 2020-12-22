package policies

import (
	"github.com/huaweicloud/golangsdk"
)

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"display_name"`
	Catalog     string `json:"catalog"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Policy      Policy `json:"policy" required:"true"`
	DomainId    string `json:"domain_id"`
	References  int    `json:"references"`
}

type roleResult struct {
	golangsdk.Result
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as a Role.
type GetResult struct {
	roleResult
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as a Role
type CreateResult struct {
	roleResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Role.
type UpdateResult struct {
	roleResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// Extract interprets any roleResults as a Role.
func (r roleResult) Extract() (*Role, error) {
	var s struct {
		Role *Role `json:"role"`
	}
	err := r.ExtractInto(&s)
	return s.Role, err
}
