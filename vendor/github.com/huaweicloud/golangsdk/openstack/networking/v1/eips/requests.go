package eips

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//ApplyOptsBuilder is an interface by which can build the request body of public ip
//application
type ApplyOptsBuilder interface {
	ToPublicIpApplyMap() (map[string]interface{}, error)
}

//ApplyOpts is a struct which is used to create public ip
type ApplyOpts struct {
	IP                  PublicIpOpts  `json:"publicip" required:"true"`
	Bandwidth           BandwidthOpts `json:"bandwidth" required:"true"`
	EnterpriseProjectID string        `json:"enterprise_project_id,omitempty"`
}

type PublicIpOpts struct {
	Type    string `json:"type" required:"true"`
	Address string `json:"ip_address,omitempty"`
}

type BandwidthOpts struct {
	Name       string `json:"name,omitempty"`
	Size       int    `json:"size,omitempty"`
	Id         string `json:"id,omitempty"`
	ShareType  string `json:"share_type" required:"true"`
	ChargeMode string `json:"charge_mode,omitempty"`
}

func (opts ApplyOpts) ToPublicIpApplyMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Apply is a method by which can access to apply the public ip
func Apply(client *golangsdk.ServiceClient, opts ApplyOptsBuilder) (r ApplyResult) {
	b, err := opts.ToPublicIpApplyMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//Get is a method by which can get the detailed information of public ip
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

//Delete is a method by which can be able to delete a private ip
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

//UpdateOptsBuilder is an interface by which can be able to build the request
//body
type UpdateOptsBuilder interface {
	ToPublicIpUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the request body of update method
type UpdateOpts struct {
	PortID string `json:"port_id,omitempty"`
}

func (opts UpdateOpts) ToPublicIpUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "publicip")
}

//Update is a method which can be able to update the port of public ip
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPublicIpUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type ListOpts struct {
	// Specifies the resource ID of pagination query. If the parameter
	// is left blank, only resources on the first page are queried.
	Marker string `q:"marker"`

	// Specifies the number of records returned on each page. The
	// value ranges from 0 to intmax.
	Limit int `q:"limit"`

	// Value range: 4, 6, respectively, to create ipv4 and ipv6, when not created ipv4 by default
	IPVersion int `q:"ip_version"`

	// Associated port id
	PortId string `q:"port_id"`

	// Public IP address
	PublicIp string `q:"public_ip_address"`

	// enterprise_project_id
	// You can use this field to filter the elastic public IP under an enterprise project.
	EnterpriseProjectId string `q:"enterprise_project_id"`
}

type ListOptsBuilder interface {
	ToListPublicIPQuery() (string, error)
}

func (opts ListOpts) ToListPublicIPQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToListPublicIPQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PublicIPPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
