package queues

import "github.com/chnsz/golangsdk/pagination"

// PropertyResp is the structure that represents property detail of the queue.
type PropertyResp struct {
	// The valid values are as follows:
	// + computeEngine.maxInstances: Maximum number of Spark drivers can be started on this queue.
	// + computeEngine.maxPrefetchInstan: Maximum number of Spark drivers can be pre-started on this queue.
	// + job.maxConcurrent: Maximum number of tasks can be concurrently executed by a Spark driver.
	// + multipleSc.support: Whether multiple spark drivers can be configured.
	Key   string `json:"key"`
	Value string `json:"value"`
}

// QueuePropertyPage is a single page maximum result representing a query by offset page.
type QueuePropertyPage struct {
	pagination.OffsetPageBase
}

// ExtractProperties is a method to extract the list of properties.
func ExtractProperties(r pagination.Page) ([]PropertyResp, error) {
	var s []PropertyResp
	err := r.(QueuePropertyPage).Result.ExtractIntoSlicePtr(&s, "properties")
	return s, err
}

// QueuePropertyResp is the structure that represents response of UpdateQueueProperty or DeleteQueueProperties method.
type QueuePropertyResp struct {
	// Whether the request is successfully executed. Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
}
