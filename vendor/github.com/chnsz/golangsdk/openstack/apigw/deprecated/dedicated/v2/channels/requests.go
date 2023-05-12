package channels

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ChannelOpts struct {
	// Backend server list. Only one backend server is included if the VPC channel type is set to 1.
	Members []MemberInfo `json:"members" required:"true"`
	// VPC channel name.
	// A VPC channel name can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// VPC channel type.
	//     1: private network ELB channel (to be deprecated)
	//     2: fast channel with the load balancing function
	Type int `json:"type" required:"true"`
	// Health check details.
	VpcHealthConfig VpcHealthConfig `json:"vpc_health_config" required:"true"`
	// Distribution algorithm.
	// The valid values are as following:
	//     1: WRR (default)
	//     2: WLC
	//     3: SH
	//     4: URI hashing
	// This parameter is mandatory if the VPC channel type is set to 2.
	BalanceStrategy int `json:"balance_strategy,omitempty"`
	// Member type of the VPC channel.
	// The valid values are as following:
	//     ip
	//     ecs (default)
	// This parameter is required if the VPC channel type is set to 2.
	MemberType string `json:"member_type,omitempty"`
	// Host port of the VPC channel.
	// This parameter is valid only when the VPC channel type is set to 2. The value range is 1–65535.
	// This parameter is required if the VPC channel type is set to 2.
	Port int `json:"port,omitempty"`
}

type MemberInfo struct {
	// Backend server ID.
	// This parameter is valid when the member type is instance.
	// The value can contain 1 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	EcsId string `json:"ecs_id,omitempty" xor:"Host"`
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
	// It is valid only when the VPC channel type is set to 2.
	// The valid value is range from 0 to 100.
	Weight int `json:"weight,omitempty"`
}

type VpcHealthConfig struct {
	// Protocol for performing health checks on backend servers in the VPC channel.
	// The valid values are as following:
	//     TCP
	//     HTTP
	//     HTTPS
	Protocol string `json:"protocol" required:"true"`
	// Healthy threshold, which refers to the number of consecutive successful checks required for a backend server to
	// be considered healthy.
	// The valid value is range from 2 to 10.
	ThresholdNormal int `json:"threshold_normal" required:"true"`
	// Unhealthy threshold, which refers to the number of consecutive failed checks required for a backend server to be
	// considered unhealthy.
	// The valid value is range from 2 to 10.
	ThresholdAbnormal int `json:"threshold_abnormal" required:"true"`
	// Interval between consecutive checks, in second. The value must be greater than the value of timeout.
	// The valid value is range from 5 to 300.
	TimeInterval int `json:"time_interval" required:"true"`
	// Timeout for determining whether a health check fails, in second.
	// The value must be less than the value of time_interval.
	// The valid value is range from 2 to 30.
	Timeout int `json:"timeout" required:"true"`
	// Indicates whether to enable two-way authentication.
	// If this function is enabled, the certificate specified in the backend_client_certificate configuration item of
	// the gateway is used. Default to false.
	EnableClientSsl bool `json:"enable_client_ssl,omitempty"`
	// Response codes for determining a successful HTTP response.
	// The value can be any integer within 100–599 in one of the following formats:
	//     Value, for example, 200.
	//     Multiple values, for example, 200,201,202.
	//     Range, for example, 200-299.
	//     Multiple values and ranges, for example, 201,202,210-299.
	// This parameter is required if protocol is set to http.
	HttpCodes string `json:"http_code,omitempty"`
	// Request method for health checks.
	// The valid values are as following:
	//     GET (default)
	//     HEAD
	Method string `json:"method,omitempty"`
	// Destination path for health checks. This parameter is required if protocol is set to http.
	Path string `json:"path,omitempty"`
	// Destination port for health checks. By default, the host port of the VPC channel is used.
	// The valid value is range from 1 to 65535.
	Port int `json:"port,omitempty"`
}

type ChannelOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

func (opts ChannelOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create an APIG vpc channel.
func Create(client *golangsdk.ServiceClient, instanceId string, opts ChannelOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Update is a method by which to update an existing APIG vpc channel by request parameters.
func Update(client *golangsdk.ServiceClient, instanceId, chanId string, opts ChannelOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, chanId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtain an existing APIG vpc channel by channel ID.
func Get(client *golangsdk.ServiceClient, instanceId, chanId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, chanId), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// VPC channel ID.
	Id string `q:"id"`
	// VPC channel name.
	Name string `q:"name"`
	// VPC channel type.
	VpcType int `q:"vpc_type"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	Limit int `q:"limit"`
	// Parameter name (name) for exact matching.
	PreciseSearch string `q:"precise_search"`
}

type ListOptsBuilder interface {
	ToChannelListQuery() (string, error)
}

func (opts ListOpts) ToChannelListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more vpc channels by query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToChannelListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ChannelPage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing vpc channel.
func Delete(client *golangsdk.ServiceClient, instanceId, chanId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, chanId), nil)
	return
}

type BackendAddOpts struct {
	// Backend server list.
	Members []MemberInfo `json:"members" required:"true"`
}

type BackendAddOptsBuilder interface {
	ToBackendAddMap() (map[string]interface{}, error)
}

func (opts BackendAddOpts) ToBackendAddMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// AddBackendServices is a method to add a backend instance to vpc channel.
func AddBackendServices(client *golangsdk.ServiceClient, instanceId, chanId string,
	opts BackendAddOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToBackendAddMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(membersURL(client, instanceId, chanId), reqBody, &r.Body, nil)
	return
}

// BackendListOpts allows to filter list data using given parameters.
type BackendListOpts struct {
	// Cloud server name.
	Name string `q:"name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	Limit int `q:"limit"`
}

type BackendListOptsBuilder interface {
	ToBackendListQuery() (string, error)
}

func (opts BackendListOpts) ToBackendListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListBackendServices is a method to obtain an array of one or more backend services by query parameters.
func ListBackendServices(client *golangsdk.ServiceClient, instanceId, chanId string,
	opts BackendListOptsBuilder) pagination.Pager {
	url := membersURL(client, instanceId, chanId)
	if opts != nil {
		query, err := opts.ToBackendListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ChannelPage{pagination.SinglePageBase(r)}
	})
}

// RemoveBackendService is a method to remove an existing backend instance form vpc channel.
func RemoveBackendService(client *golangsdk.ServiceClient, instanceId, chanId, memberId string) (r RemoveResult) {
	_, r.Err = client.Delete(memberURL(client, instanceId, chanId, memberId), nil)
	return
}
