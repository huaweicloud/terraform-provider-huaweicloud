package customer_gateways

import "github.com/chnsz/golangsdk/pagination"

type CustomerGateway struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	IDType        string        `json:"id_type"`
	IDValue       string        `json:"id_value"`
	BGPAsn        int           `json:"bgp_asn"`
	IP            string        `json:"ip"`
	RouteMode     string        `json:"route_mode"`
	CaCertificate CaCertificate `json:"ca_certificate"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
}

type CaCertificate struct {
	SerialNumber       string `json:"serial_number"`
	SignatureAlgorithm string `json:"signature_algorithm"`
	Issuer             string `json:"issuer"`
	Subject            string `json:"subject"`
	ExpireTime         string `json:"expire_time"`
	IsUpdatable        bool   `json:"is_updatable"`
}

type listResp struct {
	CustomerGateways []CustomerGateway `json:"customer_gateways"`
	RequestId        string            `json:"request_id"`
	PageInfo         pageInfo          `json:"page_info"`
	TotalCount       int64             `json:"total_count"`
}

type pageInfo struct {
	NextMarker   string `json:"next_marker"`
	CurrentCount int    `json:"current_count"`
}

type CustomerGatewaysPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a ListResult is empty.
func (r CustomerGatewaysPage) IsEmpty() (bool, error) {
	resp, err := extractCustomerGateways(r)
	return len(resp) == 0, err
}

// LastMarker returns the last marker index in a ListResult.
func (r CustomerGatewaysPage) LastMarker() (string, error) {
	resp, err := extractPageInfo(r)
	if err != nil {
		return "", err
	}
	return resp.NextMarker, nil
}

// extractCustomerGateways is a method which to extract the response to a CustomerGateway list.
func extractCustomerGateways(r pagination.Page) ([]CustomerGateway, error) {
	var s listResp
	err := r.(CustomerGatewaysPage).Result.ExtractInto(&s)
	return s.CustomerGateways, err
}

// extractPageInfo is a method which to extract the response of the page information.
func extractPageInfo(r pagination.Page) (*pageInfo, error) {
	var s listResp
	err := r.(CustomerGatewaysPage).Result.ExtractInto(&s)
	return &s.PageInfo, err
}
