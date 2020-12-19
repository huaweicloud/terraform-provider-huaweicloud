package cloudservers

import (
	"encoding/base64"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateOpts struct {
	ImageRef string `json:"imageRef" required:"true"`

	FlavorRef string `json:"flavorRef" required:"true"`

	Name string `json:"name" required:"true"`

	UserData []byte `json:"-"`

	// AdminPass sets the root user password. If not set, a randomly-generated
	// password will be created and returned in the response.
	AdminPass string `json:"adminPass,omitempty"`

	KeyName string `json:"key_name,omitempty"`

	VpcId string `json:"vpcid" required:"true"`

	Nics []Nic `json:"nics" required:"true"`

	PublicIp *PublicIp `json:"publicip,omitempty"`

	Count int `json:"count,omitempty"`

	IsAutoRename *bool `json:"isAutoRename,omitempty"`

	RootVolume RootVolume `json:"root_volume" required:"true"`

	DataVolumes []DataVolume `json:"data_volumes,omitempty"`

	SecurityGroups []SecurityGroup `json:"security_groups,omitempty"`

	AvailabilityZone string `json:"availability_zone" required:"true"`

	ExtendParam *ServerExtendParam `json:"extendparam,omitempty"`

	MetaData *MetaData `json:"metadata,omitempty"`

	SchedulerHints *SchedulerHints `json:"os:scheduler_hints,omitempty"`

	Tags []string `json:"tags,omitempty"`

	ServerTags []ServerTags `json:"server_tags,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToServerCreateMap() (map[string]interface{}, error)
}

// ToServerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToServerCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.UserData != nil {
		var userData string
		if _, err := base64.StdEncoding.DecodeString(string(opts.UserData)); err != nil {
			userData = base64.StdEncoding.EncodeToString(opts.UserData)
		} else {
			userData = string(opts.UserData)
		}
		b["user_data"] = &userData
	}

	return map[string]interface{}{"server": b}, nil
}

type Nic struct {
	SubnetId string `json:"subnet_id" required:"true"`

	IpAddress string `json:"ip_address,omitempty"`
}

type PublicIp struct {
	Id string `json:"id,omitempty"`

	Eip *Eip `json:"eip,omitempty"`
}

type Eip struct {
	IpType string `json:"iptype" required:"true"`

	BandWidth *BandWidth `json:"bandwidth" required:"true"`

	ExtendParam *EipExtendParam `json:"extendparam,omitempty"`
}

type BandWidth struct {
	Size int `json:"size,omitempty"`

	ShareType string `json:"sharetype" required:"true"`

	ChargeMode string `json:"chargemode,omitempty"`

	Id string `json:"id,omitempty"`
}

type EipExtendParam struct {
	ChargingMode string `json:"chargingMode,omitempty"`
}

type RootVolume struct {
	VolumeType string `json:"volumetype" required:"true"`

	Size int `json:"size,omitempty"`

	ExtendParam *VolumeExtendParam `json:"extendparam,omitempty"`
}

type DataVolume struct {
	VolumeType string `json:"volumetype" required:"true"`

	Size int `json:"size" required:"true"`

	MultiAttach *bool `json:"multiattach,omitempty"`

	PassThrough *bool `json:"hw:passthrough,omitempty"`

	Extendparam *VolumeExtendParam `json:"extendparam,omitempty"`
}

type VolumeExtendParam struct {
	SnapshotId string `json:"snapshotId,omitempty"`
}

type ServerExtendParam struct {
	ChargingMode string `json:"chargingMode,omitempty"`

	RegionID string `json:"regionID,omitempty"`

	PeriodType string `json:"periodType,omitempty"`

	PeriodNum int `json:"periodNum,omitempty"`

	IsAutoRenew string `json:"isAutoRenew,omitempty"`

	IsAutoPay string `json:"isAutoPay,omitempty"`

	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`

	SupportAutoRecovery string `json:"support_auto_recovery,omitempty"`
}

type MetaData struct {
	OpSvcUserId string `json:"op_svc_userid,omitempty"`
	AgencyName  string `json:"agency_name,omitempty"`
}

type SecurityGroup struct {
	ID string `json:"id" required:"true"`
}

type SchedulerHints struct {
	Group       string `json:"group,omitempty"`
	FaultDomain string `json:"fault_domain,omitempty"`

	// Specifies whether the ECS is created on a Dedicated Host (DeH) or in a shared pool.
	Tenancy string `json:"tenancy,omitempty"`

	// DedicatedHostID specifies a DeH ID.
	DedicatedHostID string `json:"dedicated_host_id,omitempty"`
}

type ServerTags struct {
	Key   string `json:"key" required:"true"`
	Value string `json:"value,omitempty"`
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// CreatePrePaid requests a server to be provisioned to the user in the current tenant.
func CreatePrePaid(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r OrderResult) {
	reqBody, err := opts.ToServerCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Get retrieves a particular Server based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return
}

type DeleteOpts struct {
	Servers        []Server `json:"servers" required:"true"`
	DeletePublicIP bool     `json:"delete_publicip,omitempty"`
	DeleteVolume   bool     `json:"delete_volume,omitempty"`
}

type Server struct {
	Id string `json:"id" required:"true"`
}

// ToServerDeleteMap assembles a request body based on the contents of a
// DeleteOpts.
func (opts DeleteOpts) ToServerDeleteMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Delete requests a server to be deleted to the user in the current tenant.
func Delete(client *golangsdk.ServiceClient, opts DeleteOpts) (r JobResult) {
	reqBody, err := opts.ToServerDeleteMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

type DeleteOrderOpts struct {
	ResourceIds []string `json:"resourceIds" required:"true"`
	UnSubType   int      `json:"unSubType" required:"true"`
}

// ToServerDeleteOrderMap assembles a request body based on the contents of a
// DeleteOrderOpts.
func (opts DeleteOrderOpts) ToServerDeleteOrderMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// DeleteOrder requests a server to be deleted to the user in the current tenant.
func DeleteOrder(client *golangsdk.ServiceClient, opts DeleteOrderOpts) (r DeleteOrderResult) {
	reqBody, err := opts.ToServerDeleteOrderMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteOrderURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToServerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// Name of the server as a string; can be queried with regular expressions.
	// Realize that ?name=bob returns both bob and bobb. If you need to match bob
	// only, you can use a regular expression matching the syntax of the
	// underlying database server implemented for Compute.
	Name string `q:"name"`

	// Flavor is the name of the flavor in URL format.
	Flavor string `q:"flavor"`

	// Status is the value of the status of the server so that you can filter on
	// "ACTIVE" for example.
	Status string `q:"status"`

	// Specifies the ECS that is bound to an enterprise project.
	EnterpriseProjectID string `q:"enterprise_project_id"`

	// Indicates the filtering result for IPv4 addresses, which are fuzzy matched.
	// These IP addresses are private IP addresses of the ECS.
	IP string `q:"ip"`

	// Specifies the maximum number of ECSs on one page.
	// Each page contains 25 ECSs by default, and a maximum of 1000 ECSs are returned.
	Limit int `q:"limit"`

	// Specifies a page number. The default value is 1.
	// The value must be greater than or equal to 0. If the value is 0, the first page is displayed.
	Offset int `q:"offset"`
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServerListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list servers accessible to you.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToServerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type ResizeOpts struct {
	FlavorRef   string             `json:"flavorRef" required:"true"`
	Mode        string             `json:"mode,omitempty"`
	ExtendParam *ResizeExtendParam `json:"extendparam,omitempty"`
}

type ResizeExtendParam struct {
	AutoPay string `json:"isAutoPay,omitempty"`
}

// ResizeOptsBuilder allows extensions to add additional parameters to the
// Resize request.
type ResizeOptsBuilder interface {
	ToServerResizeMap() (map[string]interface{}, error)
}

// ToServerResizeMap assembles a request body based on the contents of a
// ResizeOpts.
func (opts ResizeOpts) ToServerResizeMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "resize")
}

// Resize requests a server to be resizeed.
func Resize(client *golangsdk.ServiceClient, opts ResizeOptsBuilder, serverId string) (r JobResult) {
	reqBody, err := opts.ToServerResizeMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(resizeURL(client, serverId), reqBody, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
