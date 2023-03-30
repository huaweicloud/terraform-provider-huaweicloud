package components

import "github.com/chnsz/golangsdk/pagination"

// Component is the structure that represents the detail of the application component.
type Component struct {
	// Source of the code or software package.
	Source Source `json:"source"`
	// Application component ID.
	ID string `json:"id"`
	// Application component name.
	Name         string `json:"name"`
	WorkloadName string `json:"workload_name"`
	Description  string `json:"description"`
	// Component Labels
	Labels    []KeyValue `json:"labels"`
	PodLabels []KeyValue `json:"pod_labels"`
	Build     Build      `json:"build"`
	// Component RuntimeStack
	RuntimeStack RuntimeStack `json:"runtime_stack"`
	// Component external accessed
	ExternalAccesses []ExternalAccess `json:"external_accesses"`
	// Component Status
	Status Status `json:"status"`
	// Environment info
	EnvironmentName string `json:"environment_name"`
	EnvironmentID   string `json:"environment_id"`
	// Application info
	ApplicationName string `json:"application_name"`
	ApplicationID   string `json:"application_id"`
	// Creator
	Creator string `json:"creator"`
	// PlatformType, enum: cce or cci
	PlatformType string `json:"platform_type"`
	// Version
	Version         string          `json:"version"`
	LimitCpu        int             `json:"limit_cpu"`
	LimitMemory     int             `json:"limit_memory"`
	RequestCpu      int             `json:"request_cpu"`
	RequestMemory   int             `json:"request_memory"`
	Replica         int             `json:"replica"`
	Envs            []Env           `json:"envs"`
	Storages        []Storage       `json:"storages"`
	Command         Command         `json:"command"`
	PostStart       K8sLifeCycle    `json:"post_start"`
	PreStop         K8sLifeCycle    `json:"pre_stop"`
	Timezone        string          `json:"timezone"`
	Mesher          Mesher          `json:"mesher"`
	DeployStrategy  DeployStrategy  `json:"deploy_strategy"`
	HostAliases     []HostAlias     `json:"host_aliases"`
	DnsPolicy       string          `json:"dns_policy"`
	DnsConfig       DnsConfig       `json:"dns_config"`
	SecurityContext SecurityContext `json:"security_context"`
	WorkloadKind    string          `json:"workload_kind"`
	JvmOpts         string          `json:"jvm_opts"`
	TomcatOpts      TomcatOpts      `json:"tomcat_opts"`
	Logs            []Log           `json:"logs"`
	CustomMetric    CustomMetric    `json:"custom_metric"`
	Affinity        Affinity        `json:"affinity"`
	AntiAffinity    Affinity        `json:"anti_affinity"`
	LivenessProbe   K8sProbe        `json:"liveness_probe"`
	ReadinessProbe  K8sProbe        `json:"readiness_probe"`
	ReferResources  []*Resource     `json:"refer_resources"`
	// The enterprise project ID.
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RuntimeStack struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Type       string `json:"type"`
	DeployMode string `json:"deploy_mode"`
}

type Status struct {
	// Component status
	ComponentStatus  string `json:"component_status"`
	AvailableReplica int    `json:"available_replica"`
	Replica          int    `json:"replica"`
	FailDetail       string `json:"fail_detail"`
	LastJobId        string `json:"last_job_id"`
	// Creator.
	Creator string `json:"creator"`
	// Creation time.
	CreatedAt int `json:"create_time"`
	// Update time.
	UpdatedAt int `json:"update_time"`
}

type ExternalAccess struct {
	// Protocol Enum: http or https
	Protocol string `json:"protocol"`
	// Address Max len is 256
	Address string `json:"address"`
	// ForwardPort
	ForwardPort int `json:"forward_port"`
}

// Source is an object to specified the source information of Open-Scoure codes or package storage.
type Source struct {
	// Type. Option: source code or artifact software package.
	Kind string `json:"kind" required:"true"`
	// Address of the software package or source code.
	Url     string `json:"url,omitempty"`
	Version string `json:"version"`
	// Storage mode. Value: swr or obs.
	Storage string `json:"storage,omitempty"`

	// The parameters of code are as follows:
	// Code repository. Value: GitHub, GitLab, Gitee, or Bitbucket.
	RepoType string `json:"repo_type,omitempty"`
	// Code repository URL. Example: https://github.com/example/demo.git.
	RepoUrl string `json:"repo_url,omitempty"`
	// Authorization name, which can be obtained from the authorization list.
	RepoAuth string `json:"repo_auth,omitempty"`
	// The code's organization. Value: GitHub, GitLab, Gitee, or Bitbucket.
	RepoNamespace string `json:"repo_namespace,omitempty"`
	// Code branch or tag. Default value: master.
	RepoRef string `json:"repo_ref,omitempty"`
	WebUrl  string `json:"web_url"`

	// Type. Value: package.
	Type string `json:"type,omitempty"`

	// Authentication mode. Value: iam or none. Default value: iam.
	Auth string `json:"auth,omitempty"`

	CodeartsProjectId string `json:"codearts_project_id"`
}

// Build is the component builder, the configuration details refer to 'Parameter'.
type Build struct {
	// This parameter is provided only when no ID is available during build creation.
	Parameter Parameter `json:"parameters" required:"true"`
}

// Parameter is an object to specified the building configuration of codes or package.
type Parameter struct {
	// Compilation command. By default:
	// When build.sh exists in the root directory, the command is ./build.sh.
	// When build.sh does not exist in the root directory, the command varies depending on the operating system (OS). Example:
	// Java and Tomcat: mvn clean package
	// Nodejs: npm build
	BuildCmd string `json:"build_cmd,omitempty"`
	// Address of the Docker file. By default, the Docker file is in the root directory (./).
	DockerfilePath string `json:"dockerfile_path,omitempty"`
	// Build archive organization. Default value: cas_{project_id}.
	ArtifactNamespace string `json:"artifact_namespace,omitempty"`
	// The ID of the cluster to be built.
	ClusterId     string `json:"cluster_id,omitempty"`
	EnvironmentId string `json:"environment_id"`
	// key indicates the key of the tag, and value indicates the value of the tag.
	NodeLabelSelector map[string]interface{} `json:"node_label_selector,omitempty"`
}

type Env struct {
	Name         string       `json:"name"`
	Value        string       `json:"value"`
	Inner        bool         `json:"inner"`
	EnvValueFrom EnvValueFrom `json:"value_from"`
}

type EnvValueFrom struct {
	ReferenceType string `json:"reference_type"`
	Name          string `json:"name"`
	Key           string `json:"key"`
	Optional      bool   `json:"optional"`
}

type Storage struct {
	Type       string           `json:"type"`
	Name       string           `json:"name"`
	Parameters StorageParameter `json:"parameters"`
	Mounts     StorageMounts    `json:"mounts"`
}

type StorageParameter struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	DefaultMode int    `json:"default_mode"`
	Medium      string `json:"medium"`
}

type StorageMounts struct {
	Path     string `json:"path"`
	SubPath  string `json:"sub_path"`
	Readonly bool   `json:"readonly"`
}

type Command struct {
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

type K8sBase struct {
	// Enum: http or https
	Type    string   `json:"type"`
	Scheme  string   `json:"scheme"`
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	Path    string   `json:"path"`
	Command []string `json:"command"`
}

type K8sLifeCycle struct {
	K8sBase
}

type K8sProbe struct {
	K8sBase
	Delay   int `json:"delay"`
	Timeout int `json:"timeout"`
}

type Mesher struct {
	Port int `json:"port"`
}

type DeployStrategy struct {
	Type           string         `json:"type"`
	RollingRelease RollingRelease `json:"rolling_release"`
}

type RollingRelease struct {
	Batches            int    `json:"batches"`
	TerminationSeconds int    `json:"termination_seconds"`
	FailStrategy       string `json:"fail_strategy"`
}

type GrayRelease struct {
	Type              string            `json:"type"`
	FirstBatchWeight  int               `json:"first_batch_weight"`
	FirstBatchReplica int               `json:"first_batch_replica"`
	RemainingBatch    int               `json:"remaining_batch"`
	DeploymentMode    int               `json:"deployment_mode"`
	ReplicaSurgeMode  string            `json:"replica_surge_mode"`
	RuleMatchMode     string            `json:"rule_match_mode"`
	Rules             []GrayReleaseRule `json:"rules"`
}

type GrayReleaseRule struct {
	Type      string `json:"type"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Condition string `json:"condition"`
}

type HostAlias struct {
	IP        string   `json:"ip"`
	HostNames []string `json:"host_names"`
}

type DnsConfig struct {
	Nameservers []string    `json:"nameservers"`
	Searches    []string    `json:"searches"`
	Options     []NameValue `json:"options"`
}

type SecurityContext struct {
	RunAsUser    int          `json:"run_as_user"`
	RunAsGroup   int          `json:"run_as_group"`
	Capabilities Capabilities `json:"capabilities"`
}

type Capabilities struct {
	Add  []string `json:"add"`
	Drop []string `json:"drop"`
}

type TomcatOpts struct {
	ServerXml string `json:"server_xml"`
}

type Log struct {
	LogPath        string `json:"log_path"`
	Rotate         string `json:"rotate"`
	HostPath       string `json:"host_path"`
	HostExtendPath string `json:"host_extend_path"`
}

type CustomMetric struct {
	Path       string `json:"path"`
	Port       int    `json:"port"`
	Dimensions string `json:"dimensions"`
}

type Affinity struct {
	Kind             string            `json:"kind"`
	Condition        string            `json:"condition"`
	Weight           int               `json:"weight"`
	MatchExpressions []MatchExpression `json:"match_expressions"`
}

type Resource struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Parameters struct {
		NameSpace string `json:"name_space,omitempty"`
	} `json:"parameters"`
}

type MatchExpression struct {
	KeyValue
	Operation string `json:"operation"`
}

// ComponentPage is a single page maximum result representing a query by offset page.
type ComponentPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ComponentPage struct is empty.
func (b ComponentPage) IsEmpty() (bool, error) {
	arr, err := ExtractComponents(b)
	return len(arr) == 0, err
}

// ExtractComponents is a method to extract the list of component details for ServiceStage service.
func ExtractComponents(r pagination.Page) ([]Component, error) {
	var s []Component
	err := r.(ComponentPage).Result.ExtractIntoSlicePtr(&s, "components")
	return s, err
}
