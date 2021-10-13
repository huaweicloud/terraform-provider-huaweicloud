package batches

const (
	// StateStarting is a state means the batch processing job is being started.
	StateStarting = "starting"
	// StateRunning is a state means the batch processing job is executing a task.
	StateRunning = "running"
	// StateDead is a state means the batch processing job has failed to execute.
	StateDead = "dead"
	// StateSuccess is a state means the batch processing job is successfully executed.
	StateSuccess = "success"
	// StateRecovering is a state means the batch processing job is being restored.
	StateRecovering = "recovering"
)

// CreateResp represents a result of the Create method.
type CreateResp struct {
	// ID of a batch processing job.
	ID string `json:"id"`
	// Back-end application ID of a batch processing job.
	AppId []string `json:"appId"`
	// Name of a batch processing job.
	Name string `json:"name"`
	// Owner of a batch processing job.
	Owner string `json:"owner"`
	// Proxy user (resource tenant) to which a batch processing job belongs.
	ProxyUser string `json:"proxyUser"`
	// Status of a batch processing job. For details, see Table 7 in Creating a Batch Processing Job.
	State string `json:"state"`
	// Type of a batch processing job. Only Spark parameters are supported.
	Kind string `json:"kind"`
	// Last 10 records of the current batch processing job.
	Log []string `json:"log"`
	// Type of a computing resource. If the computing resource type is customized, value CUSTOMIZED is returned.
	ScType string `json:"sc_type"`
	// Queue where a batch processing job is located.
	ClusterName string `json:"cluster_name"`
	// Queue where a batch processing job is located.
	Queue string `json:"queue"`
	// Time when a batch processing job is created. The timestamp is expressed in milliseconds.
	CreateTime int `json:"create_time"`
	// Time when a batch processing job is updated. The timestamp is expressed in milliseconds.
	UpdateTime int `json:"update_time"`
	// Job feature. Type of the Spark image used by a job.
	// basic: indicates that the basic Spark image provided by DLI is used.
	// custom: indicates that the user-defined Spark image is used.
	// ai: indicates that the AI image provided by DLI is used.
	Feature string `json:"feature"`
	// Version of the Spark component used by a job. Set this parameter when feature is set to basic or ai.
	// If this parameter is not set, the default Spark component version 2.3.2 is used.
	SparkVersion string `json:"spark_version"`
	// Custom image. The format is Organization name/Image name:Image version.
	// This parameter is valid only when feature is set to custom. You can use this parameter with the feature parameter
	// to specify a user-defined Spark image for job running. For details about how to use custom images, see the Data
	// Lake Insight User Guide.
	Image string `json:"image"`
}

// StateResp represents a result of the GetState method.
type StateResp struct {
	// ID of a batch processing job, which is in the universal unique identifier (UUID) format.
	ID string `json:"id"`
	// Status of a batch processing job.
	// The valid values are starting, running, dead, success and recovering. For detail, see Constant definition.
	State string `json:"state"`
}
