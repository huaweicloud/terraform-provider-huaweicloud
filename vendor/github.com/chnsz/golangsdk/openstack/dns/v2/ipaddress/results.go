package ipaddress

import "github.com/chnsz/golangsdk"

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type ListResult struct {
	commonResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type ListObject struct {
	Status     string `json:"status"`
	ID         string `json:"id"`
	IP         string `json:"ip"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	SubnetID   string `json:"subnet_id"`
	ErrorInfo  string `json:"error_info"`
}

type EndpointResp struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	VpcID             string `json:"vpc_id"`
	IPAddressCount    int    `json:"ipaddress_count"`
	ResolverRuleCount int    `json:"resolver_rule_count"`
	CreateTime        string `json:"create_time"`
	UpdateTime        string `json:"update_time"`
}

func (c CreateResult) Extract() (EndpointResp, error) {
	type responseBody struct {
		EndpointResp `json:"endpoint"`
	}
	var r responseBody
	err := c.ExtractInto(&r)
	return r.EndpointResp, err
}

func (l ListResult) Extract() ([]ListObject, error) {
	type result struct {
		IPAddress []ListObject `json:"ipaddresses"`
	}
	var r result
	err := l.ExtractInto(&r)
	return r.IPAddress, err
}
