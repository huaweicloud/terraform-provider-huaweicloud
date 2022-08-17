package instances

import "github.com/chnsz/golangsdk/pagination"

type JobResp struct {
	// Component instance ID.
	InstanceId string `json:"instance_id"`
	// Job ID.
	JobId string `json:"job_id"`
}

// Instance is the structure that represents the details of the ServiceStage component instance.
type Instance struct {
	// Component instance ID.
	ID string `json:"id"`
	// Component instance name.
	Name string `json:"name"`
	// Component environment ID.
	EnvironmentId string `json:"environment_id"`
	// Platform type.
	// Value: cce or vmapp.
	PlatformType string `json:"platform_type"`
	// Instance description.
	Description string `json:"description"`
	// Resource flavor.
	FlavorId string `json:"flavor_id"`
	// Artifact. key indicates the component name. In the Docker container scenario, key indicates the container name.
	Artifacts map[string]Artifact `json:"artifacts"`
	// Component version.
	Version string `json:"version"`
	// Component configurations, such as environment variables.
	Configuration ConfigurationResp `json:"configuration"`
	// Creator.
	Creator string `json:"creator"`
	// Creation time.
	CreatedAt int `json:"create_time"`
	// Update time.
	UpdatedAt int `json:"update_time"`
	// Access mode.
	ExternalAccesses []ExternalAccessResp `json:"external_accesses"`
	// Deployed resources.
	ReferResources []ReferResource `json:"refer_resources"`
	// Status details.
	StatusDetail StatusDetail `json:"status_detail"`
}

// ConfigurationResp is an object represents the configuration details of the deployment.
type ConfigurationResp struct {
	// Environment variable.
	EnvVariables []VariableResp `json:"env"`
	// Data storage configuration.
	Storages []StorageResp `json:"storage"`
	// Upgrade policy.
	Strategy StrategyResp `json:"strategy"`
	// Lifecycle.
	Lifecycle LifecycleResp `json:"lifecycle"`
	// Policy list of log collection.
	LogCollectionPolicy []LogCollectionPolicyResp `json:"logs"`
	// Scheduling policy.
	Scheduler SchedulerResp `json:"scheduler"`
	// Health check.
	Probe ProbeResp `json:"probes"`
}

// VariableResp is an object represents the detail of the environment variable.

type VariableResp struct {
	// Whether variable is internal variable.
	Internal bool `json:"internal"`
	// Environment variable name.
	// The value contains 1 to 64 characters, including letters, digits, underscores (_), hyphens (-), and dots (.),
	// and cannot start with a digit.
	Name string `json:"name"`
	// Environment variable value.
	Value string `json:"value"`
}

// StorageResp is an object represents the storage configuration of the deployment.
type StorageResp struct {
	// Storage type. Value:
	// HostPath: host path mounting.
	// EmptyDir: temporary directory mounting.
	// ConfigMap: configuration item mounting.
	// Secret: secret volume mounting.
	// PersistentVolumeClaim: cloud storage mounting.
	Type string `json:"type" required:"true"`
	// Storage parameter.
	Parameters StorageParams `json:"parameters" required:"true"`
	// Directory mounted to the container.
	Mounts []Mount `json:"mounts" required:"true"`
}

// StrategyResp is an object represents the upgrade information of the deployment.
type StrategyResp struct {
	Spec Spec `json:"spec"`
	// Upgrade policy. Value: Recreate or RollingUpdate (default).
	// The former indicates in-place upgrade while the latter indicates rolling upgrade.
	Upgrade string `json:"upgrade,omitempty"`
}

// Spec is an object represents the other upgrade details.
type Spec struct {
	// The maximum surge number.
	MaxSurge int `json:"maxSurge"`
	// The maximum unavailable upgrade count.
	MaxUnavailable int `json:"maxUnavailable"`
}

// LifecycleResp is an object represents the lifecycle of the deployment.
type LifecycleResp struct {
	// Startup command.
	Entrypoint Entrypoint `json:"entrypoint"`
	// Post-start processing.
	PostStart ProcessResp `json:"post-start"`
	// Pre-stop processing.
	PreStop ProcessResp `json:"pre-stop"`
}

// ProcessResp is an object represents the details of the post-processing or stop pre-processing.
type ProcessResp struct {
	// Process type. The value is command or http.
	// The command is to execute the command line, and http is to send an http request.
	Type string `json:"type" required:"true"`
	// Start post-processing or stop pre-processing parameters.
	Parameters ProcessParams `json:"parameters" required:"true"`
}

// LogCollectionPolicyResp is an object represents the details of the log collection policy.
type LogCollectionPolicyResp struct {
	// Container mounting path.
	LogPath string `json:"logPath"`
	// The extend host path, the valid values are as follows:
	//	None
	//	PodUID
	//	PodName
	//	PodUID/ContainerName
	//	PodName/ContainerName
	// If omited, means container mounting.
	HostExtendPath string `json:"hostExtendPath"`
	// Aging period.
	AgingPeriod string `json:"rotate"`
	// Host mounting path.
	HostPath string `json:"hostPath"`
}

// SchedulerResp is an object represents the scheduling policy.
type SchedulerResp struct {
	// Affinity.
	Affinity Affinity `json:"affinity,omitempty"`
	// Anti-affinity.
	AntiAffinity Affinity `json:"anti-affinity,omitempty"`
}

// ProbeResp is an object represents the probe members configuration of the deployment.
type ProbeResp struct {
	// Component liveness probe.
	LivenessProbe ProbeDetail `json:"livenessProbe,omitempty"`
	// Component service probe.
	ReadinessProbe ProbeDetail `json:"readinessProbe,omitempty"`
}

// ExternalAccessResp is an object represents the configuration of the external IP access.
type ExternalAccessResp struct {
	// Protocol. Value: http or https.
	Protocol string `json:"protocol"`
	// Access address.
	Address string `json:"address"`
	// Port number.
	ForwardPort int `json:"forward_port"`
	// Type.
	Type string `json:"type"`
	// Status.
	Status string `json:"status"`
	// Creation time.
	CreatedAt int `json:"create_time"`
	// Update time.
	UpdatedAt int `json:"update_time"`
}

// StorageResp is an object represents the status details of the deployment.
type StatusDetail struct {
	// Enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// Instance status.
	Status string `json:"status"`
	// Number of normal instance replicas.
	AvailableReplica int `json:"available_replica"`
	// Number of instance replicas.
	Replica int `json:"replica"`
	// Failure description.
	FailDetail string `json:"fail_detail"`
	// Latest job ID.
	LastJobId string `json:"last_job_id"`
	// Latest job status.
	LastJobStatus string `json:"last_job_status"`
}

// InstancePage is a single page maximum result representing a query by offset page.
type InstancePage struct {
	pagination.OffsetPageBase
}

// ExtractInstances is a method to extract the list of environment details for ServiceStage service.
func ExtractInstances(r pagination.Page) ([]Instance, error) {
	var s []Instance
	r.(InstancePage).Result.ExtractIntoSlicePtr(&s, "isntances")
	return s, nil
}

// Job is the structure that represents the result of the Action.
type Job struct {
	// Job ID.
	ID string `json:"job_id"`
}
