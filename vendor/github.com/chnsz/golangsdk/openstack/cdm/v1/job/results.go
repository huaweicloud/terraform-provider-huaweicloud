package job

type CreateResponse struct {
	Name             string             `json:"name"`
	ValidationResult []ValidationDetail `json:"validation-result"`
}

type ValidationDetail struct {
	Message string `json:"message"`
	// ERROR,WARNING
	Status string `json:"status"`
}

type validationResult struct {
	LinkConfig []ValidationDetail `json:"linkConfig"`
}

type JobsDetail struct {
	Total    int   `json:"total"`
	Jobs     []Job `json:"jobs"`
	PageNo   int   `json:"page_no"`
	PageSize int   `json:"page_size"`
	// Whether to return compact information. If this parameter is set to true, compact information will be returned,
	// which means only parameter names and values will be returned.
	// Attributes such as size, type, and id will not be returned.
	Simple bool `json:"simple"`
}

type UpdateResponse struct {
	ValidationResult []validationResult `json:"validation-result"`
}

type ErrorResponse struct {
	// Error code
	ErrCode string `json:"errCode"`
	// Error message
	ErrMessage string `json:"externalMessage"`
}

type StartJobResponse struct {
	Submissions []JobSubmission `json:"submissions"`
}

type JobSubmission struct {
	DeleteRows   int    `json:"delete_rows"`
	UpdateRows   int    `json:"update_rows"`
	WriteRows    int    `json:"write_rows"`
	SubmissionId int    `json:"submission-id"`
	JobName      string `json:"job-name"`
	CreationUser string `json:"creation-user"`
	CreationDate int    `json:"creation-date"`
	// Job progress. If a job fails, the value is -1. Otherwise, the value ranges from 0 to 100.
	Progress float32 `json:"progress"`
	// Job status. The options are as follows:
	// BOOTING: The job is starting.
	// FAILURE_ON_SUBMIT: The job fails to be submitted.
	// RUNNING: The job is running.
	// SUCCEEDED: The job is executed successfully.
	// FAILED: The job failed.
	// UNKNOWN: The job status is unknown.
	// NEVER_EXECUTED: The job has not been executed.
	Status             string `json:"status"`
	IsStopingIncrement string `json:"isStopingIncrement"`
	IsExecuteAuto      bool   `json:"is-execute-auto"`
	// Whether the job involves incremental data migration
	IsIncrementing bool `json:"isIncrementing"`

	LastUpdateDate int    `json:"last-update-date"`
	LastUdpateUser string `json:"last-udpate-user"`
	// Whether to delete the job after it is executed
	IsDeleteJob bool `json:"isDeleteJob"`
}

type StatusResponse struct {
	Submissions []StatusSubmission `json:"submissions"`
}

type StatusSubmission struct {
	JobSubmission

	// Job running result statistics. This parameter is available only when status is SUCCEEDED.
	Counters     Counters `json:"counters"`
	ExternalId   string   `json:"external-id"`
	ExecuteDate  int      `json:"execute-date"`
	ErrorDetails string   `json:"error-details"`
	ErrorSummary string   `json:"error-summary"`
}

type Counters struct {
	SqoopCounters Counter `json:"org.apache.sqoop.submission.counter.SqoopCounters"`
}

type Counter struct {
	BytesWritten       int `json:"BYTES_WRITTEN"`
	TotalFiles         int `json:"TOTAL_FILES"`
	RowsRead           int `json:"ROWS_READ"`
	BytesRead          int `json:"BYTES_READ"`
	RowsWritten        int `json:"ROWS_WRITTEN"`
	FilesWritten       int `json:"FILES_WRITTEN"`
	FilesRead          int `json:"FILES_READ"`
	TotalSize          int `json:"TOTAL_SIZE"`
	FilesSkipped       int `json:"FILES_SKIPPED"`
	RowsWrittenSkipped int `json:"ROWS_WRITTEN_SKIPPED"`
}

type ListSubmissionsRst struct {
	Submissions []StatusSubmission `json:"submissions"`
	Total       int                `json:"total"`
	PageNo      int                `json:"page_no"`
	PageSize    int                `json:"page_size"`
}
