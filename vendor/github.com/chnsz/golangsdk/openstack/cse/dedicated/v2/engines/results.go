package engines

// RequestResp is the structure that represents the response of the API request.
type RequestResp struct {
	// The ID of the dedicated microservice engine.
	ID string `json:"id"`
	// The name of the dedicated microservice engine.
	Name string `json:"name"`
	// The Job ID.
	JobId int `json:"jobId"`
}

// Engine is the structure that represents the details of the Microservice engine.
type Engine struct {
	// The ID of the dedicated microservice engine.
	ID string `json:"id"`
	// The name of the dedicated microservice engine.
	Name string `json:"name"`
	// The project ID to which the dedicated microservice engine belongs.
	ProjectId string `json:"projectId"`
	// The enterprise project ID to which the dedicated microservice engine belongs.
	EnterpriseProjectId string `json:"enterpriseProjectId"`
	// The enterprise project name to which the dedicated microservice engine belongs.
	EnterpriseProjectName string `json:"enterpriseProjectName"`
	// The engine type.
	Type string `json:"type"`
	// The description of the dedicated microservice engine.
	Description string `json:"description"`
	// Whether the engine is the default engine.
	IsDefault bool `json:"beDefault"`
	// The flavor of the dedicated microservice engine.
	//   cse.s1.small2: High availability 100 instance engine.
	//   cse.s1.medium2: High availability 200 instance engine.
	//   cse.s1.large2: High availability 500 instance engine.
	//   cse.s1.xlarge2: High availability 2000 instance engine.
	Flavor string `json:"flavor"`
	// The charging mode of the dedicated microservice engine.
	//   0: pre-paid.
	//   1: post-paid.
	//   2: free.
	Payment string `json:"payment"`
	// The authentication method for the dedicated microservice engine.
	// The "RBAC" is security authentication, and the "NONE" is no authentication.
	AuthType string `json:"authType"`
	// The current status.
	Status string `json:"status"`
	// The CCE flavor of the dedicated microservice engine.
	CceSpec CceSpec `json:"cceSpec"`
	// The access address of the dedicated microservice engine.
	ExternalEntrypoint ExternalEntrypoint `json:"externalEntrypoint"`
	// The IP address of the public access for the dedicated microservice engine.
	PublicAddress string `json:"publicAddress"`
	// The current version of the dedicated microservice engine.
	Version string `json:"version"`
	// The latest version of the dedicated microservice engine.
	LatestVersion string `json:"latestVersion"`
	// The creation time of the dedicated microservice engine.
	CreateTime int `json:"createTime"`
	// The creator of the dedicated microservice engine.
	CreateUser string `json:"createUser"`
	// The latest task ID of the dedicated microservice engine.
	LatestJobId int `json:"latestJobId"`
	// Additional operations allowed by the dedicated microservice engine.
	//   Delete
	//   ForceDelete
	//   Upgrade
	//   Modify
	//   Retry
	EngineAdditionalActions []string `json:"engineAdditionalActions"`
	// The deployment type of the dedicated microservice engine. The fixed value is "CSE2".
	SpecType string `json:"specType"`
	// Additional information for the dedicated microservice engine.
	Reference Reference `json:"reference"`
	//The list of virtual machine IDs used by the current microservice engine on the resource tenant side.
	VmIds []string `json:"vmIds"`
}

// CceSpec is an object that represents the configuration of the associated CCE cluster.
type CceSpec struct {
	// The CCE flavor ID.
	ID int `json:"id"`
	// The ID of the dedicated microservice engine.
	EngineId string `json:"engineId"`
	// The deployment type of the dedicated microservice engine.
	SpecType string `json:"specType"`
	// The cluster ID of the dedicated microservice engine.
	ClusterId string `json:"clusterId"`
	// The list of the CCE nodes for the dedicated microservice engine.
	ClusterNodes NodeDetails `json:"clusterNodes"`
	// The CCE cluster flavor.
	Flavor string `json:"flavor"`
	// The region where the CCE cluster is located.
	Region string `json:"region"`
	// The CCE cluster version.
	Version string `json:"version"`
	// Additional parameters for the CCE cluster.
	ExtendParam string `json:"extendParam"`
}

// NodeDetails is an object that represents the nodes list.
type NodeDetails struct {
	// The list of the CCE nodes.
	Nodes []Node `json:"clusterNodes"`
}

// Node is an object that represents the details of the CCE node.
type Node struct {
	// The node ID.
	ID string `json:"id"`
	// The availability zone where the node is located.
	AvailabilityZone string `json:"az"`
	// The node IP.
	IP string `json:"ip"`
	// The node tag.
	Label string `json:"label"`
	// The node status.
	Status string `json:"status"`
}

// ExternalEntrypoint is an object that represents the access information.
type ExternalEntrypoint struct {
	// The access address in the VPC on the tenant side of the dedicated microservice engine.
	ExternalAddress string `json:"externalAddress"`
	// The public network access address of the dedicated microservice engine, it needs to be enabled for public access.
	PublicAddress string `json:"publicAddress"`
	// Access address in the VPC on the tenant side of the component of the dedicated Microservice Engine.
	ServiceEndpoint ServiceEndpoint `json:"serviceEndpoint"`
	// The public network access address of the component of the dedicated Microservice Engine.
	// The public network access needs to be enabled.
	PublicServiceEndpoint ServiceEndpoint `json:"publicServiceEndpoint"`
}

// ServiceEndpoint is an object that represent the entrypoints of the service center and config center.
type ServiceEndpoint struct {
	// The entrypoint details of the service center.
	ServiceCenter Detail `json:"serviceCenter"`
	// The entrypoint details of the config center.
	ConfigCenter Detail `json:"kie"`
}

// Detail is an object that represent the endpoint informations of the service center or config center.
type Detail struct {
	// The main ipv4 access address in the VPC of the Microservice Engine Exclusive Edition component.
	MasterEntrypoint string `json:"masterEntrypoint"`
	// The main ipv6 access address in the VPC of the Microservice Engine Exclusive Edition component.
	MasterEntrypointIpv6 string `json:"masterEntrypointIpv6"`
	// The ipv4 standby access address in the VPC of the Microservice Engine Exclusive Edition component.
	SlaveEntrypoint string `json:"slaveEntrypoint"`
	// The ipv6 standby access address in the VPC of the Microservice Engine Exclusive Edition component.
	SlaveEntrypointIpv6 string `json:"slaveEntrypointIpv6"`
	// The component type of the dedicated microservice engine.
	Type string `json:"type"`
}

// Reference is an object that represent the additional parameters of the dedicated microservice engine.
type Reference struct {
	// The VPC name.
	Vpc string `json:"vpc"`
	// The VPC ID.
	VpcId string `json:"vpcId"`
	// List of deployment availability zones for the dedicated Microservice Engine.
	AzList []string `json:"azList"`
	// The subnet network ID of the dedicated microservice engine.
	NetworkId string `json:"networkId"`
	// The ipv4 subnet division.
	SubnetCidr string `json:"subnetCidr"`
	// The ipv6 subnet division.
	SubnetCidV6 string `json:"subnetCidV6"`
	// The subnet gateway.
	SubnetGateway string `json:"subnetGateway"`
	// The public network address ID of the dedicated microservice engine. Public network access needs to be enabled.
	PublicIpId string `json:"publicIpId"`
	// The total number of microservices that the plan can support.
	ServiceLimit string `json:"serviceLimit"`
	// The total number of microservice instances that the plan can support.
	InstanceLimit string `json:"instanceLimit"`
	// Additional parameters for the dedicated microservice engine.
	Inputs map[string]interface{} `json:"inputs"`
}

// Job is an object that represent the details of the excuting job.
type Job struct {
	// The job ID.
	Id int `json:"id"`
	// The engine ID corresponding to the currently executing job.
	EngineId string `json:"engineId"`
	// The job type.
	//   Create: create engine.
	//   Delete: delete engine.
	//   Upgrade: upgrade engine.
	//   Modify：update engine flavor.
	//   Configure：update engine configuration.
	Type string `json:"type"`
	// The job description.
	Description string `json:"description"`
	// The current execution status of the job.
	//   Init：initialization.
	//   Executing：executing.
	//   Error：failed to execute.
	//   Timeout：timeout to execute。
	//   Finished：execution complete。
	Status string `json:"status"`
	// Whether the job is executing, 0 means not executing, 1 means executing.
	Scheduling int `json:"scheduling"`
	// The job creator.
	CreateUser string `json:"createUser"`
	// The time when the job started executing.
	StartTime int `json:"startTime"`
	// The time when the job ends.
	EndTime int `json:"endTime"`
	// The job execution context.
	Context string `json:"context"`
	// The processing stages of the job.
	Tasks []Task `json:"tasks"`
}

// Task is an object that represent the details of the processing stages of the job.
type Task struct {
	// The current task name.
	Name string `json:"taskName"`
	// List of processing step names contained in the current processing stage.
	Steps []string `json:"taskNames"`
	// The current status.
	Status string `json:"status"`
	// Task processing phase start time.
	StartTime int `json:"startTime"`
	// Task processing phase end time.
	EndTime int `json:"endTime"`
	// The task metadata.
	ExecutorBrief ExecutorBrief `json:"taskExecutorBrief"`
	// Subtasks included in the processing phase.
	SubTasks []SubTask `json:"tasks"`
}

// ExecutorBrief is an object that represent the job metadata.
type ExecutorBrief struct {
	// The subtask duration.
	Duration int `json:"duration"`
	// The subtask description.
	Description string `json:"description"`
}

// SubTask is an object that represent the details of the subjob.
type SubTask struct {
	// The task ID to which the subtask belongs.
	JobId int `json:"jobId"`
	// The subtask ID.
	Id int `json:"id"`
	// The subtask type.
	Type string `json:"type"`
	// The executor of the subtask
	Assigned string `json:"assigned"`
	// The subtask name.
	TaskName string `json:"taskName"`
	// The name of the engine to which the subtask belongs.
	EngineName string `json:"engineName"`
	// The order in which the subtasks are executed, from small to large.
	TaskOrder int `json:"taskOrder"`
	// The subtask status
	Status string `json:"status"`
	// The subtask start time.
	StartTime int `json:"startTime"`
	// The subtask end time.
	EndTime int `json:"endTime"`
	// The subtask creation time.
	CreateTime int `json:"createTime"`
	// The subtask update time.
	UpdateTime int `json:"updateTime"`
	// Whether subtask is timeout.
	Timeout int `json:"timeout"`
	// The subtask details, auxiliary information generated during execution.
	Log string `json:"log"`
	// The subtask output information.
	Output string `json:"output"`
	// The subtask metadata.
	ExecutorBrief ExecutorBrief `json:"taskExecutorBrief"`
}

type ErrorResponse struct {
	// Error code
	ErrCode string `json:"error_code"`
	// Error message
	ErrMessage string `json:"error_message"`
}
