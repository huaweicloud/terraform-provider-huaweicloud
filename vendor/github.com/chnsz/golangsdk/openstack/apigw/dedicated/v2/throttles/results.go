package throttles

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// UpdateResult represents a result of the Update method.
type UpdateResult struct {
	commonResult
}

// GetResult represents a result of the Get method.
type GetResult struct {
	commonResult
}

type ThrottlingPolicy struct {
	// Number of APIs to which the request throttling policy has been bound.
	BindNum int `json:"bind_num"`
	// Indicates whether an excluded request throttling configuration has been created.
	// 1: yes
	// 2: no
	IsIncludeSpecialThrottle int `json:"is_include_special_throttle"`
	// Creation time.
	CreateTime string `json:"create_time"`
	// Description.
	Description string `json:"remark"`
	// Type of the request throttling policy.
	// 1: exclusive, limiting the maximum number of times a single API bound to the policy can be called within
	// the specified period.
	// 2: shared, limiting the maximum number of times all APIs bound to the policy can be called within the
	// specified period.
	Type int `json:"type"`
	// Period of time for limiting the number of API calls.
	TimeInterval int `json:"time_interval"`
	// Maximum number of times the API can be accessed by an IP address within the same period.
	IpCallLimits int `json:"ip_call_limits"`
	// Maximum number of times the API can be accessed by an app within the same period.
	AppCallLimits int `json:"app_call_limits"`
	// Request throttling policy name.
	Name string `json:"name"`
	// Time unit for limiting the number of API calls.
	// The valid values are as following:
	//     SECOND
	//     MINUTE
	//     HOUR
	//     DAY
	TimeUnit string `json:"time_unit"`
	// Maximum number of times an API can be accessed within a specified period.
	ApiCallLimits int `json:"api_call_limits"`
	// Request throttling policy ID.
	Id string `json:"id"`
	// Maximum number of times the API can be accessed by a user within the same period.
	UserCallLimits int `json:"user_call_limits"`
}

func (r commonResult) Extract() (*ThrottlingPolicy, error) {
	var s ThrottlingPolicy
	err := r.ExtractInto(&s)
	return &s, err
}

// The ThorttlePage represents the result of a List operation.
type ThorttlePage struct {
	pagination.SinglePageBase
}

// ExtractPolicies its Extract method to interpret it as a throttling policy array.
func ExtractPolicies(r pagination.Page) ([]ThrottlingPolicy, error) {
	var s []ThrottlingPolicy
	err := r.(ThorttlePage).Result.ExtractIntoSlicePtr(&s, "throttles")
	return s, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

// SpecThrottle is a struct that represents the result of CreateSpecThrottle, UpdateSpecThrottle and
// ListSpecThrottles methods.
type SpecThrottle struct {
	// Maximum number of times the excluded object can access an API within the throttling period.
	CallLimits int `json:"call_limits"`
	// Name of the app to which the excluded request throttling configuration applies.
	AppName string `json:"app_name"`
	// Name of an app or a tenant to which the excluded request throttling configuration applies.
	ObjectName string `json:"object_name"`
	// ID of an object specified in the excluded request throttling configuration.
	ObjectId string `json:"object_id"`
	// Request throttling policy ID.
	ThrottleId string `json:"throttle_id"`
	// Time when the excluded request throttling configuration is created.
	ApplyTime string `json:"apply_time"`
	// Excluded request throttling configuration ID.
	ID string `json:"id"`
	// ID of the app to which the excluded request throttling configuration applies.
	AppId string `json:"app_id"`
	// Excluded object type, which can be APP or USER.
	ObjectType string `json:"object_type"`
}

// The SpecThrottleResult represents the base result of the each special throttling polciy methods.
type SpecThrottleResult struct {
	commonResult
}

// The CreateSpecThrottleResult represents the result of the CreateSpecThrottle method.
type CreateSpecThrottleResult struct {
	SpecThrottleResult
}

// The UpdateSpecThrottleResult represents the result of the UpdateSpecThrottle method.
type UpdateSpecThrottleResult struct {
	SpecThrottleResult
}

// Extract is a method which to extract the response to a special throttling policy.
func (r SpecThrottleResult) Extract() (*SpecThrottle, error) {
	var s SpecThrottle
	err := r.ExtractInto(&s)
	return &s, err
}

// The SpecThrottlePage represents the result of a List operation.
type SpecThrottlePage struct {
	pagination.SinglePageBase
}

// ExtractSpecThrottles its Extract method to interpret it as a special throttling policy array.
func ExtractSpecThrottles(r pagination.Page) ([]SpecThrottle, error) {
	var s []SpecThrottle
	err := r.(SpecThrottlePage).Result.ExtractIntoSlicePtr(&s, "throttle_specials")
	return s, err
}

type BindResp struct {
	// The ID of the binding relationship.
	ID string `json:"id"`
	// The publish ID.
	PublishId string `json:"publish_id"`
	// Policy scope, the values are as follows:
	// 1: The entire API
	// 2: single user
	// 3: Single APP
	// Currently only 1 is supported.
	Scope int `json:"scope"`
	// Throttling policy ID
	StrategyId string `json:"strategy_id"`
	// Apply time.
	ApplyTime string `json:"apply_time"`
}

type ApiForThrottle struct {
	// API authentication method.
	AuthType string `json:"auth_type"`
	// Group name to which the API belongs.
	GroupName string `json:"group_name"`
	// API publish record ID.
	PublishId string `json:"publish_id"`
	// The number of the binding relationship with the throttling policy.
	ThrottleApplyId string `json:"throttle_apply_id"`
	// Binding time of the bound throttling policy.
	ApplyTime string `json:"apply_time"`
	// API description.
	Remark string `json:"remark"`
	// The ID of the environment published by the API.
	RunEnvId string `json:"run_env_id"`
	// API type
	Type int `json:"type"`
	// Binding throttling policy name.
	ThrottleName string `json:"throttle_name"`
	// Request URI of the API.
	ReqUri string `json:"req_uri"`
	// The name of the environment published by the API.
	RunEnvName string `json:"run_env_name"`
	// Group ID to which the API belogns.
	GroupId string `json:"group_id"`
	// API Name.
	Name string `json:"name"`
	// API ID.
	ID string `json:"id"`
}

type BatchResp struct {
	// Number of API and throttling policy bindings that have been successfully unbound.
	SuccessCount int `json:"success_count"`
	// Unbinding failed API and throttling policy binding relationship and error information.
	Failure []UnbindResult `json:"failure"`
}

type UnbindResult struct {
	// API and throttling policy binding relationship ID that failed to unbind.
	BindId string `json:"bind_id"`
	// Unbind failed error code
	ErrorCode string `json:"error_code"`
	// Unbind failed error message
	ErrorMsg string `json:"error_msg"`
	// ID of the API that failed to unbind
	ApiId string `json:"api_id"`
	// The name of the API that failed to unbind
	ApiName string `json:"api_name"`
}
