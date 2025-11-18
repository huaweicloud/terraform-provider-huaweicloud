package flinkjob

type CreateJobResp struct {
	IsSuccess bool      `json:"is_success,string"`
	Message   string    `json:"message"`
	Job       JobStatus `json:"job"`
}

type JobStatus struct {
	JobId      int    `json:"job_id"`
	StatusName string `json:"status_name"`
	StatusDesc string `json:"status_desc"`
}

type UpdateJobResp struct {
	IsSuccess bool              `json:"is_success,string"`
	Message   string            `json:"message"`
	Job       UpdateJobResp_job `json:"job"`
}

type UpdateJobResp_job struct {
	UpdateTime int `json:"update_time"`
}

type CommonResp struct {
	IsSuccess bool   `json:"is_success,string"`
	Message   string `json:"message"`
}

type GetJobResp struct {
	IsSuccess bool   `json:"is_success,string"`
	Message   string `json:"message"`
	JobDetail Job    `json:"job_detail"`
}

type Job struct {
	// Job ID.
	JobId int `json:"job_id"`
	// Name of the job. Length range: 0 to 57 characters.
	Name string `json:"name"`
	// Job description. Length range: 0 to 512 characters.
	Desc string `json:"desc"`
	// Job type.
	// flink_sql_job: Flink SQL job
	// flink_opensource_sql_job: Flink OpenSource SQL job
	// flink_jar_job: User-defined Flink job
	JobType string `json:"job_type"`
	// Job status.
	// Available job statuses are as follows:
	// job_init: The job is in the draft status.
	// job_submitting: The job is being submitted.
	// job_submit_fail: The job fails to be submitted.
	// job_running: The job is running. (The billing starts. After the job is submitted, a normal result is returned.)
	// job_running_exception (The billing stops. The job stops running due to an exception.)
	// job_downloading: The job is being downloaded.
	// job_idle: The job is idle.
	// job_canceling: The job is being stopped.
	// job_cancel_success: The job has been stopped.
	// job_cancel_fail: The job fails to be stopped.
	// job_savepointing: The savepoint is being created.
	// job_arrearage_stopped: The job is stopped because the account is in arrears.
	//  (The billing ends. The job is stopped because the user account is in arrears.)
	// job_arrearage_recovering: The recharged job is being restored.
	//  (The account in arrears is recharged, and the job is being restored).
	// job_finish: The job is completed.
	Status string `json:"status"`
	// Description of job status.
	StatusDesc string `json:"status_desc"`
	// Time when a job is created.
	CreateTime int `json:"create_time"`
	// Time when a job is started.
	StartTime int `json:"start_time"`
	// ID of the user who creates the job.
	UserId string `json:"user_id"`
	// Name of a queue. Length range: 1 to 128 characters.
	QueueName string `json:"queue_name"`
	// ID of the project to which a job belongs.
	ProjectId string `json:"project_id"`
	// Stream SQL statement.
	SqlBody string `json:"sql_body"`
	// Job running mode. The options are as follows:
	// shared_cluster: indicates that the job is running on a shared cluster.
	// exclusive_cluster: indicates that the job is running on an exclusive cluster.
	// edge_node: indicates that the job is running on an edge node.
	RunMode string `json:"run_mode"`
	// Job configurations. Refer to Table 4 for details.
	JobConfig JobConf `json:"job_config"`
	// Main class of a JAR package, for example, org.apache.spark.examples.streaming.JavaQueueStream.
	MainClass string `json:"main_class"`
	// Running parameter of a JAR package job. Multiple parameters are separated by spaces.
	EntrypointArgs string `json:"entrypoint_args"`
	// Job execution plan.
	ExecutionGraph string `json:"execution_graph"`
	// Time when a job is updated.
	UpdateTime int `json:"update_time"`
	// User-defined job feature. Type of the Flink image used by a job.
	// basic: indicates that the basic Flink image provided by DLI is used.
	// custom: indicates that the user-defined Flink image is used.
	Feature string `json:"feature"`
	// Flink version. This parameter is valid only when feature is set to basic. You can use this parameter with the
	// feature parameter to specify the version of the DLI basic Flink image used for job running.
	FlinkVersion string `json:"flink_version"`
	// Custom image. The format is Organization name/Image name:Image version.
	// This parameter is valid only when feature is set to custom. You can use this parameter with the feature
	// parameter to specify a user-defined Flink image for job running. For details about how to use custom images.
	Image string `json:"image"`
}

type JobConfBase struct {
	// Whether to enable the automatic job snapshot function.
	// true: The automatic job snapshot function is enabled.
	// false: The automatic job snapshot function is disabled.
	// The default value is false.
	CheckpointEnabled bool `json:"checkpoint_enabled"`
	// Snapshot mode. There are two options:
	// exactly_once: indicates that data is processed only once.
	// at_least_once: indicates that data is processed at least once.
	// The default value is exactly_once.
	CheckpointMode string `json:"checkpoint_mode"`
	// Snapshot interval. The unit is second. The default value is 10.
	CheckpointInterval int `json:"checkpoint_interval"`
	// Whether to enable the log storage function. The default value is false.
	LogEnabled bool `json:"log_enabled"`
	// Name of an OBS bucket.
	ObsBucket string `json:"obs_bucket"`
	// SMN topic name. If a job fails, the system will send a message to users subscribed to the SMN topic.
	SmnTopic string `json:"smn_topic"`
	// Parent job ID.
	RootId int `json:"root_id"`
	// List of edge computing group IDs. Use commas (,) to separate multiple IDs.
	EdgeGroupIds []string `json:"edge_group_ids"`
	// Number of CUs of the management unit. The default value is 1.
	ManagerCuNumber int `json:"manager_cu_number"`
	// Number of CUs selected for a job. This parameter is valid only when show_detail is set to true.
	// Minimum value: 2
	// Maximum value: 400
	// The default value is 2.
	CuNumber int `json:"cu_number"`
	// Number of concurrent jobs set by a user. This parameter is valid only when show_detail is set to true.
	// Minimum value: 1
	// Maximum value: 2000
	// The default value is 1.
	ParallelNumber int `json:"parallel_number"`
	// Whether to enable the function of restart upon exceptions.
	RestartWhenException bool `json:"restart_when_exception"`
	// Expiration time.
	IdleStateRetention int `json:"idle_state_retention"`
	// Name of the package that has been uploaded to the DLI resource management system. The UDF Jar file of the SQL
	// job is uploaded through this parameter.
	UdfJarUrl string `json:"udf_jar_url"`
	// Dirty data policy of a job.
	// 2:obsDir: Save. obsDir specifies the path for storing dirty data.
	// 1: Trigger a job exception
	// 0: Ignore
	DirtyDataStrategy string `json:"dirty_data_strategy"`
	// Name of the package that has been uploaded to the DLI resource management system.
	// This parameter is used to customize the JAR file where the job main class is located.
	Entrypoint string `json:"entrypoint"`
	// Name of the package that has been uploaded to the DLI resource management system.
	// This parameter is used to customize other dependency packages.
	DependencyJars []string `json:"dependency_jars"`
	// Name of the resource package that has been uploaded to the DLI resource management system.
	// This parameter is used to customize dependency files.
	DependencyFiles []string `json:"dependency_files"`
	// Number of compute nodes in a job.
	ExecutorNumber int `json:"executor_number"`
	// Number of CUs in a compute node.
	ExecutorCuNumber int `json:"executor_cu_number"`
	// Whether to restore data from the latest checkpoint when the system automatically restarts upon an exception.
	// The default value is false.
	ResumeCheckpoint bool   `json:"resume_checkpoint"`
	TmCus            int    `json:"tm_cus"`
	TmSlotNum        int    `json:"tm_slot_num"`
	ResumeMaxNum     int    `json:"resume_max_num"`
	CheckpointPath   string `json:"checkpoint_path"`
	Feature          string `json:"feature"`
	FlinkVersion     string `json:"flink_version"`
	Image            string `json:"image"`
	// Degree of parallelism (DOP) of an operator.
	OperatorConfig string `json:"operator_config"`
	// The traffic or hit rate configuration of each operator.
	StaticEstimatorConfig string `json:"static_estimator_config"`
	// The name of the delegation authorized to DLI.
	ExecutionAgencyUrn string `json:"execution_agency_urn"`
	// The version of the resource configuration.
	ResourceConfigVersion string `json:"resource_config_version"`
	// The resource configuration of the Flink job.
	ResourceConfig ResourceConfig `json:"resource_config"`
}

type ResourceConfig struct {
	// The number of slots in the JobManager.
	MaxSlot int `json:"max_slot"`
	// The parallel number of the job.
	ParallelNumber int `json:"parallel_number"`
	// The resource specification of the JobManager.
	JobManagerResourceSpec ManagerResourceSpec `json:"jobmanager_resource_spec"`
	// The resource specification of the TaskManager.
	TaskManagerResourceSpec ManagerResourceSpec `json:"taskmanager_resource_spec"`
}

type ManagerResourceSpec struct {
	// The number of CPU cores that the JobManager or TaskManager can use.
	CPU float64 `json:"cpu"`
	// The memory size that the JobManager or TaskManager can use, in MB or GB (default).
	Memory string `json:"memory"`
}

type ListResp struct {
	IsSuccess bool          `json:"is_success,string"`
	Message   string        `json:"message"`
	JobList   JobListWapper `json:"job_list"`
}

type JobListWapper struct {
	TotalCount int        `json:"total_count"`
	Jobs       []Job4List `json:"jobs"`
}

type Job4List struct {
	JobId int    `json:"job_id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	// Job description. Length range: 0 to 512 characters.
	Username   string `json:"username"`
	JobType    string `json:"job_type"`
	Status     string `json:"status"`
	StatusDesc string `json:"status_desc"`
	CreateTime int    `json:"create_time"`
	StartTime  int    `json:"start_time"`
	// Running duration of a job. Unit: ms. This parameter is valid only when show_detail is set to false.
	Duration int `json:"duration"`
	// Parent job ID. This parameter is valid only when show_detail is set to false.
	RootId int `json:"root_id"`
	// ID of the user who creates the job. This parameter is valid only when show_detail is set to true.
	UserId string `json:"user_id"`
	// This parameter is valid only when show_detail is set to true.
	ProjectId string `json:"project_id"`
	// Stream SQL statement. This parameter is valid only when show_detail is set to false.
	SqlBody string `json:"sql_body"`
	// Job running mode. The options are as follows: The value can be shared_cluster, exclusive_cluster, or edge_node.
	// This parameter is valid only when show_detail is set to true.
	// shared_cluster: indicates that the job is running on a shared cluster.
	// exclusive_cluster: indicates that the job is running on an exclusive cluster.
	// edge_node: indicates that the job is running on an edge node.
	RunMode string `json:"run_mode"`
	// Job configuration. This parameter is valid only when show_detail is set to false.
	JobConfig JobConfBase `json:"job_config"`
	//Main class of a JAR package. This parameter is valid only when show_detail is set to false.
	MainClass string `json:"main_class"`
	// Job running parameter of the JAR file. Multiple parameters are separated by spaces.
	// This parameter is valid only when show_detail is set to true.
	EntrypointArgs string `json:"entrypoint_args"`
	// Job execution plan. This parameter is valid only when show_detail is set to false.
	ExecutionGraph string `json:"execution_graph"`
	// Time when a job is updated. This parameter is valid only when show_detail is set to false.
	UpdateTime int `json:"update_time"`
}

type JobConf struct {
	JobConfBase
	// Customizes optimization parameters when a Flink job is running.
	RuntimeConfig string `json:"runtime_config"`
}

type DliError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}
