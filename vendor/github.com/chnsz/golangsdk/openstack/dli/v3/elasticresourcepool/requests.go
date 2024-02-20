package elasticresourcepool

import (
	"github.com/chnsz/golangsdk"
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
