package responses

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ResponseOpts allows to create a new custom response or update the existing custom response using given parameters.
type ResponseOpts struct {
	// APIG group name, which can contain 1 to 64 characters, only letters, digits, hyphens (-) and
	// underscores (_) are allowed.
	Name string `json:"name" required:"true"`
	// Response type definition, which includes a key and value. Options of the key:
	//     AUTH_FAILURE: Authentication failed.
	//     AUTH_HEADER_MISSING: The identity source is missing.
	//     AUTHORIZER_FAILURE: Custom authentication failed.
	//     AUTHORIZER_CONF_FAILURE: There has been a custom authorizer error.
	//     AUTHORIZER_IDENTITIES_FAILURE: The identity source of the custom authorizer is invalid.
	//     BACKEND_UNAVAILABLE: The backend service is unavailable.
	//     BACKEND_TIMEOUT: Communication with the backend service timed out.
	//     THROTTLED: The request was rejected due to request throttling.
	//     UNAUTHORIZED: The app you are using has not been authorized to call the API.
	//     ACCESS_DENIED: Access denied.
	//     NOT_FOUND: No API is found.
	//     REQUEST_PARAMETERS_FAILURE: The request parameters are incorrect.
	//     DEFAULT_4XX: Another 4XX error occurred.
	//     DEFAULT_5XX: Another 5XX error occurred.
	// Each error type is in JSON format.
	Responses map[string]ResponseInfo `json:"responses,omitempty"`
	// APIG dedicated instance ID.
	InstanceId string `json:"-"`
	// APIG group ID.
	GroupId string `json:"-"`
}

type ResponseInfo struct {
	// Response body template.
	Body string `json:"body" required:"true"`
	// HTTP status code of the response. If omitted, the status will be cancelled.
	Status int `json:"status,omitempty"`
	// The configuration of the custom response headers.
	Headers []ResponseInfoHeader `json:"headers,omitempty"`
	// Indicates whether the response is the default response.
	// Only the response of the API call is supported.
	IsDefault bool `json:"default,omitempty"`
}

type ResponseInfoHeader struct {
	// The key name of the response header.
	// The valid length is limited from 1 to 128, only English letters, digits and hyphens (-) are allowed.
	Key string `json:"key,omitempty"`
	// The value for the specified response header key.
	// The valid length is limited from 1 to 1,024.
	Value string `json:"value,omitempty"`
}

type ResponseOptsBuilder interface {
	ToResponseOptsMap() (map[string]interface{}, error)
	GetInstanceId() string
	GetGroupId() string
}

func (opts ResponseOpts) ToResponseOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts ResponseOpts) GetInstanceId() string {
	return opts.InstanceId
}

func (opts ResponseOpts) GetGroupId() string {
	return opts.GroupId
}

// Create is a method by which to create function that create a new custom response.
func Create(client *golangsdk.ServiceClient, opts ResponseOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToResponseOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, opts.GetInstanceId(), opts.GetGroupId()),
		reqBody, &r.Body, nil)
	return
}

// Update is a method by which to create function that udpate the existing custom response.
func Update(client *golangsdk.ServiceClient, respId string, opts ResponseOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToResponseOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, opts.GetInstanceId(), opts.GetGroupId(), respId),
		reqBody, &r.Body, &golangsdk.RequestOpts{
			OkCodes: []int{200},
		})
	return
}

// Get is a method to obtain the specified custom response according to the instanceId, appId and respId.
func Get(client *golangsdk.ServiceClient, instanceId, groupId, respId string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, groupId, respId), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// APIG dedicated instance ID.
	InstanceId string `json:"-"`
	// APIG group ID.
	GroupId string `json:"-"`
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
	GetInstanceId() string
	GetGroupId() string
}

func (opts ListOpts) GetInstanceId() string {
	return opts.InstanceId
}

func (opts ListOpts) GetGroupId() string {
	return opts.GroupId
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more custom reponses according to the query parameters.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, opts.GetInstanceId(), opts.GetGroupId())
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResponsePage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete the existing custom response.
func Delete(client *golangsdk.ServiceClient, instanceId, groupId, respId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, groupId, respId), nil)
	return
}

// SpecRespOpts is used to build the APIG response url. All parameters are required.
type SpecRespOpts struct {
	InstanceId string
	GroupId    string
	RespId     string
}

type SpecRespOptsBuilder interface {
	GetInstanceId() string
	GetGroupId() string
	GetResponseId() string
}

func (opts SpecRespOpts) GetInstanceId() string {
	return opts.InstanceId
}

func (opts SpecRespOpts) GetGroupId() string {
	return opts.GroupId
}

func (opts SpecRespOpts) GetResponseId() string {
	return opts.RespId
}

// GetSpecResp is a method to get the specifies custom response configuration from an group.
func GetSpecResp(client *golangsdk.ServiceClient, respType string, opts SpecRespOptsBuilder) (r GetSpecRespResult) {
	url := specResponsesURL(client, opts.GetInstanceId(), opts.GetGroupId(), opts.GetResponseId(), respType)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

type UpdateSpecRespBuilder interface {
	ToSpecRespUpdateMap() (map[string]interface{}, error)
}

func (opts ResponseInfo) ToSpecRespUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// UpdateSpecResp is a method to update an existing custom response configuration from an group.
func UpdateSpecResp(client *golangsdk.ServiceClient, respType string, specOpts SpecRespOptsBuilder,
	respOpts UpdateSpecRespBuilder) (r UpdateSpecRespResult) {
	reqBody, err := respOpts.ToSpecRespUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	url := specResponsesURL(client, specOpts.GetInstanceId(), specOpts.GetGroupId(), specOpts.GetResponseId(), respType)
	_, r.Err = client.Put(url, reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// DeleteSpecResp is a method to delete an existing custom response configuration from an group.
func DeleteSpecResp(client *golangsdk.ServiceClient, respType string, specOpts SpecRespOptsBuilder) (r DeleteResult) {
	url := specResponsesURL(client, specOpts.GetInstanceId(), specOpts.GetGroupId(), specOpts.GetResponseId(), respType)
	_, r.Err = client.Delete(url, nil)
	return
}
