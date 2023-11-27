package resolverrule

import "github.com/chnsz/golangsdk"

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type ListResult struct {
	commonResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type ResponseBody struct {
	ResolverRule `json:"resolver_rule"`
}

type ResolverRule struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	DomainName     string      `json:"domain_name"`
	EndpointID     string      `json:"endpoint_id"`
	Status         string      `json:"status"`
	RuleType       string      `json:"rule_type"`
	IPAddressCount int         `json:"ipaddress_count"`
	IPAddresses    []IPAddress `json:"ipaddresses"`
	Routers        []Router    `json:"routers"`
	CreatedAt      string      `json:"create_time"`
	UpdatedAt      string      `json:"update_time"`
}

type Router struct {
	RouterID     string `json:"router_id"`
	RouterRegion string `json:"router_region"`
	Status       string `json:"status"`
}

func (lr ListResult) Extract() ([]ResolverRule, error) {
	type Metadata struct {
		TotalCount int `json:"total_count"`
	}
	var l struct {
		ResolverRules []ResolverRule `json:"resolver_rules"`
		Metadata      `json:"metedata"`
	}
	err := lr.Result.ExtractInto(&l)
	return l.ResolverRules, err
}

func (r commonResult) Extract() (*ResponseBody, error) {
	var s *ResponseBody
	err := r.ExtractInto(&s)
	return s, err
}
