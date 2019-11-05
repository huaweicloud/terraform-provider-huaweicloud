package domains

import (
	"github.com/huaweicloud/golangsdk"
)

// ExtensionOpts allows extensions to add parameters to some requests
// the possible requests include get,delete,enable or disable.
type ExtensionOpts struct {
	// specifies the enterprise_project_id.
	EnterpriseProjectId string `q:"enterprise_project_id"`
}

// ToExtensionQuery formats a ExtensionOpts into a query string.
func (opts ExtensionOpts) ToExtensionQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// SourcesOpts specifies the domain name or the IP address of the origin server
type SourcesOpts struct {
	IporDomain    string `json:"ip_or_domain" required:"true"`
	OriginType    string `json:"origin_type" required:"true"`
	ActiveStandby int    `json:"active_standby" required:"true"`
}

// CreateOpts specifies the attributes used to create a CDN domain.
type CreateOpts struct {
	// the acceleration domain name, the length of a label is within 50 characters.
	DomainName string `json:"domain_name" required:"true"`
	// the service type, valid values are web, downlaod and video
	BusinessType string `json:"business_type" required:"true"`
	// the domain name or the IP address of the origin server
	Sources []SourcesOpts `json:"sources" required:"true"`
	// the enterprise project ID
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToCdnDomainCreateMap() (map[string]interface{}, error)
}

// ToCdnDomainCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToCdnDomainCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "domain")
}

// OriginOpts specifies the attributes used to modify the orogin server.
type OriginOpts struct {
	// the domain name or the IP address of the origin server
	Sources []SourcesOpts `json:"sources" required:"true"`
}

// OriginOptsBuilder allows extensions to add additional parameters to the
// Origin request.
type OriginOptsBuilder interface {
	ToCdnDomainOriginMap() (map[string]interface{}, error)
}

// ToCdnDomainOriginMap assembles a request body based on the contents of a
// OriginOpts.
func (opts OriginOpts) ToCdnDomainOriginMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "origin")
}

// Create implements a CDN domain create request.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToCdnDomainCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Get retrieves a particular CDN domain based on its unique ID.
func Get(client *golangsdk.ServiceClient, id string, opts *ExtensionOpts) (r GetResult) {
	url := getURL(client, id)
	if opts != nil {
		query, err := opts.ToExtensionQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Delete requests a CDN domain to be deleted to the user in the current tenant.
func Delete(client *golangsdk.ServiceClient, id string, opts *ExtensionOpts) (r DeleteResult) {
	url := deleteURL(client, id)
	if opts != nil {
		query, err := opts.ToExtensionQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	// Delete requests will response 'domain' body, so we use DeleteWithResponse
	_, r.Err = client.DeleteWithResponse(url, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Enable implements a CDN domain enable request.
func Enable(client *golangsdk.ServiceClient, id string, opts *ExtensionOpts) (r EnableResult) {
	url := enableURL(client, id)
	if opts != nil {
		query, err := opts.ToExtensionQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Put(url, nil, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Disable implements a CDN domain disable request.
func Disable(client *golangsdk.ServiceClient, id string, opts *ExtensionOpts) (r DisableResult) {
	url := disableURL(client, id)
	if opts != nil {
		query, err := opts.ToExtensionQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	_, r.Err = client.Put(url, nil, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Modifying Information About the Origin Server
func Origin(client *golangsdk.ServiceClient, id string, opts *ExtensionOpts, req OriginOptsBuilder) (r OriginResult) {
	url := originURL(client, id)
	if opts != nil {
		// build url with enterprise_project_id
		query, err := opts.ToExtensionQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	// build request body
	reqBody, err := req.ToCdnDomainOriginMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(url, reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
