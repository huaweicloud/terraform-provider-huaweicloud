package groups

import (
	"time"

	"github.com/huaweicloud/golangsdk"
)

type UrlDomain struct {
	ID          string `json:"id"`
	Domain      string `json:"domain"`
	CnameStatus int    `json:"cname_status"`
	SslID       string `json:"ssl_id"`
	SslName     string `json:"ssl_name"`
}

// Group contains all the information associated with a API group.
type Group struct {
	// Unique identifier for the Group.
	ID string `json:"id"`
	// Human-readable display name for the Group.
	Name string `json:"name"`
	// Description of the API group.
	Remark string `json:"remark"`
	// Status of the Group.
	Status int `json:"status"`
	// Indicates whether the API group has been listed on the marketplace.
	OnSellStatus int `json:"on_sell_status"`
	// Subdomain name automatically allocated by the system to the API group.
	SlDomain string `json:"sl_domain"`
	// Total number of times all APIs in the API group can be accessed.
	CallLimits int `json:"call_limits"`
	// The type of Group to create, either SATA or SSD.
	TimeInterval int `json:"time_interval"`
	// Time unit for limiting the number of API calls
	TimeUnit string `json:"time_unit"`
	// List of independent domain names bound to the API group
	UrlDomains []UrlDomain `json:"url_domains"`
	// Time when the API group is created
	RegisterTime time.Time `json:"-"`
	// Time when the API group was last modified
	UpdateTime time.Time `json:"-"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract will get the Group object out of the commonResult object.
func (r commonResult) Extract() (*Group, error) {
	var s Group
	err := r.ExtractInto(&s)
	return &s, err
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}
