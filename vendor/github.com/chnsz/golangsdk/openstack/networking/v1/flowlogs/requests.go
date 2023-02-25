package flowlogs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlowLogsListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the subnet attributes you want to see returned.
type ListOpts struct {
	// Specifies the VPC flow log UUID.
	ID string `q:"id"`

	// Specifies the VPC flow log name.
	Name string `q:"name"`

	// Specifies the type of resource on which to create the VPC flow log..
	ResourceType string `q:"resource_type"`

	// Specifies the unique resource ID.
	ResourceID string `q:"resource_id"`

	// Specifies the type of traffic to log.
	TrafficType string `q:"traffic_type"`

	// Specifies the log group ID..
	LogGroupID string `q:"log_group_id"`

	// Specifies the log topic ID.
	LogTopicID string `q:"log_topic_id"`

	// Specifies the VPC flow log status, the value can be ACTIVE, DOWN or ERROR.
	Status string `q:"status"`

	//Specifies the number of records returned on each page.
	//The value ranges from 0 to intmax.
	Limit int `q:"limit"`

	//Specifies the resource ID of pagination query.
	//If the parameter is left blank, only resources on the first page are queried.
	Marker string `q:"marker"`
}

// ToFlowLogsListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlowLogsListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// VPC flow logs. It accepts a ListOpts struct, which allows you to filter
//  and sort the returned collection for greater efficiency.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToFlowLogsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FlowLogPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type CreateOpts struct {
	// Specifies the VPC flow log name. The value is a string of no more than 64
	// characters that can contain letters, digits, underscores (_), hyphens (-) and periods (.).
	Name string `json:"name,omitempty"`

	// Provides supplementary information about the VPC flow log.
	// The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).
	Description string `json:"description,omitempty"`

	// Specifies the type of resource on which to create the VPC flow log.
	// The value can be Port, VPC, and Network.
	ResourceType string `json:"resource_type" required:"true"`

	// Specifies the unique resource ID.
	ResourceID string `json:"resource_id" required:"true"`

	//Specifies the type of traffic to log. The value can be all, accept and reject.
	TrafficType string `json:"traffic_type" required:"true"`

	// Specifies the log group ID.
	LogGroupID string `json:"log_group_id" required:"true"`

	// Specifies the log topic ID.
	LogTopicID string `json:"log_topic_id" required:"true"`
}

type CreateOptsBuilder interface {
	ToCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "flow_log")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(CreateURL(client), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

func Delete(client *golangsdk.ServiceClient, flId string) (r DeleteResult) {
	url := DeleteURL(client, flId)
	_, r.Err = client.Delete(url, nil)
	return
}

func Get(client *golangsdk.ServiceClient, flId string) (r GetResult) {
	url := GetURL(client, flId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{})
	return
}

type UpdateOpts struct {
	// Specifies the VPC flow log name. The value is a string of no more than 64
	// characters that can contain letters, digits, underscores (_), hyphens (-) and periods (.).
	Name string `json:"name,omitempty"`

	// Provides supplementary information about the VPC flow log.
	// The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).
	Description string `json:"description,omitempty"`

	// Specifies whether to enable the VPC flow log function.
	AdminState bool `json:"admin_state"`
}

type UpdateOptsBuilder interface {
	ToUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "flow_log")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Update(client *golangsdk.ServiceClient, flId string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UpdateURL(client, flId), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
