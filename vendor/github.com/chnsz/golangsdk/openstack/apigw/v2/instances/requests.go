package instances

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts allows to create an APIG dedicated instance using given parameters.
type CreateOpts struct {
	// Name of the APIG dedicated instance. The name can contains of 3 to 64 characters.
	Name string `json:"instance_name" required:"true"`
	// Edition of the APIG dedicated instance. Currently, the editions are support:
	// (IPv4): BASIC, PROFESSIONAL, ENTERPRISE, PLATINUM
	// (IPv6): BASIC_IPV6, PROFESSIONAL_IPV6, ENTERPRISE_IPV6, PLATINUM_IPV6
	Edition string `json:"spec_id" required:"true"`
	// VPC ID.
	VpcId string `json:"vpc_id" required:"true"`
	// Subnet network ID.
	SubnetId string `json:"subnet_id" required:"true"`
	// ID of the security group to which the APIG dedicated instance belongs to.
	SecurityGroupId string `json:"security_group_id" required:"true"`
	// ID of the APIG dedicated instance, which will be automatically generated if you do not specify this parameter.
	Id string `json:"instance_id,omitempty"`
	// Description about the APIG dedicated instance.
	Description string `json:"description,omitempty"`
	// Start time of the maintenance time window in the format "xx:00:00".
	// The value of xx can be 02, 06, 10, 14, 18, or 22.
	MaintainBegin string `json:"maintain_begin,omitempty"`
	// End time of the maintenance time window in the format "xx:00:00".
	// There is a 4-hour difference between the start time and end time.
	MaintainEnd string `json:"maintain_end,omitempty"`
	// EIP ID.
	EipId string `json:"eip_id,omitempty"`
	// Outbound access bandwidth. This parameter is required if public outbound access is enabled for the APIG
	// dedicated instance.
	// Zero means turn off the egress access.
	BandwidthSize int `json:"bandwidth_size"`
	// Enterprise project ID. This parameter is required if you are using an enterprise account.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// AZs.
	AvailableZoneIds []string `json:"available_zone_ids,omitempty"`
	// Whether public access with an IPv6 address is supported.
	Ipv6Enable bool `json:"ipv6_enable,omitempty"`
}

type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a APIG dedicated instance.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), reqBody, &r.Body, nil)
	return
}

// Get is a method to obtain the specified APIG dedicated instance according to the instance Id.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// ID of the APIG dedicated instance.
	Id string `q:"instance_id"`
	// Name of the APIG dedicated instance.
	Name string `q:"instance_name"`
	// Instance status.
	Status string `q:"status"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	Limit int `q:"limit"`
}

type ListOptsBuilder interface {
	ToInstanceListQuery() (string, error)
}

func (opts ListOpts) ToInstanceListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more APIG dedicated instance according to the query parameters.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.SinglePageBase(r)}
	})
}

// UpdateOpts allows to update an existing APIG dedicated instance using given parameters.
type UpdateOpts struct {
	// Description about the APIG dedicated instance.
	Description string `json:"description,omitempty"`
	// Start time of the maintenance time window in the format "xx:00:00".
	// The value of xx can be 02, 06, 10, 14, 18, or 22.
	MaintainBegin string `json:"maintain_begin,omitempty"`
	// End time of the maintenance time window in the format "xx:00:00".
	// There is a 4-hour difference between the start time and end time.
	MaintainEnd string `json:"maintain_end,omitempty"`
	// Description about the APIG dedicated instance.
	Name string `json:"instance_name,omitempty"`
	// ID of the security group to which the APIG dedicated instance belongs to.
	SecurityGroupId string `json:"security_group_id,omitempty"`
}

type UpdateOptsBuilder interface {
	ToInstanceUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToInstanceUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is a method by which to update an existing APIG dedicated instance.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToInstanceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete is a method to delete an existing APIG dedicated instance
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

// EgressAccessOpts allows the bandwidth size of an existing APIG dedicated instance to be updated with the given
// parameters.
type EgressAccessOpts struct {
	// Outbound access bandwidth, in Mbit/s.
	BandwidthSize string `json:"bandwidth_size,omitempty"`
}

type EgressAccessOptsBuilder interface {
	ToEgressAccessMap() (map[string]interface{}, error)
}

func (opts EgressAccessOpts) ToEgressAccessMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// EnableEgressAccess is a method by which to enable the egress access of an existing APIG dedicated instance.
func EnableEgressAccess(client *golangsdk.ServiceClient, id string, opts EgressAccessOptsBuilder) (r EnableEgressResult) {
	reqBody, err := opts.ToEgressAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(egressURL(client, id), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// UpdateEgressBandwidth is a method by which to update the egress bandwidth size of an existing APIG dedicated instance.
func UpdateEgressBandwidth(client *golangsdk.ServiceClient, id string, opts EgressAccessOptsBuilder) (r UdpateEgressResult) {
	reqBody, err := opts.ToEgressAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(egressURL(client, id), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// DisableEgressAccess is a method by which to disable the egress access of an existing APIG dedicated instance.
func DisableEgressAccess(client *golangsdk.ServiceClient, id string) (r DisableEgressResult) {
	_, r.Err = client.Delete(egressURL(client, id), nil)
	return
}

// IngressAccessOpts allows binding and updating the eip associated with an existing APIG dedicated instance with the
// given parameters.
type IngressAccessOpts struct {
	// EIP ID
	EipId string `json:"eip_id,omitempty"`
}

type IngressAccessOptsBuilder interface {
	ToEnableIngressAccessMap() (map[string]interface{}, error)
}

func (opts IngressAccessOpts) ToEnableIngressAccessMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// UpdateIngressAccess is a method to bind and update the eip associated with an existing APIG dedicated instance.
func EnableIngressAccess(client *golangsdk.ServiceClient, id string, opts IngressAccessOptsBuilder) (r EnableIngressResult) {
	reqBody, err := opts.ToEnableIngressAccessMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(ingressURL(client, id), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// DisableIngressAccess is a method to unbind the eip associated with an existing APIG dedicated instance.
func DisableIngressAccess(client *golangsdk.ServiceClient, id string) (r DisableIngressResult) {
	_, r.Err = client.Delete(ingressURL(client, id), &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
