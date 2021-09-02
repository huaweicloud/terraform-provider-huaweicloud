package users

import (
	"github.com/chnsz/golangsdk"
)

// User represents a User in the OpenStack Identity Service.
type User struct {
	ID               string `json:"id"`
	DomainID         string `json:"domain_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	AreaCode         string `json:"areacode"`
	Phone            string `json:"phone"`
	Description      string `json:"description"`
	AccessMode       string `json:"access_mode"`
	Enabled          bool   `json:"enabled"`
	PasswordStatus   bool   `json:"pwd_status"`
	PasswordStrength string `json:"pwd_strength"`
	CreateAt         string `json:"create_time"`
	UpdateAt         string `json:"update_time"`
	LastLogin        string `json:"last_login_time"`
	XuserID          string `json:"xuser_id"`
	XuserType        string `json:"xuser_type"`

	// Links contains referencing links to the user.
	Links map[string]interface{} `json:"links"`
}

type userResult struct {
	golangsdk.Result
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as a User.
type GetResult struct {
	userResult
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as a User.
type CreateResult struct {
	userResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a User.
type UpdateResult struct {
	userResult
}

// Extract interprets any user results as a User.
func (r userResult) Extract() (*User, error) {
	var s struct {
		User *User `json:"user"`
	}
	err := r.ExtractInto(&s)
	return s.User, err
}
