package domains

import (
	"github.com/chnsz/golangsdk"
)

type Domain struct {
	// Domain ID
	Id string `json:"id"`
	// Domain name
	HostName string `json:"hostname"`
	// CNAME value
	PolicyId string `json:"policyid"`
	// Access Code
	AccessCode string `json:"access_code"`
	// WAF mode: 0 - disabled, 1 - enabled, -1 - bypassed.
	ProtectStatus int `json:"protect_status"`
	// Whether a domain name is connected to WAF
	AccessStatus int `json:"access_status"`
	// Protocol type
	Protocol string `json:"protocol,omitempty"`
	// Certificate ID
	CertificateId string `json:"certificateid"`
	// Certificate name
	CertificateName string `json:"certificatename"`
	//The original server information
	Servers []Server `json:"server"`
	//Whether proxy is configured
	Proxy *bool `json:"proxy"`
	// the time when the domain is created in unix timestamp
	Timestamp int `json:"timestamp"`
}

type Server struct {
	// Protocol type of the client
	FrontProtocol string `json:"front_protocol"`
	// Protocol used by WAF to forward client requests to the server
	BackProtocol string `json:"back_protocol"`
	// IP address or domain name of the web server that the client accesses.
	Address string `json:"address"`
	// Port number used by the web server
	Port int `json:"port"`
	// The type of network: ipv4, ipv6. Default: ipv4
	Type string `json:"type"`
	// VPC ID where the site is located
	VpcId string `json:"vpc_id"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a domain.
func (r commonResult) Extract() (*Domain, error) {
	var response Domain
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Domain.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Domain.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Domain.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
