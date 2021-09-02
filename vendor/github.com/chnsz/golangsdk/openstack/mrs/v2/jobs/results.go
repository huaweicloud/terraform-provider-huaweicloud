package jobs

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	golangsdk.Result
}

// CreateResp is a struct that represents the result of Create methods.
type CreateResp struct {
	// Job execution result
	JobSubmitResult JobResp `json:"job_submit_result"`
	// Error message
	ErrorMsg string `json:"error_msg"`
	// Error code
	ErrorCode string `json:"error_code"`
}

// JobResp is an object struct that represents the details of the submit result.
type JobResp struct {
	// Job ID
	JobId string `json:"job_id"`
	// Job submission status.
	// COMPLETE: The job is submitted.
	// JOBSTAT_SUBMIT_FAILED: Failed to submit the job.
	State string `json:"state"`
}

// Extract is a method which to extract a creation response.
func (r CreateResult) Extract() (*CreateResp, error) {
	var s CreateResp
	err := r.ExtractInto(&s)
	return &s, err
}

// Job is a struct that represents the details of the mapreduce job.
type Job struct {
	// Job ID.
	JobId string `json:"job_id"`
	// Name of the user who submits a job.
	User string `json:"user"`
	// Job name. It contains 1 to 64 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	JobName string `json:"job_name"`
	// Final result of a job.
	//   FAILED: indicates that the job fails to be executed.
	//   KILLED: indicates that the job is manually terminated during execution.
	//   UNDEFINED: indicates that the job is being executed.
	//   SUCCEEDED: indicates that the job has been successfully executed.
	JobResult string `json:"job_result"`
	// Execution status of a job.
	//   FAILED: failed
	//   KILLED: indicates that the job is terminated.
	//   New: indicates that the job is created.
	//   NEW_SAVING: indicates that the job has been created and is being saved.
	//   SUBMITTED: indicates that the job is submitted.
	//   ACCEPTED: indicates that the job is accepted.
	//   RUNNING: indicates that the job is running.
	//   FINISHED: indicates that the job is completed.
	JobState string `json:"job_state"`
	// Job execution progress.
	JobProgress float32 `json:"job_progress"`
	// Type of a job, which support:
	//   MapReduce
	//   SparkSubmit
	//   HiveScript
	//   HiveSql
	//   DistCp, importing and exporting data
	//   SparkScript
	//   SparkSql
	//   Flink
	JobType string `json:"job_type"`
	// Start time to run a job. Unit: ms.
	StartedTime int `json:"started_time"`
	// Time when a job is submitted. Unit: ms.
	SubmittedTime int `json:"submitted_time"`
	// End time to run a job. Unit: ms.
	FinishedTime int `json:"finished_time"`
	// Running duration of a job. Unit: ms.
	ElapsedTime int `json:"elapsed_time"`
	// Running parameter. The parameter contains a maximum of 4,096 characters,
	// excluding special characters such as ;|&>'<$, and can be left blank.
	Arguments string `json:"arguments"`
	// Configuration parameter, which is used to configure -d parameters.
	// The parameter contains a maximum of 2,048 characters, excluding special characters such as ><|'`&!\,
	// and can be left blank.
	Properties string `json:"properties"`
	// Launcher job ID.
	LauncherId string `json:"launcher_id"`
	// Actual job ID.
	AppId string `json:"app_id"`
}

// GetResult represents a result of the Get method.
type GetResult struct {
	commonResult
}

func (r commonResult) Extract() (*Job, error) {
	var s Job
	err := r.ExtractIntoStructPtr(&s, "job_detail")
	return &s, err
}

// JobPage represents the response pages of the List operation.
type JobPage struct {
	pagination.SinglePageBase
}

// ExtractJobs is a method which to extract a job list by job pages.
func ExtractJobs(r pagination.Page) ([]Job, error) {
	var s []Job
	err := r.(JobPage).Result.ExtractIntoSlicePtr(&s, "job_list")
	return s, err
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	golangsdk.ErrResult
}
