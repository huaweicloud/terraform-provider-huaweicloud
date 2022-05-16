package jobs

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
