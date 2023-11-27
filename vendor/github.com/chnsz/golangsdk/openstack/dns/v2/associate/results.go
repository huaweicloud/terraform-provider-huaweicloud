package associate

import "github.com/chnsz/golangsdk"

type AssociateResult struct {
	commonResult
}

type DisAssociateResult struct {
	commonResult
}

type commonResult struct {
	golangsdk.Result
}

type ResponseRouter struct {
	RouterID     string `json:"router_id"`
	RouterRegion string `json:"router_region"`
	Status       string `json:"status"`
}

func (r commonResult) Extract() (*ResponseRouter, error) {
	var s *ResponseRouter
	err := r.ExtractInto(&s)
	return s, err
}
