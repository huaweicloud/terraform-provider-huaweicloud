package elasticresourcepool

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// AssociateQueueOpts is the structure that used to associate a queue with an elastic resource pool.
type AssociateQueueOpts struct {
	// The name of the elastic resource pool.
	ElasticResourcePoolName string `json:"-" required:"true"`
	// The name of queue.
	QueueName string `json:"queue_name" required:"true"`
}

// AssociateQueue is a method to associate a queue with an elastic resource pool using given parameters.
func AssociateQueue(c *golangsdk.ServiceClient, opts AssociateQueueOpts) (*AssociateQueueResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r AssociateQueueResp
	_, err = c.Post(associateQueueURl(c, opts.ElasticResourcePoolName), b, &r, nil)
	return &r, err
}

// UpdateQueuePolicyOpts is the structure that used to modify the scaling policy of the queue associated with an elastic resource pool.
type UpdateQueuePolicyOpts struct {
	// The name of the elastic resource pool.
	ElasticResourcePoolName string `json:"-" required:"true"`
	// The name of queue.
	QueueName string `json:"-" required:"true"`
	// The list of the queue scaling policies.
	QueueScalingPolicies []QueueScalingPolicy `json:"queue_scaling_policies" required:"true"`
}

// QueueScalingPolicy is the structure that represents the policy detail of specified queue associated with the specified elastic resource pool.
type QueueScalingPolicy struct {
	// The priority of the queue scaling policy. The valid value ranges from 1 to 100.
	// The larger the value, the higher the priority.
	Priority int `json:"priority" required:"true"`
	// The effective time of the queue scaling policy.
	ImpactStartTime string `json:"impact_start_time" required:"true"`
	// The expiration time of the queue scaling policy.
	ImpactStopTime string `json:"impact_stop_time" required:"true"`
	// The minimum number of CUs allowed by the scaling policy.
	MinCu int `json:"min_cu" required:"true"`
	// The maximum number of CUs allowed by the scaling policy.
	MaxCu int `json:"max_cu" required:"true"`
}

// UpdateElasticResourcePoolQueuePolicy is a method to modify the scaling policy of the queue associated with
// specified elastic resource pool using given parameters.
func UpdateElasticResourcePoolQueuePolicy(c *golangsdk.ServiceClient, opts UpdateQueuePolicyOpts) (*AssociateQueueResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r AssociateQueueResp
	_, err = c.Put(queueScalingPolicyURL(c, opts.ElasticResourcePoolName, opts.QueueName), b, &r, nil)
	return &r, err
}

// ListElasticResourcePoolQueuesOpts is the structure used to query the list of queue scaling policies in an elastic resource pool.
type ListElasticResourcePoolQueuesOpts struct {
	// The name of the elastic resource pool.
	ElasticResourcePoolName string `json:"-" required:"true"`
	// The name of the queue.
	QueueName string `q:"queue_name"`
	// Offset from which the query starts. Default to 0.
	Offset int `q:"offset"`
	// The number of items displayed on each page. Default to 100.
	Limit int `q:"limit"`
}

// ListElasticResourcePoolQueues is a method to query all queue scaling policies in an elastic resource pool using given parameters.
func ListElasticResourcePoolQueues(c *golangsdk.ServiceClient, opts ListElasticResourcePoolQueuesOpts) ([]Queue, error) {
	url := associateQueueURl(c, opts.ElasticResourcePoolName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := QueuePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}

	return ExtractQueues(pages)
}
