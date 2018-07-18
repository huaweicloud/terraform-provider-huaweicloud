package stacktemplates

import (
	"encoding/json"

	"github.com/huaweicloud/golangsdk"
)

// GetResult represents the result of a Get operation.
type GetResult struct {
	golangsdk.Result
}

// Extract returns the JSON template and is called after a Get operation.
func (r GetResult) Extract() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	template, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		return nil, err
	}
	return template, nil
}
