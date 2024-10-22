package cloudservers

import (
	"encoding/base64"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
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

	AvailabilityZone string `json:"availability_zone,omitempty"`

	ExtendParam *ServerExtendParam `json:"extendparam,omitempty"`

	MetaData *MetaData `json:"metadata,omitempty"`

	SchedulerHints *SchedulerHints `json:"os:scheduler_hints,omitempty"`

	Tags []string `json:"tags,omitempty"`

	ServerTags []tags.ResourceTag `json:"server_tags,omitempty"`

	Description string `json:"description,omitempty"`

	AutoTerminateTime string `json:"auto_terminate_time,omitempty"`
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
	SubnetId  string `json:"subnet_id" required:"true"`
	IpAddress string `json:"ip_address,omitempty"`

	// enable ipv6 or not
	Ipv6Enable bool `json:"ipv6_enable,omitempty"`
	// bandWidth id when ipv6 is enabled
	BandWidth *Ipv6BandWidth `json:"ipv6_bandwidth,omitempty"`
}

type Ipv6BandWidth struct {
	ID string `json:"id,omitempty"`
}

type PublicIp struct {
	Id string `json:"id,omitempty"`

	Eip *Eip `json:"eip,omitempty"`

	DeleteOnTermination bool `json:"delete_on_termination,omitempty"`
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

	// The iops of evs volume. Only required when volume_type is `GPSSD2` or `ESSD2`
	IOPS int `json:"iops,omitempty"`
	// The throughput of evs volume. Only required when volume_type is `GPSSD2`
	Throughput int `json:"throughput,omitempty"`

	ExtendParam *VolumeExtendParam `json:"extendparam,omitempty"`

	Metadata *VolumeMetadata `json:"metadata,omitempty"`

	ClusterId string `json:"cluster_id,omitempty"`
	// The cluster type is default to DSS
	ClusterType string `json:"cluster_type,omitempty"`
}

type DataVolume struct {
	VolumeType string `json:"volumetype" required:"true"`

	Size int `json:"size" required:"true"`

	MultiAttach *bool `json:"multiattach,omitempty"`

	PassThrough *bool `json:"hw:passthrough,omitempty"`

	// The iops of evs volume. Only required when volume_type is `GPSSD2` or `ESSD2`
	IOPS int `json:"iops,omitempty"`
	// The throughput of evs volume. Only required when volume_type is `GPSSD2`
	Throughput int `json:"throughput,omitempty"`

	Extendparam *VolumeExtendParam `json:"extendparam,omitempty"`

	Metadata *VolumeMetadata `json:"metadata,omitempty"`

	ClusterId string `json:"cluster_id,omitempty"`
	// The cluster type is default to DSS
	ClusterType string `json:"cluster_type,omitempty"`
}

type VolumeExtendParam struct {
	SnapshotId string `json:"snapshotId,omitempty"`
}

type VolumeMetadata struct {
	SystemEncrypted string `json:"__system__encrypted,omitempty"`
	SystemCmkid     string `json:"__system__cmkid,omitempty"`
}

type ServerExtendParam struct {
	ChargingMode        string `json:"chargingMode,omitempty"`
	RegionID            string `json:"regionID,omitempty"`
	PeriodType          string `json:"periodType,omitempty"`
	PeriodNum           int    `json:"periodNum,omitempty"`
	IsAutoRenew         string `json:"isAutoRenew,omitempty"`
	IsAutoPay           string `json:"isAutoPay,omitempty"`
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	SupportAutoRecovery string `json:"support_auto_recovery,omitempty"`

	// Specifies whether to support the function of creating a disk and then ECS: true of false
	DiskPrior string `json:"diskPrior,omitempty"`

	// When creating a spot ECS, set the parameter value to "spot"
	MarketType string `json:"marketType,omitempty"`
	// Specifies the highest price per hour you accept for a spot ECS
	SpotPrice string `json:"spotPrice,omitempty"`
	// Specifies the service duration of the spot ECS in hours
	SpotDurationHours int `json:"spot_duration_hours,omitempty"`
	// Specifies the number of time periods in the service duration
	SpotDurationCount int `json:"spot_duration_count,omitempty"`
	// Specifies the spot ECS interruption policy, which can only be set to "immediate" currently
	InterruptionPolicy string `json:"interruption_policy,omitempty"`
}

type MetaData struct {
	OpSvcUserId string `json:"op_svc_userid,omitempty"`
	AgencyName  string `json:"agency_name,omitempty"`
	AgentList   string `json:"__support_agent_list,omitempty"`
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

	// Indicates the filtering result for IPv4 addresses, which are accurate matched.
	// These IP addresses are private IP addresses of the ECS.
	IPEqual string `q:"ip_eq"`

	// Indicates the tags. The format is key=value, for example:
	// GET /v1/{project_id}/cloudservers/detail?tags=foo%3Dbar
	// In the preceding information, = needs to be escaped to %3D, foo indicates the tag key, and bar indicates the tag
	// value.
	Tags []string `q:"tags"`

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
		return ServerPage{pagination.PageSizeBase{PageResult: r}}
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

// ChangeAdminPassword alters the administrator or root password for a specified
// server.
func ChangeAdminPassword(client *golangsdk.ServiceClient, id, newPassword string) (r PasswordResult) {
	b := map[string]interface{}{
		"reset-password": map[string]string{
			"new_password": newPassword,
		},
	}
	_, r.Err = client.Put(passwordURL(client, id), b, nil, &golangsdk.RequestOpts{OkCodes: []int{204}})
	return
}

// UpdateOptsBuilder allows extensions to add additional attributes to the
// Update request.
type UpdateOptsBuilder interface {
	ToServerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts specifies the base attributes that may be updated on an existing
// server.
type UpdateOpts struct {
	Name        string  `json:"name,omitempty"`
	Hostname    string  `json:"hostname,omitempty"`
	UserData    []byte  `json:"-"`
	Description *string `json:"description,omitempty"`
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToServerUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var userData string
	if _, err := base64.StdEncoding.DecodeString(string(opts.UserData)); err != nil {
		userData = base64.StdEncoding.EncodeToString(opts.UserData)
	} else {
		userData = string(opts.UserData)
	}
	b["user_data"] = &userData

	return map[string]interface{}{"server": b}, nil
}

// Update requests that various attributes of the indicated server be changed.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToServerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// UpdateMetadata updates (or creates) all the metadata specified by opts for
// the given server ID. This operation does not affect already-existing metadata
// that is not specified by opts.
func UpdateMetadata(client *golangsdk.ServiceClient, id string, opts map[string]interface{}) (r UpdateMetadataResult) {
	b := map[string]interface{}{"metadata": opts}
	_, r.Err = client.Post(metadataURL(client, id), b, &r.Body, nil)
	return
}

// DeleteMetadatItem will delete the key-value pair with the given key for the given server ID.
func DeleteMetadatItem(client *golangsdk.ServiceClient, id, key string) (r DeleteMetadatItemResult) {
	_, r.Err = client.Delete(metadatItemURL(client, id, key), nil)
	return
}

// update auto terminate time for the given server ID.
func UpdateAutoTerminateTime(client *golangsdk.ServiceClient, id, terminateTime string) (r UpdateResult) {
	body := map[string]interface{}{
		"auto_terminate_time": terminateTime,
	}
	_, r.Err = client.Post(updateAutoTerminateTimeURL(client, id), body, nil, &golangsdk.RequestOpts{OkCodes: []int{204}})
	return
}
