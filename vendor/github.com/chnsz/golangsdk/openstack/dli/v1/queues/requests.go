package queues

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateOpts contains the options for create a service. This object is passed to Create().
type CreateOpts struct {
	// Name of a newly created resource queue. The name can contain only digits, letters, and underscores (_),
	//but cannot contain only digits or start with an underscore (_). Length range: 1 to 128 characters.
	QueueName string `json:"queue_name" required:"true"`

	// Indicates the queue type. The options are as follows:
	// sql
	// general
	// all
	// NOTE:
	// If the type is not specified, the default value sql is used.
	QueueType string `json:"queue_type,omitempty"`

	// Description of a queue.
	Description string `json:"description,omitempty"`

	// Minimum number of CUs that are bound to a queue. Currently, the value can only be 16, 64, or 256.
	CuCount int `json:"cu_count" required:"true"`

	// Billing mode of a queue. This value can only be set to 1, indicating that the billing is based on the CUH used.
	ChargingMode int `json:"charging_mode,omitempty"`

	// Enterprise project ID. The value 0 indicates the default enterprise project.
	// NOTE:
	// Users who have enabled Enterprise Management can set this parameter to bind a specified project.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`

	// CPU architecture of queue computing resources.
	// x86_64 (default)
	// aarch64
	Platform string `json:"platform,omitempty"`

	// Queue resource mode. The options are as follows:
	// 0: indicates the shared resource mode.
	// 1: indicates the exclusive resource mode.
	ResourceMode int `json:"resource_mode"`

	// Specifies the tag information of the queue to be created,
	// including the JSON character string indicating whether the queue is Dual-AZ. Currently,
	// only the value 2 is supported, indicating that two queues are created.
	Labels []string `json:"labels,omitempty"`

	// Indicates the queue feature. The options are as follows:
	// basic: basic type
	// ai: AI-enhanced (Only the SQL x86_64 dedicated queue supports this option.)
	// The default value is basic.
	// NOTE:
	// For an enhanced AI queue, an AI image is loaded in the background.
	// The image integrates AI algorithm packages based on the basic image.
	Feature string `json:"feature,omitempty"`

	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

type ListOpts struct {
	QueueType      string `q:"queue_type"`
	WithPriv       bool   `q:"with-priv"`
	WithChargeInfo bool   `q:"with-charge-info"`
	Tags           string `q:"tags"`
}

type ActionOpts struct {
	Action    string `json:"action" required:"true"` //Operations to be performed: restart; scale_out;scale_in
	Force     bool   `json:"force,omitempty"`        //when action= restart,can Specifies whether to forcibly restart
	CuCount   int    `json:"cu_count,omitempty"`     // Number of CUs to be scaled in or out.
	QueueName string `json:"-" required:"true"`
}

type UpdateCidrOpts struct {
	Cidr string `json:"cidr_in_vpc,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToDomainCreateMap() (map[string]interface{}, error)
}

// ToDomainCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToDomainCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

/*
This API is used to create a queue.

@cloudAPI-URL: POST /v1.0/{project_id}/queues
@cloudAPI-ResourceType: dli
@cloudAPI-version: v1.0

@since: 2021-07-07 12:12:12

*/
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	requstbody, err := opts.ToDomainCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(createURL(c), requstbody, &r.Body, reqOpt)
	return
}

/*
This API is used to delete a queue.

@cloudAPI-URL: Delete /v1.0/{project_id}/queues/{queueName}
@cloudAPI-ResourceType: dli
@cloudAPI-version: v1.0

@since: 2021-07-07 12:12:12

*/
func Delete(c *golangsdk.ServiceClient, queueName string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Delete(resourceURL(c, queueName), reqOpt)
	return
}

/*
This API is used to query all Queue  list

@cloudAPI-URL: GET /v1.0/{project_id}/queues
@cloudAPI-ResourceType: dli
@cloudAPI-version: v1.0

@since: 2021-07-07 12:12:12

*/
func List(c *golangsdk.ServiceClient, listOpts ListOptsBuilder) (r GetResult) {
	listResult := new(ListResult)

	url := queryAllURL(c)
	if listOpts != nil {
		query, err := listOpts.ToListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Get(queryAllURL(c), &listResult, reqOpt)
	r.Body = listResult
	return r
}

/*
This API is used to query the Details of a Queue

@cloudAPI-URL: GET /v1.0/{project_id}/queues/{queueName}
@cloudAPI-ResourceType: dli
@cloudAPI-version: v1.0

@since: 2021-07-07 12:12:12

*/
func Get(c *golangsdk.ServiceClient, queueName string) (r GetResult) {
	result := new(Queue4Get)

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Get(resourceURL(c, queueName), &result, reqOpt)

	if result != nil {
		r.Body = result
	}

	return r
}

func ScaleOrRestart(c *golangsdk.ServiceClient, opts ActionOpts) (r PutResult) {
	requstbody, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(ActionURL(c, opts.QueueName), requstbody, &r.Body, reqOpt)
	return
}

func UpdateCidr(c *golangsdk.ServiceClient, queueName string, opts UpdateCidrOpts) (*UpdateCidrResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst UpdateCidrResp
	_, err = c.Put(resourceURL(c, queueName), b, &rst, nil)
	return &rst, err
}
