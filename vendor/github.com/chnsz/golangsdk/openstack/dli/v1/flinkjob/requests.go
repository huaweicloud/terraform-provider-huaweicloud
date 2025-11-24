package flinkjob

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const (
	JobTypeFlinkSql           = "flink_sql_job"
	JobTypeFlinkOpenSourceSql = "flink_opensource_sql_job"
	JobTypeFlinkEdgeSql       = "flink_sql_edge_job"
	JobTypeFlinkJar           = "flink_jar_job"

	RunModeSharedCluster    = "shared_cluster"
	RunModeExclusiveCluster = "exclusive_cluster"
	RunModeEdgeNode         = "edge_node"

	CheckpointModeExactlyOnce = "exactly_once"
	CheckpointModeAtLeastOnce = "at_least_once"
)

type CreateSqlJobOpts struct {
	// Name of the job. Length range: 0 to 57 characters.
	Name string `json:"name" required:"true"`
	// Job description. Length range: 0 to 512 characters.
	Desc string `json:"desc,omitempty"`
	// Template ID.
	// If both template_id and sql_body are specified, sql_body is used. If template_id is specified but sql_body is
	// not, fill sql_body with the template_id value.
	TemplateId *int `json:"template_id,omitempty"`
	// Name of a queue. Length range: 1 to 128 characters.
	QueueName string `json:"queue_name,omitempty"`
	// Stream SQL statement, which includes at least the following three parts: source, query, and sink.
	// Length range: 1024x1024 characters.
	SqlBody string `json:"sql_body,omitempty"`
	// Job running mode. The options are as follows:
	// shared_cluster: indicates that the job is running on a shared cluster.
	// exclusive_cluster: indicates that the job is running on an exclusive cluster.
	// edge_node: indicates that the job is running on an edge node.
	// The default value is shared_cluster.
	RunMode string `json:"run_mode,omitempty"`
	// Number of CUs selected for a job. The default value is 2.
	CuNumber *int `json:"cu_number,omitempty"`
	// Number of parallel jobs set by a user. The default value is 1.
	ParallelNumber *int `json:"parallel_number,omitempty"`
	// Whether to enable the automatic job snapshot function.
	// true: indicates to enable the automatic job snapshot function.
	// false: indicates to disable the automatic job snapshot function.
	// Default value: false
	CheckpointEnabled *bool `json:"checkpoint_enabled,omitempty"`
	// Snapshot mode. There are two options:
	// 1: ExactlyOnce, indicates that data is processed only once.
	// 2: AtLeastOnce, indicates that data is processed at least once.
	// The default value is 1.
	CheckpointMode *int `json:"checkpoint_mode,omitempty"`
	// Snapshot interval. The unit is second. The default value is 10.
	CheckpointInterval *int `json:"checkpoint_interval,omitempty"`
	// OBS path where users are authorized to save the snapshot. This parameter is valid only when checkpoint_enabled
	// is set to true.
	// OBS path where users are authorized to save the snapshot. This parameter is valid only when log_enabled
	// is set to true.
	ObsBucket string `json:"obs_bucket,omitempty"`
	// Whether to enable the function of uploading job logs to users' OBS buckets. The default value is false.
	LogEnabled *bool `json:"log_enabled,omitempty"`
	// SMN topic. If a job fails, the system will send a message to users subscribed to the SMN topic.
	SmnTopic string `json:"smn_topic,omitempty"`
	// Whether to enable the function of automatically restarting a job upon job exceptions. The default value is false.
	RestartWhenException *bool `json:"restart_when_exception,omitempty"`
	// Retention time of the idle state. The unit is hour. The default value is 1.
	IdleStateRetention *int     `json:"idle_state_retention,omitempty"`
	EdgeGroupIds       []string `json:"edge_group_ids,omitempty"`
	// Job type. This parameter can be set to flink_sql_job, and flink_opensource_sql_job.
	// If run_mode is set to shared_cluster or exclusive_cluster, this parameter must be flink_sql_job.
	// The default value is flink_sql_job.
	JobType string `json:"job_type,omitempty"`
	// Dirty data policy of a job.
	// 2:obsDir: Save. obsDir specifies the path for storing dirty data.
	// 1: Trigger a job exception
	// 0: Ignore
	// The default value is 0.
	DirtyDataStrategy string `json:"dirty_data_strategy,omitempty"`
	// Name of the resource package that has been uploaded to the DLI resource management system.
	// The UDF Jar file of the SQL job is specified by this parameter.
	UdfJarUrl string `json:"udf_jar_url,omitempty"`
	// Number of CUs in the JobManager selected for a job. The default value is 1.
	ManagerCuNumber *int `json:"manager_cu_number"`
	// Number of CUs for each TaskManager. The default value is 1.
	TmCus *int `json:"tm_cus,omitempty"`
	// Number of slots in each TaskManager. The default value is (parallel_number*tm_cus)/(cu_number-manager_cu_number).
	TmSlotNum *int `json:"tm_slot_num,omitempty"`
	// Whether the abnormal restart is recovered from the checkpoint.
	ResumeCheckpoint *bool `json:"resume_checkpoint,omitempty"`
	// Maximum number of retry times upon exceptions. The unit is times/hour. Value range: -1 or greater than 0.
	// The default value is -1, indicating that the number of times is unlimited.
	ResumeMaxNum *int `json:"resume_max_num,omitempty"`
	// Customizes optimization parameters when a Flink job is running.
	RuntimeConfig string `json:"runtime_config,omitempty"`
	// Label of a Flink SQL job. For details, see Table 3.
	Tags []tags.ResourceTag `json:"tags"`
	// Flink version. The valid value is `1.1`0 or `1.12`.
	FlinkVersion string `json:"flink_version,omitempty"`
}

type UpdateSqlJobOpts struct {
	// Name of a job. Length range: 0 to 57 characters.
	Name string `json:"name,omitempty"`
	// Job description. Length range: 0 to 512 characters.
	Desc string `json:"desc,omitempty"`
	// Name of a queue. Length range: 1 to 128 characters.
	QueueName string `json:"queue_name,omitempty"`
	// Stream SQL statement, which includes at least the following three parts: source, query, and sink.
	// Length range: 0 to 1024x1024 characters.
	SqlBody string `json:"sql_body,omitempty"`
	// Job running mode. The options are as follows:
	// shared_cluster: indicates that the job is running on a shared cluster.
	// exclusive_cluster: indicates that the job is running on an exclusive cluster.
	// edge_node: indicates that the job is running on an edge node.
	// The default value is shared_cluster.
	RunMode string `json:"run_mode,omitempty"`
	// Number of CUs selected for a job. The default value is 2.
	CuNumber *int `json:"cu_number,omitempty"`
	// Number of parallel jobs set by a user. The default value is 1.
	ParallelNumber *int `json:"parallel_number,omitempty"`
	// Whether to enable the automatic job snapshot function.
	// true: indicates to enable the automatic job snapshot function.
	// false: indicates to disable the automatic job snapshot function.
	// Default value: false
	CheckpointEnabled *bool `json:"checkpoint_enabled,omitempty"`
	// Snapshot mode. There are two options:
	// 1: ExactlyOnce, indicates that data is processed only once.
	// 2: at_least_once, indicates that data is processed at least once.
	// The default value is 1.
	CheckpointMode *int `json:"checkpoint_mode,omitempty"`
	// Snapshot interval. The unit is second. The default value is 10.
	CheckpointInterval *int `json:"checkpoint_interval,omitempty"`
	// OBS path where users are authorized to save the snapshot.
	// This parameter is valid only when checkpoint_enabled is set to true.
	// OBS path where users are authorized to save the snapshot.
	// This parameter is valid only when log_enabled is set to true.
	ObsBucket string `json:"obs_bucket,omitempty"`
	// Whether to enable the function of uploading job logs to users' OBS buckets. The default value is false.
	LogEnabled *bool `json:"log_enabled,omitempty"`
	// SMN topic. If a job fails, the system will send a message to users subscribed to the SMN topic.
	SmnTopic string `json:"smn_topic,omitempty"`
	// Whether to enable the function of automatically restarting a job upon job exceptions. The default value is false.
	RestartWhenException *bool `json:"restart_when_exception,omitempty"`
	// Expiration time, in seconds. The default value is 3600.
	IdleStateRetention *int `json:"idle_state_retention,omitempty"`
	// List of edge computing group IDs. Use commas (,) to separate multiple IDs.
	EdgeGroupIds []string `json:"edge_group_ids,omitempty"`
	// Dirty data policy of a job.
	// 2:obsDir: Save. obsDir specifies the path for storing dirty data.
	// 1: Trigger a job exception
	// 0: Ignore
	// The default value is 0.
	DirtyDataStrategy string `json:"dirty_data_strategy,omitempty"`
	// Name of the resource package that has been uploaded to the DLI resource management system.
	// The UDF Jar file of the SQL job is specified by this parameter.
	UdfJarUrl string `json:"udf_jar_url,omitempty"`
	// Number of CUs in the JobManager selected for a job. The default value is 1.
	ManagerCuNumber *int `json:"manager_cu_number,omitempty"`
	// Number of CUs for each TaskManager. The default value is 1.
	TmCus *int `json:"tm_cus,omitempty"`
	// Number of slots in each TaskManager. The default value is (parallel_number*tm_cus)/(cu_number-manager_cu_number).
	TmSlotNum *int `json:"tm_slot_num,omitempty"`
	// Degree of parallelism (DOP) of an operator.
	OperatorConfig string `json:"operator_config"`
	// Whether the abnormal restart is recovered from the checkpoint.
	ResumeCheckpoint *bool `json:"resume_checkpoint,omitempty"`
	// Maximum number of retry times upon exceptions. The unit is times/hour. Value range: -1 or greater than 0.
	// The default value is -1, indicating that the number of times is unlimited.
	ResumeMaxNum *int `json:"resume_max_num,omitempty"`
	// Traffic or hit ratio of each operator, which is a character string in JSON format.
	StaticEstimatorConfig string `json:"static_estimator_config"`
	// Customizes optimization parameters when a Flink job is running.
	RuntimeConfig string `json:"runtime_config,omitempty"`
	// Flink version. The valid value is `1.1`0 or `1.12`.
	FlinkVersion string `json:"flink_version,omitempty"`
}

type ListOpts struct {
	Name                     string `q:"name"`
	UserName                 string `q:"user_name"`
	QueueName                string `q:"queue_name"`
	Status                   string `q:"status"`
	JobType                  string `q:"job_type"`
	Tags                     string `q:"tags"`
	SysEnterpriseProjectName string `q:"sys_enterprise_project_name"`
	ShowDetail               *bool  `q:"show_detail"`
	Order                    string `q:"order"`
	Offset                   *int   `q:"offset"`
	Limit                    *int   `q:"limit"` // default 10
	//Specifies parent job id of Edge job to query Edge subJob
	// empty: will dont query Edge subJob
	RootJobId *int `q:"root_job_id"`
}

type RunJobOpts struct {
	ResumeSavepoint *bool `json:"resume_savepoint,omitempty"`
	JobIds          []int `json:"job_ids" required:"true"`
}

type ObsBucketsOpts struct {
	Buckets []string `json:"obs_buckets" required:"true"`
}

type CreateJarJobOpts struct {
	// Name of the job. Length range: 0 to 57 characters.
	Name string `json:"name" required:"true"`
	// Job description. Length range: 0 to 512 characters.
	Desc string `json:"desc,omitempty"`
	// Name of a queue. Length range: 1 to 128 characters.
	QueueName string `json:"queue_name,omitempty"`
	// Number of CUs selected for a job.
	CuNumber *int `json:"cu_number,omitempty"`
	// Number of CUs on the management node selected by the user for a job,
	// which corresponds to the number of Flink job managers. The default value is 1.
	ManagerCuNumber int `json:"manager_cu_number,omitempty"`
	// Number of parallel operations selected for a job.
	ParallelNumber int `json:"parallel_number,omitempty"`
	// Whether to enable the job log function.
	// true: indicates to enable the job log function.
	// false: indicates to disable the job log function.
	// Default value: false
	LogEnabled *bool `json:"log_enabled,omitempty"`
	// OBS bucket where users are authorized to save logs when log_enabled is set to true.
	ObsBucket string `json:"obs_bucket,omitempty"`
	// SMN topic. If a job fails, the system will send a message to users subscribed to the SMN topic.
	SmnTopic string `json:"smn_topic,omitempty"`
	// Job entry class.
	MainClass string `json:"main_class,omitempty"`
	// Job entry parameter. Multiple parameters are separated by spaces.
	EntrypointArgs string `json:"entrypoint_args,omitempty"`
	// Whether to enable the function of restart upon exceptions. The default value is false.
	RestartWhenException *bool `json:"restart_when_exception,omitempty"`
	// Name of the package that has been uploaded to the DLI resource management system.
	// This parameter is used to customize the JAR file where the job main class is located.
	Entrypoint string `json:"entrypoint,omitempty"`
	// Name of the package that has been uploaded to the DLI resource management system.
	// This parameter is used to customize other dependency packages.
	// Example: myGroup/test.jar,myGroup/test1.jar.
	DependencyJars []string `json:"dependency_jars,omitempty"`
	// Name of the resource package that has been uploaded to the DLI resource management system.
	// This parameter is used to customize dependency files.
	// Example: myGroup/test.cvs,myGroup/test1.csv.
	// You can add the following content to the application to access the corresponding dependency file:
	// In the command, fileName indicates the name of the file to be accessed,
	// and ClassName indicates the name of the class that needs to access the file.
	// ClassName.class.getClassLoader().getResource("userData/fileName")
	DependencyFiles []string `json:"dependency_files,omitempty"`
	// Number of CUs for each TaskManager. The default value is 1.
	TmCus int `json:"tm_cus,omitempty"`
	// Number of slots in each TaskManager. The default value is (parallel_number*tm_cus)/(cu_number-manager_cu_number).
	TmSlotNum *int `json:"tm_slot_num,omitempty"`
	// Job feature. Type of the Flink image used by a job.
	// basic: indicates that the basic Flink image provided by DLI is used.
	// custom: indicates that the user-defined Flink image is used.
	Feature string `json:"feature,omitempty"`
	// Flink version. This parameter is valid only when feature is set to basic. You can use this parameter with the
	// feature parameter to specify the version of the DLI basic Flink image used for job running.
	FlinkVersion string `json:"flink_version,omitempty"`
	// Custom image. The format is Organization name/Image name:Image version.
	// This parameter is valid only when feature is set to custom. You can use this parameter with the feature
	// parameter to specify a user-defined Flink image for job running. For details about how to use custom images
	Image string `json:"image,omitempty"`
	// Whether the abnormal restart is recovered from the checkpoint.
	ResumeCheckpoint *bool `json:"resume_checkpoint,omitempty"`
	// Maximum number of retry times upon exceptions. The unit is times/hour. Value range: -1 or greater than 0.
	// The default value is -1, indicating that the number of times is unlimited.
	ResumeMaxNum *int `json:"resume_max_num,omitempty"`
	// Storage address of the checkpoint in the JAR file of the user. The path must be unique.
	CheckpointPath string `json:"checkpoint_path,omitempty"`
	// Label of a Flink JAR job. For details, see Table 3.
	Tags []tags.ResourceTag `json:"tags"`
	// Customizes optimization parameters when a Flink job is running.
	RuntimeConfig string `json:"runtime_config,omitempty"`
	// Whether to enable the automatic job snapshot function，default value is **false**.
	// + true: indicates to enable the automatic job snapshot function.
	// + false: indicates to disable the automatic job snapshot function.
	CheckpointEnabled *bool `json:"checkpoint_enabled,omitempty"`
	// The mode of the snapshot.
	// +1: ExactlyOnce, indicates that data is processed only once.
	// +2: AtLeastOnce, indicates that data is processed at least once.
	// The default value is `1`.
	CheckpointMode int `json:"checkpoint_mode,omitempty"`
	// The interval of the snapshot. The unit is second. The default value is **10**.
	CheckpointInterval int `json:"checkpoint_interval,omitempty"`
	// The name of the delegation authorized to DLI. This parameter is supported only when the Flink version is 1.15.
	ExecutionAgencyUrn string `json:"execution_agency_urn,omitempty"`
	// The version of the resource configuration.
	// + v1 (default)
	// + v2
	// The v2 version is not supported to set the CU number, and supports setting Job Manager Memory and Task Manager Memory directly.
	// The v1 version is supported by Flink 1.12, Flink 1.13, and Flink 1.15.
	// The v2 version is supported by Flink 1.13, Flink 1.15, and Flink 1.17.
	// It is recommended to use the v2 version of the parameter settings.
	ResourceConfigVersion string `json:"resource_config_version,omitempty"`
	// The resource configuration of the Flink job.
	// This parameter is valid only when the resource_config_version is set to "v2".
	ResourceConfig *ResourceConfigOpts `json:"resource_config,omitempty"`
}

type ResourceConfigOpts struct {
	// The parameter is used to set the number of parallel tasks that a single TaskManager can provide.
	// Each Task Slot can execute a task in parallel. Increasing Task Slots can improve the parallel processing
	// capability of TaskManager, but also increase resource consumption.
	// The number of Task Slots is related to the CPU number of TaskManager, because each CPU can provide a Task Slot.
	// The default value is 1. The minimum parallel number cannot be less than 1.
	MaxSlot int `json:"max_slot,omitempty"`
	// The parallel number of the job, which is the number of sub-tasks executed in parallel by each operator in the job.
	// The number of sub-tasks of an operator is the parallel degree of the operator. The default value is "1".
	ParallelNumber int `json:"parallel_number,omitempty"`
	// The resource specification of the JobManager.
	JobManagerResourceSpec *JobOrTaskManagerResourceSpecOpts `json:"jobmanager_resource_spec,omitempty"`
	// The resource specification of the TaskManager.
	TaskManagerResourceSpec *JobOrTaskManagerResourceSpecOpts `json:"taskmanager_resource_spec,omitempty"`
}

type JobOrTaskManagerResourceSpecOpts struct {
	// The number of CPU cores that the JobManager or TaskManager can use.
	// The default value is 1.0. The minimum value cannot be less than 0.5.
	CPU float64 `json:"cpu,omitempty"`
	// The memory size that the JobManager or TaskManager can use, in MB or GB (default).
	// The default value is 4GB. The minimum value cannot be less than 2GB.
	Memory string `json:"memory,omitempty"`
}

type TaskManagerResourceSpecOpts struct {
	// The resource specification of the TaskManager.
	ResourceSpec *int `json:"resource_spec,omitempty"`
}

type UpdateJarJobOpts struct {
	// Name of the job. Length range: 0 to 57 characters.
	Name string `json:"name,omitempty"`
	// Job description. Length range: 0 to 512 characters.
	Desc string `json:"desc,omitempty"`
	// Name of a queue. Length range: 1 to 128 characters.
	QueueName string `json:"queue_name,omitempty"`
	// Number of CUs selected for a job. The default value is 2.
	CuNumber *int `json:"cu_number,omitempty"`
	// Number of CUs on the management node selected by the user for a job, which corresponds to the number of Flink
	// job managers. The default value is 1.
	ManagerCuNumber int `json:"manager_cu_number,omitempty"`
	// Number of parallel operations selected for a job. The default value is 1.
	ParallelNumber int `json:"parallel_number,omitempty"`
	// Whether to enable the job log function.
	// true: indicates to enable the job log function.
	// false: indicates to disable the job log function.
	// Default value: false
	LogEnabled *bool `json:"log_enabled,omitempty"`
	// OBS path where users are authorized to save logs when log_enabled is set to true.
	ObsBucket string `json:"obs_bucket,omitempty"`
	// SMN topic. If a job fails, the system will send a message to users subscribed to the SMN topic.
	SmnTopic string `json:"smn_topic,omitempty"`
	// Job entry class.
	MainClass string `json:"main_class,omitempty"`
	// Job entry parameter. Multiple parameters are separated by spaces.
	EntrypointArgs string `json:"entrypoint_args,omitempty"`
	// Whether to enable the function of restart upon exceptions. The default value is false.
	RestartWhenException *bool `json:"restart_when_exception,omitempty"`
	// Name of the package that has been uploaded to the DLI resource management system. This parameter is used to
	// customize the JAR file where the job main class is located.
	Entrypoint string `json:"entrypoint,omitempty"`
	// Name of the package that has been uploaded to the DLI resource management system. This parameter is used to
	// customize other dependency packages.
	// Example: myGroup/test.jar,myGroup/test1.jar.
	DependencyJars []string `json:"dependency_jars,omitempty"`
	// Name of the resource package that has been uploaded to the DLI resource management system. This parameter is
	// used to customize dependency files.
	// Example: myGroup/test.cvs,myGroup/test1.csv.
	DependencyFiles []string `json:"dependency_files,omitempty"`
	// Number of CUs for each TaskManager. The default value is 1.
	TmCus int `json:"tm_cus,omitempty"`
	// Number of slots in each TaskManager. The default value is (parallel_number*tm_cus)/(cu_number-manager_cu_number).
	TmSlotNum *int `json:"tm_slot_num,omitempty"`
	// Job feature. Type of the Flink image used by a job.
	// basic: indicates that the basic Flink image provided by DLI is used.
	// custom: indicates that the user-defined Flink image is used.
	Feature string `json:"feature,omitempty"`
	// Flink version. This parameter is valid only when feature is set to basic. You can use this parameter with the
	// feature parameter to specify the version of the DLI basic Flink image used for job running.
	FlinkVersion string `json:"flink_version,omitempty"`
	// Custom image. The format is Organization name/Image name:Image version.
	// This parameter is valid only when feature is set to custom. You can use this parameter with the feature
	// parameter to specify a user-defined Flink image for job running. For details about how to use custom images.
	Image string `json:"image,omitempty"`
	// Whether the abnormal restart is recovered from the checkpoint.
	ResumeCheckpoint *bool `json:"resume_checkpoint,omitempty"`
	// Maximum number of retry times upon exceptions. The unit is times/hour. Value range: -1 or greater than 0.
	// The default value is -1, indicating that the number of times is unlimited.
	ResumeMaxNum *int `json:"resume_max_num,omitempty"`
	// Storage address of the checkpoint in the JAR file of the user. The path must be unique.
	CheckpointPath string `json:"checkpoint_path,omitempty"`
	// Customizes optimization parameters when a Flink job is running.
	RuntimeConfig string `json:"runtime_config,omitempty"`
	// Whether to enable the automatic job snapshot function，default value is **false**.
	// + true: indicates to enable the automatic job snapshot function.
	// + false: indicates to disable the automatic job snapshot function.
	CheckpointEnabled *bool `json:"checkpoint_enabled,omitempty"`
	// The mode of the snapshot.
	// +1: ExactlyOnce, indicates that data is processed only once.
	// +2: AtLeastOnce, indicates that data is processed at least once.
	// The default value is `1`.
	CheckpointMode int `json:"checkpoint_mode,omitempty"`
	// The interval of the snapshot. The unit is second. The default value is **10**.
	CheckpointInterval int `json:"checkpoint_interval,omitempty"`
	// The name of the delegation authorized to DLI. This parameter is supported only when the Flink version is 1.15.
	ExecutionAgencyUrn string `json:"execution_agency_urn,omitempty"`
	// The version of the resource configuration.
	// + v1 (default)
	// + v2
	// The v2 version is not supported to set the CU number, and supports setting Job Manager Memory and Task Manager Memory directly.
	// The v1 version is supported by Flink 1.12, Flink 1.13, and Flink 1.15.
	// The v2 version is supported by Flink 1.13, Flink 1.15, and Flink 1.17.
	// It is recommended to use the v2 version of the parameter settings.
	ResourceConfigVersion string `json:"resource_config_version,omitempty"`
	// The resource configuration of the Flink job.
	// This parameter is valid only when the resource_config_version is set to "v2".
	ResourceConfig *ResourceConfigOpts `json:"resource_config,omitempty"`
}

type StopFlinkJobInBatch struct {
	TriggerSavepoint *bool `json:"trigger_savepoint,omitempty"`
	JobIds           []int `json:"job_ids" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func CreateSqlJob(c *golangsdk.ServiceClient, opts CreateSqlJobOpts) (*CreateJobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CreateJobResp
	_, err = c.Post(createFlinkSqlUrl(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func UpdateSqlJob(c *golangsdk.ServiceClient, jobId int, opts UpdateSqlJobOpts) (*UpdateJobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst UpdateJobResp
	_, err = c.Put(updateFlinkSqlURL(c, jobId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Run(c *golangsdk.ServiceClient, opts RunJobOpts) (*[]CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst []CommonResp
	_, err = c.Post(runFlinkJobURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Get(c *golangsdk.ServiceClient, jobId int) (*GetJobResp, error) {
	var rst GetJobResp
	_, err := c.Get(getURL(c, jobId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func List(c *golangsdk.ServiceClient, opts ListOpts) (*ListResp, error) {
	url := listURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst ListResp
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Delete(c *golangsdk.ServiceClient, jobId int) (*CommonResp, error) {
	var rst CommonResp
	_, err := c.DeleteWithResponse(deleteURL(c, jobId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func AuthorizeBucket(c *golangsdk.ServiceClient, opts ObsBucketsOpts) (*CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CommonResp
	_, err = c.Post(authorizeBucketURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func CreateJarJob(c *golangsdk.ServiceClient, opts CreateJarJobOpts) (*CreateJobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CreateJobResp
	_, err = c.Post(createJarJobURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func UpdateJarJob(c *golangsdk.ServiceClient, jobId int, opts UpdateJarJobOpts) (*UpdateJobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst UpdateJobResp
	_, err = c.Put(updateJarJobURL(c, jobId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Stop(c *golangsdk.ServiceClient, opts StopFlinkJobInBatch) (*[]CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst []CommonResp
	_, err = c.Post(stopJobURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}
