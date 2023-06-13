package domains

import (
	"github.com/chnsz/golangsdk"
)

// GetResponse response
type GetResponse struct {
	Domains []Domain `json:"domains"`
}

type Domain struct {
	//Domain ID
	Id string `json:"id"`
	//Domain Name
	Name string `json:"name"`
	//Domain Description
	Description string `json:"description"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*GetResponse, error) {
	var s GetResponse
	err := r.Result.ExtractInto(&s)
	return &s, err
}
