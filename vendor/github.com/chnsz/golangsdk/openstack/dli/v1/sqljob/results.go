package sqljob

import "github.com/chnsz/golangsdk/openstack/common/tags"

const (
	JobTypeDDL           = "DDL"
	JobTypeDCL           = "DCL"
	JobTypeImport        = "IMPORT"
	JobTypeExport        = "EXPORT"
	JobTypeQuery         = "QUERY"
	JobTypeInsert        = "INSERT"
	JobTypeDataMigration = "DATA_MIGRATION"
	JobTypeUpdate        = "UPDATE"
	JobTypeDelete        = "DELETE"
	JobTypeRestartQueue  = "RESTART_QUEUE"
	JobTypeScaleQueue    = "SCALE_QUEUE"

	JobModeSync  = "synchronous"
	JobModeAsync = "asynchronous"

	JobStatusLaunching = "LAUNCHING"
	JobStatusRunning   = "RUNNING"
	JobStatusFinished  = "FINISHED"
	JobStatusFailed    = "FAILED"
	JobStatusCancelled = "CANCELLED"
)

type SubmitJobResult struct {
	// Indicates whether the request is successfully sent. Value true indicates that the request is successfully sent.
	IsSuccess bool `json:"is_success"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message"`
	// ID of a job returned after a job is generated and submitted by using SQL statements.
	// The job ID can be used to query the job status and results.
	JobId string `json:"job_id"`
	// Type of a job. Job types include the following:
	// DDL
	// DCL
	// IMPORT
	// EXPORT
	// QUERY
	// INSERT
	JobType string `json:"job_type"`
	// If the statement type is DDL, the column name and type of DDL are displayed.
	Schema []map[string]string `json:"schema"`
	// When the statement type is DDL, results of the DDL are displayed.
	Rows [][]string `json:"rows"`
	// Job execution mode. The options are as follows:
	// async: asynchronous
	// sync: synchronous
	JobMode string `json:"job_mode"`
}

type CommonResp struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
}

type ListJobsResp struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	JobCount  int    `json:"job_count"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	// Job ID.
	JobId string `json:"job_id"`
	// Type of a job.
	JobType string `json:"job_type"`
	// Queue to which a job is submitted.
	QueueName string `json:"queue_name"`
	// User who submits a job.
	Owner string `json:"owner"`
	// Time when a job is started. The timestamp is expressed in milliseconds.
	StartTime int `json:"start_time"`
	// Job running duration (unit: millisecond).
	Duration int `json:"duration"`
	// Status of a job, including LAUNCHING, RUNNING, FINISHED, FAILED, and CANCELLED.
	Status string `json:"status"`
	// Number of records scanned during the Insert job execution.
	InputRowCount int `json:"input_row_count"`
	// Number of error records scanned during the Insert job execution.
	BadRowCount int `json:"bad_row_count"`
	// Size of scanned files during job execution.
	InputSize int `json:"input_size"`
	// Total number of records returned by the current job or total number of records inserted by the Insert job.
	ResultCount int `json:"result_count"`
	// Name of the database where the target table resides.
	// database_name is valid only for jobs of the Import and Export types.
	DatabaseName string `json:"database_name"`
	// Name of the target table. table_name is valid only for jobs of the Import and Export types.
	TableName string `json:"table_name"`
	// Import jobs, which record whether the imported data contains column names.
	WithColumnHeader bool `json:"with_column_header"`
	// JSON character string of related columns queried by using SQL statements.
	Detail string `json:"detail"`
	// SQL statements of a job.
	Statement string             `json:"statement"`
	Tags      []tags.ResourceTag `json:"tags"`
}

type JobStatus struct {
	// Whether the request is successfully executed. Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success" required:"true"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message" required:"true"`
	// Job ID.
	JobId string `json:"job_id" required:"true"`
	// Type of a job, Includes DDL, DCL, IMPORT, EXPORT, QUERY, INSERT, DATA_MIGRATION, UPDATE, DELETE, RESTART_QUEUE and SCALE_QUEUE.
	JobType string `json:"job_type" required:"true"`
	// Job execution mode. The options are as follows:
	// async: asynchronous
	// sync: synchronous
	JobMode string `json:"job_mode" required:"true"`
	// Name of the queue where the job is submitted.
	QueueName string `json:"queue_name" required:"true"`
	// User who submits a job.
	Owner string `json:"owner" required:"true"`
	// Time when a job is started. The timestamp is expressed in milliseconds.
	StartTime int `json:"start_time" required:"true"`
	// Job running duration (unit: millisecond).
	Duration int `json:"duration"`
	// Status of a job, including RUNNING, SCALING, LAUNCHING, FINISHED, FAILED, and CANCELLED.
	Status string `json:"status" required:"true"`
	// Number of records scanned during the Insert job execution.
	InputRowCount int `json:"input_row_count"`
	// Number of error records scanned during the Insert job execution.
	BadRowCount int `json:"bad_row_count"`
	// Size of scanned files during job execution (unit: byte).
	InputSize int `json:"input_size" required:"true"`
	// Total number of records returned by the current job or total number of records inserted by the Insert job.
	ResultCount int `json:"result_count" required:"true"`
	// Name of the database where the target table resides. database_name is valid only for jobs of the IMPORT EXPORT, and QUERY types.
	DatabaseName string `json:"database_name"`
	// Name of the target table. table_name is valid only for jobs of the IMPORT EXPORT, and QUERY types.
	TableName string `json:"table_name"`
	// JSON character string for information about related columns.
	Detail string `json:"detail" required:"true"`
	// SQL statements of a job.
	Statement string             `json:"statement" required:"true"`
	Tags      []tags.ResourceTag `json:"tags"`
}

type JobDetail struct {
	// Whether the request is successfully executed. Value true indicates that the request is successfully executed.
	IsSuccess bool `json:"is_success" required:"true"`
	// System prompt. If execution succeeds, the parameter setting may be left blank.
	Message string `json:"message" required:"true"`
	// Job ID.
	JobId string `json:"job_id" required:"true"`
	// User who submits a job.
	Owner string `json:"owner" required:"true"`
	// Time when a job is started. The timestamp is expressed in milliseconds.
	StartTime int `json:"start_time" required:"true"`
	// Duration for executing the job (unit: millisecond).
	Duration int `json:"duration" required:"true"`
	// Specified export mode during data export and query result saving.
	ExportMode string `json:"export_mode"`
	// Path to imported or exported files.
	DataPath string `json:"data_path" required:"true"`
	// Type of data to be imported or exported. Currently, only CSV and JSON are supported.
	DataType string `json:"data_type" required:"true"`
	// Name of the database where the table, where data is imported or exported, resides.
	DatabaseName string `json:"database_name" required:"true"`
	// Name of the table where data is imported or exported.
	TableName string `json:"table_name" required:"true"`
	// Whether the imported data contains the column name during the execution of an import job.
	WithColumnHeader bool `json:"with_column_header"`
	// User-defined data delimiter set when the import job is executed.
	Delimiter string `json:"delimiter"`
	// User-defined quotation character set when the import job is executed.
	QuoteChar string `json:"quote_char"`
	// User-defined escape character set when the import job is executed.
	EscapeChar string `json:"escape_char"`
	// Table date format specified when the import job is executed.
	DateFormat string `json:"date_format"`
	// Table time format specified when the import job is executed.
	TimestampFormat string `json:"timestamp_format"`
	// Compression mode specified when the export job is executed.
	Compress string `json:"compress"`
}

type CheckSqlResult struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	// Type of a job. Job types include the following: DDL, DCL, IMPORT, EXPORT, QUERY, and INSERT.
	JobType string `json:"job_type"`
}

type JobResp struct {
	IsSuccess bool   `json:"is_success" required:"true"`
	Message   string `json:"message" required:"true"`
	// ID of a job returned after a job is generated and submitted by using SQL statements.
	// The job ID can be used to query the job status and results.
	JobId string `json:"job_id"`
	// Job execution mode. The options are as follows:
	// async: asynchronous
	// sync: synchronous
	JobMode string `json:"job_mode"`
}

type JobProgress struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	JobId     string `json:"job_id"`
	Status    string `json:"status"`
	// ID of a subjob that is running. If the subjob is not running or it is already finished,
	// the subjob ID may be empty.
	SubJobId int `json:"sub_job_id"`
	// Progress of a running subjob or the entire job. The value can only be a rough estimate of the subjob progress
	// and does not indicate the detailed job progress.

	// If the job is just started or being submitted, the progress is displayed as 0. If the job execution is complete,
	//  the progress is displayed as 1. In this case, progress indicates the running progress of the entire job.
	//  Because no subjob is running, sub_job_id is not displayed.
	// If a subjob is running, the running progress of the subjob is displayed. The calculation method of progress is as
	//  follows: Number of completed tasks of the subjob/Total number of tasks of the subjob. In this case,
	//  progress indicates the running progress of the subjob, and sub_job_id indicates the subjob ID.
	Progress int `json:"progress"`
	// Details about a subjob of a running job. A job may contain multiple subjobs. For details
	SubJobs []SubJob `json:"sub_jobs"`
}

type SubJob struct {
	// Subjob ID, corresponding to jobId of the open-source spark JobData.
	Id int `json:"id"`
	// Subjob name, corresponding to the name of the open-source spark JobData.
	Name string `json:"name"`
	// Description of a subjob, corresponding to the description of the open-source spark JobData.
	Description string `json:"description"`
	// Submission time of a subjob, corresponding to the submissionTime of open-source Spark JobData.
	SubmissionTime string `json:"submission_time"`
	// Completion time of a subjob, corresponding to the completionTime of the open-source Spark JobData.
	CompletionTime string `json:"completion_time"`
	// Stage ID of the subjob, corresponding to the stageIds of the open-source spark JobData.
	StageIds []int `json:"stage_ids"`
	// ID of a DLI job, corresponding to the jobGroup of open-source Spark JobData.
	JobGroup string `json:"job_group"`
	// Subjob status, corresponding to the status of open-source spark JobData.
	Status string `json:"status"`
	// Number of subjobs, corresponding to numTasks of the open-source Spark JobData.
	NumTasks int `json:"num_tasks"`
	// Number of running tasks in a subjob, corresponding to numActiveTasks of the open-source Spark JobData.
	NumActiveTasks int `json:"num_active_tasks"`
	// Number of tasks that have been completed in a subjob, corresponding to numCompletedTasks of open-source Spark JobData.
	NumCompletedTasks int `json:"num_completed_tasks"`
	// Number of tasks skipped in a subjob, corresponding to numSkippedTasks of open-source Spark JobData.
	NumSkippedTasks int `json:"num_skipped_tasks"`
	// Number of subtasks that fail to be skipped, corresponding to numFailedTasks of open-source Spark JobData.
	NumFailedTasks int `json:"num_failed_tasks"`
	// Number of tasks killed in the subjob, corresponding to numKilledTasks of the open-source Spark JobData.
	NumKilledTasks int `json:"num_killed_tasks"`
	// Subjob completion index, corresponding to the numCompletedIndices of the open-source Spark JobData.
	NumCompletedIndices int `json:"num_completed_indices"`
	// Number of stages that are running in the subjob, corresponding to numActiveStages of the open-source Spark JobData.
	NumActiveStages int `json:"num_active_stages"`
	// Number of stages that have been completed in the subjob, corresponding to numCompletedStages of the open-source Spark JobData.
	NumCompletedStages int `json:"num_completed_stages"`
	// Number of stages skipped in the subjob, corresponding to numSkippedStages of the open-source Spark JobData.
	NumSkippedStages int `json:"num_skipped_stages"`
	// Number of failed stages in a subjob, corresponding to numFailedStages of the open-source Spark JobData.
	NumFailedStages int `json:"num_failed_stages"`
	// Summary of the killed tasks in the subjob, corresponding to killedTasksSummary of open-source spark JobData.
	KilledTasksSummary map[string]int `json:"killed_tasks_summary"`
}
