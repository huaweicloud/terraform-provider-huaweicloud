package apis

import (
	"github.com/huaweicloud/golangsdk"
)

type VpcChannel struct {
	VpcId        string `json:"vpc_id" required:"true"`
	VpcProxyHost string `json:"vpc_proxy_host,omitempty"`
}

type BackendOpts struct {
	Protocol     string     `json:"req_protocol" required:"true"`
	Method       string     `json:"req_method" required:"true"`
	Uri          string     `json:"req_uri" required:"true"`
	Timeout      int        `json:"timeout" required:"true"`
	VpcStatus    int        `json:"vpc_status,omitempty"`
	VpcInfo      VpcChannel `json:"vpc_info,omitempty"`
	UrlDomain    string     `json:"url_domain,omitempty"`
	Version      string     `json:"version,omitempty"`
	Remark       string     `json:"remark,omitempty"`
	AuthorizerId string     `json:"authorizer_id,omitempty"`
}
type FunctionOpts FunctionInfo
type MockOpts MockInfo

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAPICreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a API. This object is passed to
// the APIs Create function.
type CreateOpts struct {
	// ID of the API group to which the API to be created will belong
	// The value cannot be modified when updating
	GroupId string `json:"group_id" required:"true"`
	// Name of the API
	Name string `json:"name" required:"true"`
	// Type of the API. 1: public, 2: private
	Type int `json:"type" required:"true"`
	// Request method
	ReqMethod string `json:"req_method" required:"true"`
	// Access address
	ReqUri string `json:"req_uri" required:"true"`
	// Request protocol, defaults to HTTPS
	ReqProtocol string `json:"req_protocol,omitempty"`
	// Request parameter list
	ReqParams []RequestParameter `json:"req_params,omitempty"`
	// Security authentication mode, which can be: None, App and IAM
	AuthType string `json:"auth_type" required:"true"`
	// backend type, which can be: HTTP, Function and MOCK
	BackendType string `json:"backend_type" required:"true"`
	// Backend parameter list
	BackendParams []BackendParameter `json:"backend_params,omitempty"`
	//Web backend details
	BackendOpts BackendOpts `json:"backend_api,omitempty"`
	// FunctionGraph backend details
	FunctionOpts FunctionOpts `json:"func_info,omitempty"`
	// Mock backend details
	MockOpts MockOpts `json:"mock_info,omitempty"`
	// Example response for a successful request
	ResultNormalSample string `json:"result_normal_sample" required:"true"`
	// Example response for a failed request
	ResultFailureSample string `json:"result_failure_sample,omitempty"`
	// ID of customer authorizer
	AuthorizerId string `json:"authorizer_id,omitempty"`
	// Version of the API
	Version string `json:"version,omitempty"`
	// Route matching mode
	MatchMode string `json:"match_mode,omitempty"`
	// Description of the API
	Remark string `json:"remark,omitempty"`
	// tags of the API
	Tags []string `json:"tags,omitempty"`
	// Description of the API request body
	BodyRemark string `json:"body_remark,omitempty"`
	// whether CORS is supported
	Cors bool `json:"cors,omitempty"`
}

// ToAPICreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToAPICreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new API based on the values in CreateOpts.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAPICreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// Update will update the API with provided information. To extract the updated
// API from the response, call the Extract method on the UpdateResult.
// parameters of update is same with parameters of create. (group_id cannot be modified)
func Update(client *golangsdk.ServiceClient, id string, opts CreateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAPICreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(groupURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will delete the existing API with the provided ID.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(groupURL(client, id), nil)
	return
}

// Get retrieves the API with the provided ID. To extract the API object
// from the response, call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(groupURL(client, id), &r.Body, nil)
	return
}
