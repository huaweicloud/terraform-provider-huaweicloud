package jobs

// Job is the structure of the API response Detail.
type Job struct {
	// Specifies the task status.
	//   SUCCESS: indicates the task is successfully executed.
	//   RUNNING: indicates that the task is in progress.
	//   FAIL: indicates that the task failed.
	//   INIT: indicates that the task is being initialized.
	//   PENDING_PAYMENT: indicates that a yearly/monthly order is to be paid.
	// NOTE:
	//   The PENDING_PAYMENT status is displayed after the request for creating a yearly/monthly ECS or modifying the
	//   specifications of yearly/monthly ECS has been submitted and before the order is paid. If the order is canceled,
	//   the status will not be automatically updated. The task will be automatically deleted 14 days later.
	Status string `json:"status"`
	// Specifies the object of the task.
	// The value of this parameter varies depending on the type of the task. If the task is an ECS-related operation,
	// the value is server_id. If the task is a NIC operation, the value is nic_id. If a sub-Job is available, details
	// about the sub-job are displayed.
	Entities JobEntity `json:"entities"`
	// Specifies the ID of an asynchronous request task.
	ID string `json:"job_id"`
	// Specifies the type of an asynchronous request task.
	Type string `json:"job_type"`
	// Specifies the time when the task started.
	BeginTime string `json:"begin_time"`
	// Specifies the time when the task finished.
	EndTime string `json:"end_time"`
	// Specifies the returned error code when the task execution fails.
	// After the task is executed successfully, the value of this parameter is null.
	ErrorCode string `json:"error_code"`
	// Specifies the cause of the task execution failure.
	// After the task is executed successfully, the value of this parameter is null.
	FailReason string `json:"fail_reason"`
	// Specifies the error message returned when an error occurs in the request to query a task.
	Message string `json:"message"`
	// Specifies the error code returned when an error occurs in the request to query a task.
	// For details about the error code, see Returned Values for General Requests.
	Code string `json:"code"`
}

// JobEntity is structure that shows all sub-jobs and the total count.
type JobEntity struct {
	// Specifies the number of subtasks.
	SubJobsTotal int `json:"sub_jobs_total"`
	// Specifies the execution information of a subtask.
	SubJobs []SubJob `json:"sub_jobs"`
}

// SubJob is the structure of the execution details for each subn-job under the Job.
type SubJob struct {
	// Specifies the task status.
	// SUCCESS: indicates the task is successfully executed.
	// RUNNING: indicates that the task is in progress.
	// FAIL: indicates that the task failed.
	// INIT: indicates that the task is being initialized.
	Status string `json:"status"`
	// Specifies the object of the task. The value of this parameter varies depending on the type of the task. If the task is an ECS-related operation, the value is server_id. If the task is a NIC operation, the value is nic_id. For details, see Table 5.
	Entities SubJobEntity `json:"entities"`
	// Specifies the subtask ID.
	ID string `json:"job_id"`
	// Specify the subtask type.
	Type string `json:"job_type"`
	// Specifies the time when the task started.
	BeginTime string `json:"begin_time"`
	// Specifies the time when the task finished.
	EndTime string `json:"end_time"`
	// Specifies the returned error code when the task execution fails.
	// After the task is executed successfully, the value of this parameter is null.
	ErrorCode string `json:"error_code"`
	// Specifies the cause of the task execution failure.
	// After the task is executed successfully, the value of this parameter is null.
	FailReason string `json:"fail_reason"`
}

// SubJobEntity is structure of the sub-job operation.
type SubJobEntity struct {
	// If the task is an ECS-related operation, the value is server_id.
	ServerId string `json:"server_id"`
	// If the task is a NIC-related operation, the value is nic_id.
	NicId string `json:"nic_id"`
	// Indicates the cause of a subtask execution failure.
	ErrorcodeMessage string `json:"errorcode_message"`
}
