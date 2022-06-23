package throttles

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ThrottlingPolicyOpts struct {
	// Maximum number of times an API can be accessed within a specified period.
	// The value of this parameter cannot exceed the default limit 200 TPS.
	// This value must be a positive integer and cannot exceed 2,147,483,647.
	ApiCallLimits int `json:"api_call_limits" required:"true"`
	// Request throttling policy name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits, and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// Period of time for limiting the number of API calls.
	// This parameter applies with each of the preceding three API call limits.
	// This value must be a positive integer and cannot exceed 2,147,483,647.
	TimeInterval int `json:"time_interval" required:"true"`
	// Time unit for limiting the number of API calls.
	// The valid values are as following:
	//     SECOND
	//     MINUTE
	//     HOUR
	//     DAY
	TimeUnit string `json:"time_unit" required:"true"`
	// Maximum number of times the API can be accessed by an app within the same period.
	// The value of this parameter must be less than that of user_call_limits.
	// This value must be a positive integer and cannot exceed 2,147,483,647.
	AppCallLimits int `json:"app_call_limits,omitempty"`
	// Description of the request throttling policy, which can contain a maximum of 255 characters.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description string `json:"remark,omitempty"`
	// Type of the request throttling policy.
	// 1: exclusive, limiting the maximum number of times a single API bound to the policy can be called within
	// the specified period.
	// 2: shared, limiting the maximum number of times all APIs bound to the policy can be called within the
	// specified period.
	Type int `json:"type,omitempty"`
	// Maximum number of times the API can be accessed by a user within the same period.
	// The value of this parameter must be less than that of api_call_limits.
	// This value must be a positive integer and cannot exceed 2,147,483,647.
	UserCallLimits int `json:"user_call_limits,omitempty"`
	// Maximum number of times the API can be accessed by an IP address within the same period.
	// The value of this parameter must be less than that of api_call_limits.
	// This value must be a positive integer and cannot exceed 2,147,483,647.
	IpCallLimits int `json:"ip_call_limits,omitempty"`
}

type ThrottlingPolicyOptsBuilder interface {
	ToThrottlingPolicyOptsMap() (map[string]interface{}, error)
}

func (opts ThrottlingPolicyOpts) ToThrottlingPolicyOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a new throttling policy.
func Create(client *golangsdk.ServiceClient, instanceId string, opts ThrottlingPolicyOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToThrottlingPolicyOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId), reqBody, &r.Body, nil)
	return
}

// Update is a method by which to udpate an existing throttle policy.
func Update(client *golangsdk.ServiceClient, instanceId, policyId string,
	opts ThrottlingPolicyOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToThrottlingPolicyOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, policyId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtain an existing APIG throttling policy by policy ID.
func Get(client *golangsdk.ServiceClient, instanceId, policyId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, policyId), &r.Body, nil)
	return
}

type ListOpts struct {
	// Request throttling policy ID.
	Id string `q:"id"`
	// Request throttling policy name.
	Name string `q:"name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// Parameter name (name) for exact matching.
	PreciseSearch string `q:"precise_search"`
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more throttling policies according to the query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ThorttlePage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing throttling policy.
func Delete(client *golangsdk.ServiceClient, instanceId, policyId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, policyId), nil)
	return
}

// SpecThrottleCreateOpts is a struct which will be used to create a new special throttling policy.
type SpecThrottleCreateOpts struct {
	// Maximum number of times the excluded object can access an API within the throttling period.
	CallLimits int `json:"call_limits" required:"true"`
	// Excluded app ID or excluded account ID.
	ObjectId string `json:"object_id" required:"true"`
	// Excluded object type, which supports APP and USER.
	ObjectType string `json:"object_type" required:"true"`
}

// SpecThrottleCreateOptsBuilder is an interface which to support request body build of
// the special throttling policy creation.
type SpecThrottleCreateOptsBuilder interface {
	ToSpecThrottleCreateOptsMap() (map[string]interface{}, error)
}

// ToSpecThrottleCreateOptsMap is a method which to build a request body by the SpecThrottleCreateOpts.
func (opts SpecThrottleCreateOpts) ToSpecThrottleCreateOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// CreateSpecThrottle is a method by which to create a new special throttling policy.
func CreateSpecThrottle(client *golangsdk.ServiceClient, instanceId, policyId string,
	opts SpecThrottleCreateOptsBuilder) (r SpecThrottleResult) {
	reqBody, err := opts.ToSpecThrottleCreateOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(specRootURL(client, instanceId, policyId), reqBody, &r.Body, nil)
	return
}

// SpecThrottleUpdateOpts is a struct which will be used to update an existing special throttling policy.
type SpecThrottleUpdateOpts struct {
	// Maximum number of times the excluded object can access an API within the throttling period.
	CallLimits int `json:"call_limits" required:"true"`
}

// SpecThrottleUpdateOptsBuilder is an interface which to support request body build of
// the special throttling policy updation.
type SpecThrottleUpdateOptsBuilder interface {
	ToSpecThrottleUpdateOptsMap() (map[string]interface{}, error)
}

// ToSpecThrottleUpdateOptsMap is a method which to build a request body by the SpecThrottleUpdateOpts.
func (opts SpecThrottleUpdateOpts) ToSpecThrottleUpdateOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// UpdateSpecThrottle is a method by which to update an existing special throttle policy.
func UpdateSpecThrottle(client *golangsdk.ServiceClient, instanceId, policyId, strategyId string,
	opts SpecThrottleUpdateOptsBuilder) (r SpecThrottleResult) {
	reqBody, err := opts.ToSpecThrottleUpdateOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(specResourceURL(client, instanceId, policyId, strategyId), reqBody, &r.Body,
		&golangsdk.RequestOpts{
			OkCodes: []int{200},
		})
	return
}

// SpecThrottlesListOpts allows to filter list data using given parameters.
type SpecThrottlesListOpts struct {
	// Object type, which can be APP or USER.
	ObjectType string `q:"object_type"`
	// Name of an excluded app.
	AppName string `q:"app_name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
}

// SpecThrottlesListOptsBuilder is an interface which to support request query build of
// the special throttling policies search.
type SpecThrottlesListOptsBuilder interface {
	ToSpecThrottlesListQuery() (string, error)
}

// ToSpecThrottlesListQuery is a method which to build a request query by the SpecThrottlesListOpts.
func (opts SpecThrottlesListOpts) ToSpecThrottlesListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListSpecThrottles is a method to obtain an array of one or more special throttling policies
// according to the query parameters.
func ListSpecThrottles(client *golangsdk.ServiceClient, instanceId, policyId string,
	opts SpecThrottlesListOptsBuilder) pagination.Pager {
	url := specRootURL(client, instanceId, policyId)
	if opts != nil {
		query, err := opts.ToSpecThrottlesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SpecThrottlePage{pagination.SinglePageBase(r)}
	})
}

// DeleteSpecThrottle is a method to delete an existing special throttling policy.
func DeleteSpecThrottle(client *golangsdk.ServiceClient, instanceId, policyId, strategyId string) (r DeleteResult) {
	_, r.Err = client.Delete(specResourceURL(client, instanceId, policyId, strategyId), nil)
	return
}

type BindOpts struct {
	// Throttle policy ID. The valid length is range from 1 to 65.
	ThrottleId string `json:"strategy_id" required:"true"`
	// The IDs of the API publish record.
	PublishIds []string `json:"publish_ids" required:"true"`
	// Instnace ID to which the API belongs.
	InstanceId string `json:"-" required:"true"`
}

// Bind is a method to bind the policy to one or more APIs.
func Bind(c *golangsdk.ServiceClient, opts BindOpts) ([]BindResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r struct {
		BindList []BindResp `json:"throttle_applys"`
	}
	_, err = c.Post(bindURL(c, opts.InstanceId), b, &r, nil)
	return r.BindList, err
}

// Unbind is a method to cancel the relationship between API and throttling policy.
func Unbind(c *golangsdk.ServiceClient, instanceId, bindId string) error {
	_, err := c.Delete(unbindURL(c, instanceId, bindId), nil)
	return err
}

type BatchUnbindActionOpts struct {
	// List of API and throttling policy binding relationship IDs that need to be unbound.
	Action string `q:"action"`
}

type BatchUnbindOpts struct {
	// List of API and throttling policy binding relationship IDs that need to be unbound.
	ThrottleBindings string `json:"throttle_bindings,omitempty"`
	// Instance ID.
	InstanceId string `json:"-" required:"true"`
}

// BatchUnbind is an API to unbind all throttling policies associated with the list.
func BatchUnbind(c *golangsdk.ServiceClient, unbindOpt BatchUnbindOpts,
	queryOpt BatchUnbindActionOpts) (*BatchResp, error) {
	b, err := golangsdk.BuildRequestBody(unbindOpt, "")
	if err != nil {
		return nil, err
	}

	url := batchUnbindURL(c, unbindOpt.InstanceId)
	query, err := golangsdk.BuildQueryString(queryOpt)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r BatchResp
	_, err = c.Put(url, b, &r, nil)
	return &r, err
}

type ListBindOpts struct {
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// throttling policy ID.
	ThrottleId string `q:"throttle_id"`
	// Environment ID.
	EnvId string `q:"env_id"`
	// API group ID.
	GroupId string `q:"group_id"`
	// API ID.
	ApiId string `q:"api_id"`
	// API name.
	ApiName string `q:"api_name"`
	// Instnace ID to which the API belongs.
	InstanceId string `json:"-"`
}

// ListBind is a method to obtain all API to which the throttling policy bound.
func ListBind(c *golangsdk.ServiceClient, opts ListBindOpts) ([]ApiForThrottle, error) {
	if opts.InstanceId == "" {
		return nil, fmt.Errorf("The instance ID is required.")
	}
	url := listBindURL(c, opts.InstanceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r struct {
		ApiDetails []ApiForThrottle `json:"apis"`
	}
	_, err = c.Get(url, &r, nil)
	return r.ApiDetails, err
}
