package apis

import (
	"time"

	"github.com/huaweicloud/golangsdk"
)

type BackendInfo struct {
	Protocol  string `json:"req_protocol"`
	Method    string `json:"req_method"`
	Uri       string `json:"req_uri"`
	Timeout   int    `json:"timeout"`
	VpcStatus int    `json:"vpc_status"`
	VpcInfo   string `json:"vpc_info"`
	UrlDomain string `json:"url_domain"`
	Version   string `json:"version"`
	Remark    string `json:"remark"`
}

type FunctionInfo struct {
	FunctionUrn    string `json:"function_urn" required:"true"`
	InvocationType string `json:"invocation_type" required:"true"`
	Timeout        int    `json:"timeout" required:"true"`
	Version        string `json:"version,omitempty"`
	Remark         string `json:"remark,omitempty"`
}

type MockInfo struct {
	ResultContent string `json:"result_content,omitempty"`
	Version       string `json:"version,omitempty"`
	Remark        string `json:"remark,omitempty"`
}

type BackendParameter struct {
	Name     string `json:"name" required:"true"`
	Location string `json:"location" required:"true"`
	Origin   string `json:"origin" required:"true"`
	Value    string `json:"value" required:"true"`
	Remark   string `json:"remark,omitempty"`
}

type RequestParameter struct {
	Name         string `json:"name" required:"true"`
	Type         string `json:"type" required:"true"`
	Location     string `json:"location" required:"true"`
	Required     int    `json:"required" required:"true"`
	ValidEnable  int    `json:"valid_enable" required:"true"`
	DefaultValue string `json:"default_value,omitempty"`
	SampleValue  string `json:"sample_value,omitempty"`
	Remark       string `json:"remark,omitempty"`
	Enumerations string `json:"enumerations,omitempty"`
	MinNum       int    `json:"min_num,omitempty"`
	MaxNum       int    `json:"max_num,omitempty"`
	MinSize      int    `json:"min_size,omitempty"`
	MaxSize      int    `json:"max_size,omitempty"`
}

// ApiInstance contains all the information associated with a API.
type ApiInstance struct {
	// ID of the API
	Id string `json:"id"`
	// Name of the API
	Name string `json:"name"`
	// ID of the API group to which the API to be created will belong
	GroupId string `json:"group_id"`
	// Name of the API group
	GroupName string `json:"group_name"`
	// status of the API
	Status int `json:"status"`
	// Type of the API. 1: public, 2: private
	Type int `json:"type"`
	// Version of the API
	Version string `json:"version"`
	// Request protocol
	ReqProtocol string `json:"req_protocol"`
	// Request method
	ReqMethod string `json:"req_method"`
	// Access address
	ReqUri string `json:"req_uri"`
	// Request parameter list
	ReqParams []RequestParameter `json:"req_params"`
	// Security authentication mode, which can be: None, App and IAM
	AuthType string `json:"auth_type"`
	// Route matching mode
	MatchMode string `json:"match_mode"`
	// backend type, which can be: HTTP, Function and MOCK
	BackendType string `json:"backend_type"`
	// Backend parameter list
	BackendParams []BackendParameter `json:"backend_params"`
	//Web backend details
	BackendInfo BackendInfo `json:"backend_api"`
	// FunctionGraph backend details
	FunctionInfo FunctionInfo `json:"func_info"`
	// Mock backend details
	MockInfo MockInfo `json:"mock_info"`
	// Example response for a successful request
	ResultNormalSample string `json:"result_normal_sample"`
	// Example response for a failed request
	ResultFailureSample string `json:"result_failure_sample"`
	// Description of the API
	Remark string `json:"remark"`
	// tags of the API
	Tags []string `json:"tags"`
	// Description of the API request body
	BodyRemark string `json:"body_remark"`
	// whether CORS is supported
	Cors bool `json:"cors"`

	EnvName          string    `json:"run_env_name"`
	EnvId            string    `json:"run_env_id"`
	PublishId        string    `json:"publish_id"`
	ArrangeNecessary int       `json:"arrange_necessary"`
	RegisterAt       time.Time `json:"-"`
	UpdateAt         time.Time `json:"-"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract will get the Group object out of the commonResult object.
func (r commonResult) Extract() (*ApiInstance, error) {
	var s ApiInstance
	err := r.ExtractInto(&s)
	return &s, err
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}
