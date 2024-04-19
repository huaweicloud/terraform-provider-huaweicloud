package queues

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListQueuePropertyOpts is the structure that used to query the list of the queue properties.
type ListQueuePropertyOpts struct {
	// The queue name.
	QueueName string `json:"-" required:"true"`
	// The offset number.
	Offset int `q:"offset"`
	// Number of records to be queried.
	Limit int `q:"limit"`
}

// ListQueueProperty is the method that used to query list of the queue properties using given parameters.
func ListQueueProperty(c *golangsdk.ServiceClient, opts ListQueuePropertyOpts) ([]PropertyResp, error) {
	url := propertyURL(c, opts.QueueName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	url += query.String()
	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := QueuePropertyPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})

	pages, err := pager.AllPages()
	if err != nil {
		return nil, err
	}
	return ExtractProperties(pages)
}

// Property is the structure that used to update the property of the queue.
type Property struct {
	// Maximum number of Spark drivers can be started on this queue.
	MaxInstance int `json:"computeEngine.maxInstance,omitempty"`
	// Maximum number of tasks can be concurrently executed by a Spark driver.
	MaxConcurrent int `json:"job.maxConcurrent,omitempty"`
	// Maximum number of Spark drivers can be pre-started on this queue.
	MaxPrefetchInstance *int `json:"computeEngine.maxPrefetchInstance,omitempty"`
	// The cidr of the queue.
	Cidr string `json:"network.cidrInVpc,omitempty"`
}

// UpdateQueueProperty is the method that used to update property of the queue using given parameters.
func UpdateQueueProperty(c *golangsdk.ServiceClient, queueName string, opts Property) (*QueuePropertyResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "properties")
	if err != nil {
		return nil, err
	}
	var r QueuePropertyResp
	_, err = c.Put(propertyURL(c, queueName), b, &r, &golangsdk.RequestOpts{})
	return &r, err
}

// DeleteQueueProperties is the method that used to delete properties of the queue.
func DeleteQueueProperties(c *golangsdk.ServiceClient, queueName string, opts []string) (*QueuePropertyResp, error) {
	var r QueuePropertyResp
	_, err := c.DeleteWithResponse(propertyURL(c, queueName), &r, &golangsdk.RequestOpts{
		JSONBody: map[string]interface{}{
			"keys": opts,
		},
	})
	return &r, err
}
