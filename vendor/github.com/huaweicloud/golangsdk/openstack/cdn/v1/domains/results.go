package domains

import (
	"fmt"
	"time"

	"github.com/huaweicloud/golangsdk"
)

// sources
type DomainSources struct {
	DomainID      string `json:"domain_id"`
	IporDomain    string `json:"ip_or_domain"`
	OriginType    string `json:"origin_type"`
	ActiveStandby int    `json:"active_standby"`
}

// domain_origin_host
type DomainOriginHost struct {
	DomainID        string `json:"domain_id"`
	OriginHostType  string `json:"origin_host_type"`
	CustomizeDomain string `json:"customize_domain"`
}

// CdnDomain represents a CDN domain
type CdnDomain struct {
	// the acceleration domain name ID
	ID string `json:"id"`
	// the acceleration domain name
	DomainName string `json:"domain_name"`
	// the service type, valid values are web, download, video
	BusinessType string `json:"business_type"`
	// the domain ID of the domain name's owner
	UserDomainId string `json:"user_domain_id"`
	// the status of the acceleration domain name.
	DomainStatus string `json:"domain_status"`
	// the CNAME of the acceleration domain name
	CName string `json:"cname"`
	// the domain name or the IP address of the origin server
	Sources []DomainSources `json:"sources"`
	// the configuration information of the retrieval host
	OriginHost DomainOriginHost `json:"domain_origin_host"`
	// whether the HTTPS certificate is enabled
	HttpsStatus *int `json:"https_status"`
	// whether the status is disabled
	Disabled *int `json:"disabled"`
	// whether the status is locked
	Locked *int `json:"locked"`
	// the area covered by the accelecration service
	ServiceArea string `json:"service_area"`
	// whether range-based retrieval is enabled
	RangeStatus string `json:"range_status"`
	// a thrid-party CDN
	ThridPartCDN string `json:"third_part_cdn"`
	// the id of enterprise project
	EnterpriseProjectId string `json:"enterprise_project_id"`

	CreateTime time.Time `json:"-"`
	ModifyTime time.Time `json:"-"`
}

type OriginSources struct {
	// the domain name or the IP address of the origin server
	Sources []DomainSources `json:"sources"`
}

type commonResult struct {
	golangsdk.Result
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a CDN domain.
type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*CdnDomain, error) {
	var domain CdnDomain
	err := r.ExtractInto(&domain)

	// the get request API will response OK, even if errors occurrred.
	// so we judge domain  whether is existing
	if err == nil && domain.DomainName == "" {
		err = fmt.Errorf("The CDN domain does not exist.")
	}
	return &domain, err
}

func (r GetResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "domain")
}

// CreateResult is the result of a Create request.
type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*CdnDomain, error) {
	var domain CdnDomain
	err := r.ExtractInto(&domain)

	// the create request API will response OK, even if errors occurrred.
	// so we judge domain  whether is existing
	if err == nil && domain.DomainStatus != "configuring" {
		err = fmt.Errorf("%v", r.Body)
	}
	return &domain, err
}

func (r CreateResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "domain")
}

// DeleteResult is the result of a Delete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type DeleteResult struct {
	commonResult
}

func (r DeleteResult) Extract() (*CdnDomain, error) {
	var domain CdnDomain
	err := r.ExtractInto(&domain)

	// the delete request API will response OK, even if errors occurrred.
	// so we judge domain  whether is existing
	if err == nil && domain.DomainStatus != "deleting" {
		err = fmt.Errorf("%v", r.Body)
	}
	return &domain, err
}

func (r DeleteResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "domain")
}

// EnableResult is the result of a Enable request.
type EnableResult struct {
	commonResult
}

// DisableResult is the result of a Disable request.
type DisableResult struct {
	commonResult
}

// OriginResult is the result of a origin request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type OriginResult struct {
	commonResult
}

func (r OriginResult) Extract() (*OriginSources, error) {
	var origin OriginSources
	err := r.ExtractInto(&origin)

	return &origin, err
}

func (r OriginResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "origin")
}
