package jobs

// QueryResp is the structure that represents the API response of Create method request.
type QueryResp struct {
	// The total count of job list.
	TotalCount int `json:"total_count"`
	// The job list.
	Jobs []Job `json:"jobs"`
}

// Job is an object to specified the operation detail of desktop.
type Job struct {
	// Job ID.
	ID string `json:"id"`
	// Job type.
	Type string `json:"job_type"`
	// Job entity.
	Entities Entities `json:"entities"`
	// Start time.
	BeginTime string `json:"begin_time"`
	// End time.
	EndTime string `json:"end_time"`
	// Job status.
	Status string `json:"status"`
	// Error code.
	ErrorCode string `json:"error_code"`
	// The fail reason.
	FailReason string `json:"fail_reason"`
	// The message of fail reason.
	Message string `json:"message"`
	// Sub-job ID.
	SubJobID string `json:"job_id"`
}

// Entities is an object to specified the job entity details of desktop.
type Entities struct {
	// Desktop ID.
	DesktopId string `json:"desktop_id"`
	// Product ID.
	ProductId string `json:"product_id"`
	// User name.
	UserName string `json:"user_name"`
}
