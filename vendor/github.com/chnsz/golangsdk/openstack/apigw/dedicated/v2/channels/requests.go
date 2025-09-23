package channels

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ChannelOpts is the structure that used to create a new channel.
type ChannelOpts struct {
	// The instance ID to which the channel belongs.
	InstanceId string `json:"-" requried:"true"`
	// Channel name.
	// A channel name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// Host port of the channel.
	// The valid value ranges from 1 to 65535.
	Port int `json:"port" required:"true"`
	// Distribution algorithm.
	// The valid values are as following:
	// + 1: WRR (default)
	// + 2: WLC
	// + 3: SH
	// + 4: URI hashing
	BalanceStrategy int `json:"balance_strategy" required:"true"`
	// Member type of the channel.
	// The valid values are as following:
	// + ip
	// + ecs (default)
	MemberType string `json:"member_type,omitempty"`
	// Channel type.
	// + 2: Server type.
	// + 3: Microservice type.
	Type int `json:"type,omitempty"`
	// builtin: server type
	// + microservice: microservice type
	// + reference: reference load balance channel
	// If vpc_channel_type is empty, the load balance channel type depends on the value of the type field.
	// If vpc_channel_type is non-empty and type is non-empty or non-zero, an error occurs when they are specified.
	// If vpc_channel_type is non-empty and type is empty or 0, the value of vpc_channel_type is used to specify the load balance channel type.
	VpcChannelType string `json:"vpc_channel_type,omitempty"`
	// Dictionary code of the channel.
	// The value can contain letters, digits, hyphens (-), underscores (_), and periods (.).
	DictCode string `json:"dict_code,omitempty"`
	// Backend server groups of the channel.
	// If omitted, you need to entry an empty array.
	MemberGroups []MemberGroup `json:"member_groups"`
	// Backend server list. Only one backend server is included if the channel type is set to 1.
	// If omitted, you need to entry an empty array.
	Members []MemberInfo `json:"members"`
	// Health check details.
	VpcHealthConfig *VpcHealthConfig `json:"vpc_health_config,omitempty"`
	// Microservice details.
	MicroserviceConfig *MicroserviceConfig `json:"microservice_info,omitempty"`
}

// MemberGroup is an object that represents the detail of the backend server group.
type MemberGroup struct {
	// Name of the backend server group of the channel.
	Name string `json:"member_group_name" required:"true"`
	// Description of the backend server group.
	Description string `json:"member_group_remark,omitempty"`
	// Weight of the backend server group.
	// If the server group contains servers and a weight has been set for it, the weight is automatically used to
	// assign weights to servers in this group.
	// The value is range from 0 to 100.
	Weight int `json:"member_group_weight,omitempty"`
	// Dictionary code of the backend server group.
	// The value can contain letters, digits, hyphens (-), underscores (_), and periods (.).
	DictCode string `json:"dict_code,omitempty"`
	// Version of the backend server group.
	// This parameter is supported only when the channel type is microservice.
	MicroserviceVersion string `json:"microservice_version,omitempty"`
	// Port of the backend server group.
	// This parameter is supported only when the channel type is microservice.
	// If the port number is 0, all addresses in the backend server group use the original load balancing port to
	// inherit logic.
	// The value is range from 0 to 65535.
	MicroservicePort int `json:"microservice_port,omitempty"`
	// Tags of the backend server group.
	// This parameter is supported only when the channel type is microservice.
	MicroserviceLabels []MicroserviceLabel `json:"microservice_labels,omitempty"`
	// ID of the reference load balance channel.
	// This parameter is supported only when the VPC channel type is reference (vpc_channel_type=reference).
	ReferenceVpcChannelId string `json:"reference_vpc_channel_id,omitempty"`
}

// MicroserviceLabel is an object that represents a specified microservice label.
type MicroserviceLabel struct {
	// Label name.
	Name string `json:"label_name" required:"true"`
	// Label value.
	Value string `json:"label_value" required:"true"`
}

// MemberInfo is an object that represents the backend member detail.
type MemberInfo struct {
	// Backend server ID.
	// This parameter is valid when the member type is instance.
	// The value can contain 1 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	EcsId string `json:"ecs_id,omitempty"`
	// Backend server name, which contains of 1 to 64 characters, including letters, digits, periods (.), hyphens (-)
	// and underscores (_).
	// This parameter is valid when the member type is instance.
	EcsName string `json:"ecs_name,omitempty"`
	// Backend server address.
	// This parameter is valid when the member type is IP address.
	Host string `json:"host,omitempty"`
	// Backend server weight.
	// The higher the weight is, the more requests a cloud server will receive.
	// The weight is only available for the WRR and WLC algorithms.
	// It is valid only when the channel type is set to 2.
	// The valid value ranges from 0 to 100.
	Weight *int `json:"weight,omitempty"`
	// Whether the backend service is a standby node.
	// After you enable this function, the backend service serves as a standby node.
	// It works only when all non-standby nodes are faulty.
	// This function is supported only when your gateway has been upgraded to the corresponding version.
	// If your gateway does not support this function, contact technical support.
	// Defaults to false.
	IsBackup *bool `json:"is_backup,omitempty"`
	// Backend server group name. The server group facilitates backend service address modification.
	GroupName string `json:"member_group_name,omitempty"`
	// Backend server status.
	// + 1: available
	// + 2: unavailable
	Status int `json:"status,omitempty"`
	// Backend server port.
	// The valid value ranges from 0 to 65535.
	Port *int `json:"port,omitempty"`
}

// VpcHealthConfig is an object that represents the health check configuration.
type VpcHealthConfig struct {
	// Protocol for performing health checks on backend servers in the channel.
	// The valid values are as following:
	// + TCP
	// + HTTP
	// + HTTPS
	Protocol string `json:"protocol" required:"true"`
	// Healthy threshold, which refers to the number of consecutive successful checks required for a backend server to
	// be considered healthy.
	// The valid value ranges from 1 to 10.
	ThresholdNormal int `json:"threshold_normal" required:"true"`
	// Unhealthy threshold, which refers to the number of consecutive failed checks required for a backend server to be
	// considered unhealthy.
	// The valid value range from 1 to 10.
	ThresholdAbnormal int `json:"threshold_abnormal" required:"true"`
	// Interval between consecutive checks, in second. The value must be greater than the value of timeout.
	// The valid value ranges from 1 to 300.
	TimeInterval int `json:"time_interval" required:"true"`
	// Timeout for determining whether a health check fails, in second.
	// The value must be less than the value of time_interval.
	// The valid value ranges from 1 to 30.
	Timeout int `json:"timeout" required:"true"`
	// Destination path for health checks. This parameter is required if protocol is set to http.
	Path string `json:"path,omitempty"`
	// Request method for health checks.
	// The valid values are as following:
	// + GET (default)
	// + HEAD
	Method string `json:"method,omitempty"`
	// Destination port for health checks. By default, the host port of the channel is used.
	// The valid value ranges from 1 to 65535.
	Port int `json:"port,omitempty"`
	// Response codes for determining a successful HTTP response.
	// The value can be any integer within 100–599 in one of the following formats:
	// + Value, for example, 200.
	// + Multiple values, for example, 200,201,202.
	// + Range, for example, 200-299.
	// + Multiple values and ranges, for example, 201,202,210-299.
	// This parameter is required if protocol is set to http.
	HttpCodes string `json:"http_code,omitempty"`
	// Indicates whether to enable two-way authentication.
	// If this function is enabled, the certificate specified in the backend_client_certificate configuration item of
	// the gateway is used. Defaults to false.
	EnableClientSsl *bool `json:"enable_client_ssl,omitempty"`
	// Health check result.
	// + 1: available
	// + 2: unavailable
	Status int `json:"status,omitempty"`
}

// MicroserviceConfig is an object that represents the microservice configuration.
type MicroserviceConfig struct {
	// Microservice type.
	// + CSE: CSE microservice registration center.
	// + CCE: Cloud Container Engine (CCE).
	ServiceType string `json:"service_type,omitempty"`
	// CSE microservice details. This parameter is required if service_type is set to CSE.
	CseInfo *MicroserviceCseInfo `json:"cse_info,omitempty"`
	// CCE workload details. This parameter is required if service_type is set to CCE.
	CceInfo *MicroserviceCceInfo `json:"cce_info,omitempty"`
}

// MicroserviceCseInfo is an object that represents the CSE microservice detail.
type MicroserviceCseInfo struct {
	// Microservice engine ID.
	EngineId string `json:"engine_id" required:"true"`
	// Microservice ID.
	ServiceId string `json:"service_id,omitempty"`
}

// MicroserviceCceInfo is an object that represents the CCE microservice detail.
type MicroserviceCceInfo struct {
	// CCE cluster ID.
	ClusterId string `json:"cluster_id,omitempty"`
	// CCE namespace.
	Namespace string `json:"namespace,omitempty"`
	// Workload type.
	// + deployment
	// + statefulset
	// + daemonset
	WorkloadType string `json:"workload_type,omitempty"`
	// Application name.
	AppName string `json:"app_name,omitempty"`
	// Service label key. Start with a letter or digit, and use only letters, digits, and these special
	// characters: -_./:(). (1 to 64 characters)
	LabelKey string `json:"label_key,omitempty"`
	// Service label value. Start with a letter, and include only letters, digits, periods (.), hyphens (-),
	// and underscores (_). (1 to 64 characters)
	LabelValue string `json:"label_value,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new channel using given parameters.
func Create(client *golangsdk.ServiceClient, opts ChannelOpts) (*Channel, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Channel
	_, err = client.Post(rootURL(client, opts.InstanceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain an existing channel by its ID and related instance ID.
func Get(client *golangsdk.ServiceClient, instanceId, chanId string) (*Channel, error) {
	var r Channel
	_, err := client.Get(resourceURL(client, instanceId, chanId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// The instance ID to which the channel belongs.
	InstanceId string `json:"-" requried:"true"`
	// Channel ID.
	ID string `q:"id"`
	// Channel name.
	Name string `q:"name"`
	// Dictionary code of the backend server group.
	// The value can contain letters, digits, hyphens (-), underscores (_), and periods (.).
	DictCode string `q:"dict_code"`
	// Backend service address. By default, exact match is used. Fuzzy match is not supported.
	MemberHost string `q:"member_host"`
	// Backend server port. The valid value ranges from 0 to 65535.
	MemberPort string `q:"member_port"`
	// The name of the backend server group.
	MemberGroupName string `q:"member_group_name"`
	// The ID of the backend server group.
	MemberGroupId string `q:"member_group_id"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Defaults to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	Limit int `q:"limit"`
	// Parameter name for exact matching(, only parameter 'name' and 'member_group_name' are support yet).
	PreciseSearch string `q:"precise_search"`
}

// List is a method to obtain an array of one or more channels by query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOpts) ([]Channel, error) {
	url := rootURL(client, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ChannelPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractChannels(pages)
}

// Update is a method by which to update an existing channel by request parameters.
func Update(client *golangsdk.ServiceClient, chanId string, opts ChannelOpts) (*Channel, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Channel
	_, err = client.Put(resourceURL(client, opts.InstanceId, chanId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing channel using its ID and related instance ID.
func Delete(client *golangsdk.ServiceClient, instanceId, chanId string) error {
	_, err := client.Delete(resourceURL(client, instanceId, chanId), nil)
	return err
}

type BackendAddOpts struct {
	// The instance ID to which the channel belongs.
	InstanceId string `json:"-" requried:"true"`
	// Backend server list.
	Members []MemberInfo `json:"members" required:"true"`
}

// AddBackendServices is a method to add a backend instance to channel.
func AddBackendServices(client *golangsdk.ServiceClient, chanId string, opts BackendAddOpts) ([]MemberInfo, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r struct {
		Members []MemberInfo `json:"members"`
	}
	_, err = client.Post(membersURL(client, opts.InstanceId, chanId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Members, err
}

// BackendListOpts allows to filter list data using given parameters.
type BackendListOpts struct {
	// The instance ID to which the channel belongs.
	InstanceId string `json:"-" requried:"true"`
	// Cloud server name.
	Name string `q:"name"`
	// The name of the backend server group.
	MemberGroupName string `q:"member_group_name"`
	// The ID of the backend server group.
	MemberGroupId string `q:"member_group_id"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	Limit int `q:"limit"`
	// Parameter name for exact matching(, only parameter 'name' and 'member_group_name' are support yet).
	PreciseSearch string `q:"precise_search"`
}

// ListBackendServices is a method to obtain an array of one or more backend services by query parameters.
func ListBackendServices(client *golangsdk.ServiceClient, chanId string, opts BackendListOpts) ([]MemberInfo, error) {
	url := membersURL(client, opts.InstanceId, chanId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := MemberPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractMembers(pages)
}

// RemoveBackendService is a method to remove an existing backend instance form channel.
func RemoveBackendService(client *golangsdk.ServiceClient, instanceId, chanId, memberId string) error {
	_, err := client.Delete(memberURL(client, instanceId, chanId, memberId), nil)
	return err
}
