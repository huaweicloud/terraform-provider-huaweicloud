package whitelists

import (
	"github.com/huaweicloud/golangsdk"
)

type Whitelist struct {
	ID              string `json:"id"`
	TenantId        string `json:"tenant_id"`
	ListenerId      string `json:"listener_id"`
	EnableWhitelist bool   `json:"enable_whitelist"`
	Whitelist       string `json:"whitelist"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Whitelist, error) {
	s := &Whitelist{}
	return s, r.ExtractInto(s)
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "whitelist")
}

type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}
