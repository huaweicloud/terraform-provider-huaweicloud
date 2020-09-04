package whitelists

import (
	"github.com/huaweicloud/golangsdk"
)

// WhitelistGroup is a struct that contains the whitelist parameters.
type WhitelistGroup struct {
	// the group name
	GroupName string `json:"group_name"`
	// list of IP address or range
	IPList []string `json:"ip_list"`
}

// Whitelist is a struct that contains all the whitelist parameters.
type Whitelist struct {
	// instance id
	InstanceID string `json:"instance_id"`
	// enable or disable the whitelists
	Enable bool `json:"enable_whitelist"`
	// list of whitelist groups
	Groups []WhitelistGroup `json:"whitelist"`
}

// PutResult is a struct from which can get the result of put method
type PutResult struct {
	golangsdk.ErrResult
}

// WhitelistResult is a struct from which can get the result of get method
type WhitelistResult struct {
	golangsdk.Result
}

// Extract from WhitelistResult
func (r WhitelistResult) Extract() (*Whitelist, error) {
	var s Whitelist
	err := r.Result.ExtractInto(&s)
	return &s, err
}
