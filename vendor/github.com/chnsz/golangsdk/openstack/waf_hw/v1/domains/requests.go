package domains

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToDomainCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new backup.
type CreateOpts struct {
	HostName        string       `json:"hostname" required:"true"`
	Servers         []ServerOpts `json:"server" required:"true"`
	PolicyId        string       `json:"policyid,omitempty"`
	CertificateId   string       `json:"certificateid,omitempty"`
	CertificateName string       `json:"certificatename,omitempty"`
	Proxy           *bool        `json:"proxy,omitempty"`
}

// ServerOpts contains the origin server information.
type ServerOpts struct {
	FrontProtocol string `json:"front_protocol" required:"true"`
	BackProtocol  string `json:"back_protocol" required:"true"`
	Address       string `json:"address" required:"true"`
	Port          int    `json:"port" required:"true"`
	Type          string `json:"type,omitempty"`
	VpcId         string `json:"vpc_id,omitempty"`
}

// ToDomainCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToDomainCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new Domain based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToDomainCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a Domain.
type UpdateOpts struct {
	Proxy           *bool             `json:"proxy,omitempty"`
	CertificateId   string            `json:"certificateid,omitempty"`
	CertificateName string            `json:"certificatename,omitempty"`
	Servers         []ServerOpts      `json:"server,omitempty"`
	Tls             string            `json:"tls,omitempty"`
	Cipher          string            `json:"cipher,omitempty"`
	BlockPages      []BlockPage       `json:"block_page,omitempty"`
	TrafficMarks    []TrafficMark     `json:"traffic_mark,omitempty"`
	Flag            map[string]string `json:"flag,omitempty"`
	Extend          map[string]string `json:"extend,omitempty"`
}

// BlockPage contains the alarm page information
type BlockPage struct {
	Template    string       `json:"template" required:"true"`
	CustomPages []CustomPage `json:"custom_page,omitempty"`
	RedirectUrl string       `json:"redirect_url,omitempty"`
}

// CustomPage contains the customized alarm page information
type CustomPage struct {
	StatusCode  string `json:"status_code" required:"true"`
	ContentType string `json:"content_type" required:"true"`
	Content     string `json:"content" required:"true"`
}

// TrafficMark contains the traffic identification
type TrafficMark struct {
	Sip    []string `json:"sip,omitempty"`
	Cookie string   `json:"cookie,omitempty"`
	Params string   `json:"params,omitempty"`
}

// ToDomainUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToDomainUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a Domain.The response code from api is 200
func Update(c *golangsdk.ServiceClient, domainID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToDomainUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, domainID), b, nil, reqOpt)
	return
}

// Get retrieves a particular Domain based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, reqOpt)
	return
}

// DeleteOptsBuilder allows extensions to add additional parameters to the
// delete request.
type DeleteOptsBuilder interface {
	ToDeleteQuery() (string, error)
}

// DeleteOpts contains all the values needed to delete a domain.
type DeleteOpts struct {
	// KeepPolicy specifies whether to retain the policy when deleting a domain name
	// the default value is false
	KeepPolicy bool `q:"keepPolicy"`
}

// ToDeleteQuery builds a delete request body from DeleteOpts.
func (opts DeleteOpts) ToDeleteQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// Delete will permanently delete a particular Domain based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := resourceURL(c, id)
	if opts != nil {
		var query string
		query, r.Err = opts.ToDeleteQuery()
		if r.Err != nil {
			return
		}
		url += query
	}

	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Delete(url, reqOpt)
	return
}
