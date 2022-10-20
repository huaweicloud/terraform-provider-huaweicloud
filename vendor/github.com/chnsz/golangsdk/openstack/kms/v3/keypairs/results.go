package keypairs

import (
	"github.com/chnsz/golangsdk"
)

// TaskResp is the response body of Associate or Disassociate
type TaskResp struct {
	// the task ID
	ID string `json:"task_id"`
}

// Task contains all the information about a keypair task
type Task struct {
	// the task ID
	ID string `json:"task_id"`
	// the ECS instance ID
	ServerID string `json:"server_id"`
	// the keypair processing state
	Status string `json:"task_status"`
}

// GetResult contains the response body and error from a GetTask request.
type GetResult struct {
	golangsdk.Result
}

// Extract the Task from result
func (r GetResult) Extract() (Task, error) {
	var s Task
	err := r.ExtractInto(&s)
	return s, err
}
