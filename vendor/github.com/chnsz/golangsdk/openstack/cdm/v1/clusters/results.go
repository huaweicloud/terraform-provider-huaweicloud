package clusters

const (
	ActionProgressCreating     = "CREATING"
	ActionProgressGrowing      = "GROWING"
	ActionProgressRestoring    = "RESTORING"
	ActionProgressSnapshotting = "SNAPSHOTTING"
	ActionProgressRepairing    = "REPAIRING"

	StatusCreating       = "100"
	StatusNormal         = "200"
	StatusFailed         = "300"
	StatusCreationFailed = "303"
	StatusForzen         = "800"
	StatusStopped        = "900"
	StatusStopping       = "910"
	StatusStarting       = "920"
)

type ClusterCreateResult struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ClustersRepsonse struct {
	Clusters []Cluster4List `json:"clusters"`
}

type Cluster4List struct {
	ClusterCommon
	// Whether to enable the scheduled startup/shutdown function.
	// The scheduled startup/shutdown and auto shutdown functions cannot be enabled at the same time.
	IsScheduleBootOff bool                    `json:"isScheduleBootOff"`
	Version           string                  `json:"version"`
	FailedReasons     map[string]FailedDetail `json:"failedReasons"`
}

type CustomerConfig struct {
	// Failure notification
	FailureRemind string `json:"failureRemind"`
	ClusterName   string `json:"clusterName"`
	// Service provisioning
	ServiceProvider string `json:"serviceProvider"`
	// Whether the disk is a local disk
	LocalDisk string `json:"localDisk"`
	// Whether to enable SSL
	Ssl string `json:"ssl"`
}

type Instance struct {
	// VM flavor of a node.
	Flavor Flavor `json:"flavor"`
	// Disk information of a node.
	Volume Volume `json:"volume"`
	// Node status:
	//   - 100:creating
	//   - 200:normal
	//   - 300:failed
	//   - 303:creation failed
	//   - 400:deleted
	//   - 800:forzen
	Status string `json:"status"`
	// Cluster operation status list. The options are as follows:
	//   - REBOOTING:restarting
	//   - RESTORING:restoring
	//   - REBOOT_FAILURE:restart failed
	Actions []string `json:"actions"`
	// Node type. Currently, only cdm is available.
	Type string `json:"type"`
	// Node VM ID
	Id string `json:"id"`
	// Name of the VM on the node
	Name string `json:"name"`
	// Whether the node is frozen. The value can be 0 (not frozen) or 1 (frozen).
	IsFrozen string `json:"isFrozen"`
	// Cluster configuration status. The options are as follows:
	//  In-Sync: configuration synchronized
	//  Applying: being configured
	//  Sync-Failure: configuration failed
	ConfigStatus string         `json:"config_status"`
	Role         string         `json:"role"`
	Group        string         `json:"group"`
	Links        []ClusterLinks `json:"links"`
	// 	Group ID
	ParamsGroupId string `json:"paramsGroupId"`
	// Public IP address
	PublicIp string `json:"publicIp"`
	// Management IP address
	ManageIp string `json:"manageIp"`
	// Traffic IP address
	TrafficIp string `json:"trafficIp"`
	// Shard ID
	ShardId string `json:"shard_id"`
	// Management fix IP address
	ManageFixIp string `json:"manage_fix_ip"`
	// Private IP address
	PrivateIp string `json:"private_ip"`
	// Internal IP address
	InternalIp string     `json:"internal_ip"`
	Resource   []Resource `json:"resource"`
}

type Flavor struct {
	Id string `q:"id"`
}

type Volume struct {
	// 	Type of disks on the node. Only local disks are supported.
	Type string `q:"type"`
	// Size of the disk on the node (GB)
	Size int `q:"size"`
}

type ClusterLinks struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Resource struct {
	ResourceId   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type ClusterTask struct {
	Description string `json:"description"`
	Id          string `json:"id"`
	Name        string `json:"name"`
}

type FailedDetail struct {
	ErrorCode string `q:"errorCode"`
	ErrorMsg  string `q:"errorMsg"`
}

type Job struct {
	JobId string `json:"jobId"`
}

type ActionResponse struct {
	JobId []string `json:"jobId"`
}

type ClusterCommon struct {
	CustomerConfig CustomerConfig `json:"customerConfig"`
	Datastore      Datastore      `json:"datastore"`
	Instances      []Instance     `json:"instances"`
	AzName         string         `json:"azName"`
	Dbuser         string         `json:"dbuser"`
	FlavorName     string         `json:"flavorName"`
	// Number of events
	RecentEvent int         `json:"recentEvent"`
	IsAutoOff   bool        `json:"isAutoOff"`
	ClusterMode string      `json:"clusterMode"`
	Namespace   string      `json:"namespace"`
	Task        ClusterTask `json:"task"`
	// EIP bound to the cluster
	PublicEndpoint string `json:"publicEndpoint"`
	// Cluster operation progress, which consists of a key and a value. The key indicates an ongoing task,
	// and the value indicates the progress of the ongoing task. An example is "action_progress":{"SNAPSHOTTING":"16%"}.
	ActionProgress map[string]string `json:"actionProgress"`
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	// Cluster creation time in ISO8601: YYYY-MM-DDThh:mm:ssZ format
	Created string `json:"created"`
	// Time when a cluster is updated. The format is YYYY-MM-DDThh:mm:ssZ (ISO 8601).
	Updated string `json:"updated"`
	// Cluster status. The options are as follows:
	//   - 100:creating
	//   - 200:normal
	//   - 300:failed
	//   - 303:creation failed
	//   - 800:forzen
	//   - 900:stopped
	//   - 910:stopping
	//   - 920:starting
	Status       string `json:"status"`
	StatusDetail string `json:"statusDetail"`
	// Whether the cluster is frozen. The value can be 0 (not frozen) or 1 (frozen).
	IsFrozen string `json:"isFrozen"`
	// Cluster configuration status. The options are as follows:
	// In-Sync: configuration synchronized; Applying: being configured; Sync-Failure: configuration failed
	ConfigStatus string `json:"config_status"`
	// Cluster links
	Links []ClusterLinks `json:"links"`

	IsScheduleBootOff bool   `json:"isScheduleBootOff"`
	ScheduleBootTime  string `json:"scheduleBootTime"`
	ScheduleOffTime   string `json:"scheduleOffTime"`
}

type Cluster struct {
	ClusterCommon

	SecurityGroupId string `json:"security_group_id"`
	SubnetId        string `json:"subnet_id"`
	VpcId           string `json:"vpc_id"`
	// EIP domain name bound to the cluster
	PublicEndpointDomainName string `json:"publicEndpointDomainName"`
	// Start time
	BakExpectedStartTime string `json:"bakExpectedStartTime"`
	// Retention duration
	BakKeepDay int `json:"bakKeepDay"`
	// Maintenance window
	MaintainWindow MaintainWindow `json:"maintainWindow"`
	// EIP ID
	EipId string `json:"eipId"`
	// EIP status
	PublicEndpointStatus PublicEndpointStatus `json:"publicEndpointStatus"`
	// Cluster configuration status. The options are as follows: In-Sync: configuration synchronized;
	// Applying: being configured; Sync-Failure: configuration failed
	Actions []string `json:"actions"`
}

type PublicEndpointStatus struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type MaintainWindow struct {
	// Day of a week
	Day string `json:"day"`
	// Start time
	StartTime string `json:"startTime"`
	// End time
	EndTime string `json:"endTime"`
}
