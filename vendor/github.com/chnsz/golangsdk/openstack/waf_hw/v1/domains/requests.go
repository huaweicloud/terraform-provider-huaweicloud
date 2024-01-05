package domains

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/utils"
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
	HostName            string            `json:"hostname" required:"true"`
	Servers             []ServerOpts      `json:"server" required:"true"`
	PolicyId            string            `json:"policyid,omitempty"`
	CertificateId       string            `json:"certificateid,omitempty"`
	CertificateName     string            `json:"certificatename,omitempty"`
	Proxy               *bool             `json:"proxy,omitempty"`
	PaidType            string            `json:"paid_type,omitempty"`
	EnterpriseProjectId string            `q:"enterprise_project_id" json:"-"`
	ForwardHeaderMap    map[string]string `json:"forward_header_map,omitempty"`
	Description         string            `json:"description,omitempty"`
	LbAlgorithm         string            `json:"lb_algorithm,omitempty"`
	WebTag              string            `json:"web_tag,omitempty"`
}

// ServerOpts contains the origin server information.
type ServerOpts struct {
	FrontProtocol string `json:"front_protocol" required:"true"`
	BackProtocol  string `json:"back_protocol" required:"true"`
	Address       string `json:"address" required:"true"`
	Port          int    `json:"port" required:"true"`
	Type          string `json:"type,omitempty"`
	VpcId         string `json:"vpc_id,omitempty"`
	Weight        int    `json:"weight,omitempty"`
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
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c)+query.String(), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a Domain.
type UpdateOpts struct {
	Proxy               *bool             `json:"proxy,omitempty"`
	CertificateId       string            `json:"certificateid,omitempty"`
	CertificateName     string            `json:"certificatename,omitempty"`
	Servers             []ServerOpts      `json:"server,omitempty"`
	Tls                 string            `json:"tls,omitempty"`
	Cipher              string            `json:"cipher,omitempty"`
	Http2Enable         *bool             `json:"http2_enable,omitempty"`
	Ipv6Enable          *bool             `json:"ipv6_enable,omitempty"`
	WebTag              *string           `json:"web_tag,omitempty"`
	ExclusiveIp         *bool             `json:"exclusive_ip,omitempty"`
	PaidType            string            `json:"paid_type,omitempty"`
	BlockPage           *BlockPage        `json:"block_page,omitempty"`
	TrafficMark         *TrafficMark      `json:"traffic_mark,omitempty"`
	Flag                *Flag             `json:"flag,omitempty"`
	Extend              map[string]string `json:"extend,omitempty"`
	TimeoutConfig       *TimeoutConfig    `json:"timeout_config,omitempty"`
	ForwardHeaderMap    map[string]string `json:"forward_header_map,omitempty"`
	EnterpriseProjectId string            `q:"enterprise_project_id" json:"-"`
	Description         *string           `json:"description,omitempty"`
	LbAlgorithm         *string           `json:"lb_algorithm,omitempty"`
}

// BlockPage contains the alarm page information
type BlockPage struct {
	Template    string      `json:"template" required:"true"`
	CustomPage  *CustomPage `json:"custom_page,omitempty"`
	RedirectUrl string      `json:"redirect_url,omitempty"`
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

type TimeoutConfig struct {
	ConnectTimeout *int `json:"connect_timeout,omitempty"`
	SendTimeout    *int `json:"send_timeout,omitempty"`
	ReadTimeout    *int `json:"read_timeout,omitempty"`
}

type Flag struct {
	Pci3ds   string `json:"pci_3ds,omitempty"`
	PciDss   string `json:"pci_dss,omitempty"`
	Cname    string `json:"cname,omitempty"`
	IsDualAz string `json:"is_dual_az,omitempty"`
	Ipv6     string `json:"ipv6,omitempty"`
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
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Patch(resourceURL(c, domainID)+query.String(), b, nil, reqOpt)
	return
}

// updateProtectStatusOpts the struct for updating the protect status of domain.
type updateProtectStatusOpts struct {
	ProtectStatus *int `json:"protect_status" required:"true"`
}

// UpdateProtectStatus update the protect status of domain.
func UpdateProtectStatus(c *golangsdk.ServiceClient, protectStatus int, instanceId, epsId string) (*HostProtectStatus, error) {
	opts := updateProtectStatusOpts{
		ProtectStatus: &protectStatus,
	}

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(protectStatusURL(c, instanceId)+utils.GenerateEpsIDQuery(epsId), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r HostProtectStatus
		err := rst.ExtractInto(&r)
		return &r, err
	}
	return nil, err
}

// Get retrieves a particular Domain based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	return GetWithEpsID(c, id, "")
}

func GetWithEpsID(c *golangsdk.ServiceClient, id, epsID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	}
	_, r.Err = c.Get(resourceURL(c, id)+utils.GenerateEpsIDQuery(epsID), &r.Body, reqOpt)
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
	KeepPolicy          bool   `q:"keepPolicy"`
	EnterpriseProjectId string `q:"enterprise_project_id"`
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
