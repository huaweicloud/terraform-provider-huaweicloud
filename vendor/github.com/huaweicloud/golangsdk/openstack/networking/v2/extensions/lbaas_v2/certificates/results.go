package certificates

import (
	"github.com/huaweicloud/golangsdk"
)

type Certificate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Domain      string `json:"domain"`
	PrivateKey  string `json:"private_key"`
	Certificate string `json:"certificate"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Certificate, error) {
	s := &Certificate{}
	return s, r.ExtractInto(s)
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
