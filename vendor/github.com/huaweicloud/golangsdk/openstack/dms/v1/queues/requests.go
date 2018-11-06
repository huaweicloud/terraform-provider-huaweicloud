package queues

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CreateOpsBuilder is used for creating queue parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToQueueCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	// Indicates the unique name of a queue.
	// A string of 1 to 64 characters that contain
	// a-z, A-Z, 0-9, hyphens (-), and underscores (_).
	// The name cannot be modified once specified.
	Name string `json:"name" required:"true"`

	// Indicates the queue type. Default value: NORMAL. Options:
	// NORMAL: Standard queue. Best-effort ordering.
	// FIFO: First-ln-First-out (FIFO) queue. FIFO delivery.
	// KAFKA_HA: High-availability Kafka queue.
	// KAFKA_HT: High-throughput Kafka queue.
	// AMQP: Advanced Message Queuing Protocol (AMQP) queue.
	QueueMode string `json:"queue_mode,omitempty"`

	// Indicates the basic information about a queue.
	// The queue description must be 0 to 160 characters in length,
	// and does not contain angle brackets (<) and (>).
	Description string `json:"description,omitempty"`

	// This parameter is mandatory only when queue_mode is NORMAL or FIFO.
	// Indicates whether to enable dead letter messages.
	// Default value: disable. Options: enable, disable.
	RedrivePolicy string `json:"redrive_policy,omitempty"`

	// This parameter is mandatory only when
	// redrive_policy is set to enable.
	// This parameter indicates the maximum number
	// of allowed message consumption failures.
	// Value range: 1-100.
	MaxConsumeCount int `json:"max_consume_count,omitempty"`

	// This parameter is mandatory only when
	// queue_mode is set to KAFKA_HA or KAFKA_HT.
	// This parameter indicates the retention time
	// of messages in Kafka queues.
	// Value range: 1 to 72 hours.
	RetentionHours int `json:"retention_hours,omitempty"`
}

// ToQueueCreateMap is used for type convert
func (ops CreateOps) ToQueueCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a queue with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToQueueCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})

	return
}

// Delete a queue by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get a queue with detailed information by id
func Get(client *golangsdk.ServiceClient, id string, includeDeadLetter bool) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id, includeDeadLetter), &r.Body, nil)
	return
}

// List all the queues
func List(client *golangsdk.ServiceClient, includeDeadLetter bool) pagination.Pager {
	url := listURL(client, includeDeadLetter)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return QueuePage{pagination.SinglePageBase(r)}
	})
}
