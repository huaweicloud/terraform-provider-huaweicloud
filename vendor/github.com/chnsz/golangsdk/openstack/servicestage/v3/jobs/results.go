package jobs

import "github.com/chnsz/golangsdk/pagination"

// JobResp is the structure that represents the detail of the deployment job and task list.
type JobResp struct {
	// Number of tasks.
	TaskCount int `json:"task_count"`
	// Job parameters.
	Job Job `json:"job"`
	// Task parameters.
	Tasks []Task `json:"tasks"`
}

// Job is the structure that represents the detail of the deployment action.
type Job struct {
	// Creator.
	Creator string `json:"created_by"`
	// Execution status.
	ExecutionStatus string `json:"execution_status"`
	// Job description.
	Description string `json:"job_desc"`
	// Job ID.
	ID string `json:"job_id"`
	// Job name.
	Name string `json:"job_name"`
	// Type.
	Type string `json:"job_type"`
	// Order ID.
	OrderId string `json:"order_id"`
	// Tenant's project ID.
	ProjectId string `json:"project_id"`
	// Instance ID.
	InstanceId string `json:"service_instance_id"`
}

// Task is the structure that represents the detail of the deployment task.
type Task struct {
	// Creation time.
	CreatedAt string `json:"created_at"`
	// Health check time.
	LastHealthCheck string `json:"last_health_check"`
	// Message.
	Messages string `json:"messages"`
	// Creator ID.
	OwnerId string `json:"owner_id"`
	// Task ID.
	ID string `json:"task_id"`
	// Task index.
	Index int `json:"task_index"`
	// Task name.
	Name string `json:"task_name"`
	// Task status.
	Status string `json:"task_status"`
	// Task type.
	Type string `json:"task_type"`
}

// TaskPage is a single page maximum result representing a query by offset page.
type TaskPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a TaskPage is empty.
func (b TaskPage) IsEmpty() (bool, error) {
	arr, err := ExtractTasks(b)
	return len(arr) == 0, err
}

// ExtractTasks is a method to extract the list of task details for ServiceStage component.
func ExtractTasks(r pagination.Page) ([]Task, error) {
	var s []Task
	err := r.(TaskPage).Result.ExtractIntoSlicePtr(&s, "tasks")
	return s, err
}
