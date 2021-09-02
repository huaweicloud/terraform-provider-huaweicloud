package apigroups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// GetResult represents a result of the Get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents a result of the Update operation.
type UpdateResult struct {
	commonResult
}

type Group struct {
	// List of independent domain names bound to the API group.
	UrlDomians []UrlDomian `json:"url_domains"`
	// Time when the API group was last modified.
	UpdateTime string `json:"update_time"`
	// API group name.
	Name string `json:"name"`
	// Indicates whether the API group has been listed on the marketplace.
	//     1: listed
	//     2: not listed
	//     3: under review.
	OnSellStatus int `json:"on_sell_status"`
	// Description.
	Description string `json:"remark"`
	// Subdomain name that API Gateway automatically allocates to the API group.
	Subdomain string `json:"sl_domain"`
	// Subdomain names that API Gateway automatically allocates to the API group.
	Subdomains []string `json:"sl_domains"`
	// ID.
	Id string `json:"id"`
	// Registraion time.
	RegistraionTime string `json:"register_time"`
	// group status.
	//     1: valid
	Status int `json:"status"`
	// Indicates whether the API group is the default group.
	IsDefault int `json:"is_default"`
}

type UrlDomian struct {
	// Domain ID.
	Id string `json:"id"`
	// Domain ID.
	DomainName string `json:"domain"`
	// CNAME resolution status of the domain name.
	//     1: not resolved
	//     2: resolving
	//     3: resolved
	//     4: resolving failed
	ResolutionStatus int `json:"cname_status"`
	// SSL certificate ID.
	SSLId string `json:"ssl_id"`
	// SSL certificate name.
	SSLName string `json:"ssl_name"`
	// Minimum SSL version. TLS 1.1 and TLS 1.2 are supported.
	// Enumeration values:
	//     TLSv1.1
	//     TLSv1.2
	MinSSLVersion string `json:"min_ssl_version"`
}

func (r commonResult) Extract() (*Group, error) {
	var s Group
	err := r.ExtractInto(&s)
	return &s, err
}

// GroupPage represents the response pages of the List operation.
type GroupPage struct {
	pagination.SinglePageBase
}

func ExtractGroups(r pagination.Page) ([]Group, error) {
	var s []Group
	err := r.(GroupPage).Result.ExtractIntoSlicePtr(&s, "groups")
	return s, err
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	golangsdk.ErrResult
}
