package elasticresourcepool

import (
	"github.com/chnsz/golangsdk/pagination"
)

// AssociateQueueResp is the structure that represents response of the AssociateElasticResourcePool or UpdateElasticResourcePoolQueuePolicy method.
type AssociateQueueResp struct {
	// Whether the request is successfully sent. Value true indicates that the request is successfully sent.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
}

// QueuePage is a page structure that represents each page information.
type QueuePage struct {
	pagination.OffsetPageBase
}

// QueuesResp is the structure that represents response of the ListElasticResourcePoolQueues method.
type QueuesResp struct {
	// Whether the request is successfully sent. Value true indicates that the request is successfully sent.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
	// The list of the queues associated with the specified elastic resource pool.
	Queues []Queue `json:"queues"`
	// The number of queues bound to the elastic resource pool.
	Count int `json:"count"`
}

// QueuesResp is the structure that represents the queue detail associated with the specified elastic resource pool.
type Queue struct {
	// The queue name.
	QueueName string `json:"queue_name"`
	// The enterprise project ID of the queue.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// The queue type.
	QueueType string `json:"queue_type"`
	// The list of scaling policies of the queue.
	QueueScalingPolicies []QueueScalingPolicy `json:"queue_scaling_policies"`
	// The owner of the queue.
	Owner string `json:"owner"`
	// Tht creation time of the queue.
	CreatedAt int `json:"create_time"`
	// The engine type ot the queue.
	Engine string `json:"engine"`
}

// ExtractQueues is a method to extract the list of queues associated with elastic resource pool.
func ExtractQueues(r pagination.Page) ([]Queue, error) {
	var s []Queue
	err := r.(QueuePage).Result.ExtractIntoSlicePtr(&s, "queues")
	return s, err
}
