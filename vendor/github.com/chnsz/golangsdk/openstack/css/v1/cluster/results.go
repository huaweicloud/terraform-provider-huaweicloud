package cluster

import "github.com/chnsz/golangsdk/openstack/common/tags"

// CreateClusterResponse This is a auto create Response Object
type CreateResponse struct {
	Cluster CreateResponseBody `json:"cluster"`
}

// CreateResponse
type CreateResponseBody struct {
	Id   string `json:"id"`   // Cluster ID
	Name string `json:"name"` // Cluster name
}

type ClusterResponse struct {
	Id string `json:"id"` // Cluster ID
}

type ClusterDetailResponse struct {
	// Type of the data search engine. For details, see Table 3.
	Datastore ClusterDetailDatastore `json:"datastore"`
	// List of node objects.
	Instances []ClusterDetailInstances `json:"instances"`
	// Last modification time of a cluster. The format is ISO8601: CCYY-MM-DDThh:mm:ss.
	Updated string `json:"updated"`
	// Cluster name.
	Name string `json:"name"`
	// Time when a cluster is created. The format is ISO8601: CCYY-MM-DDThh:mm:ss.
	Created string `json:"created"`
	// Cluster ID.
	Id     string `json:"id"`
	Status string `json:"status"` //100:The operation is in progress.;200: available.;303: unavailable.
	// Indicates the IP address and port number of the user used to access the VPC.
	Endpoint string `json:"endpoint"`
	// Cluster operation progress, which indicates the progress of cluster creation and expansion in percentage.
	ActionProgress map[string]interface{} `json:"actionProgress"`
	// Current behavior on a cluster. Value REBOOTING indicates that the cluster is being restarted, GROWING indicates
	// that capacity expansion is being performed on the cluster, RESTORING indicates that the cluster is being
	// restored, and SNAPSHOTTING indicates that the snapshot is being created.
	Actions []string `json:"actions"`
	// Failure cause. If the cluster is in the Available state, this parameter is not returned.
	FailedReasons ClusterDetailFailedReasons `json:"failed_reasons"`
	// Whether to enable authentication. Available values include true and false. Authentication is disabled by
	// default. When authentication is enabled, httpsEnable must be set to true.
	// Value true indicates that authentication is enabled for the cluster.
	// Value false indicates that authentication is disabled for the cluster.
	AuthorityEnable bool `json:"authorityEnable"`
	// Whether disks are encrypted.
	// Value true indicates that disks are encrypted.
	// Value false indicates that disks are not encrypted.
	DiskEncrypted bool `json:"diskEncrypted"`
	// Key ID used for disk encryption.
	CmkId string `json:"cmkId"`
	// ID of the enterprise project to which a cluster belongs.
	// If the user of the cluster does not enable the enterprise project, the setting of this parameter is not returned.
	EnterpriseProjectId string             `json:"enterprise_project_id"`
	Tags                []tags.ResourceTag `json:"tags"`
}

type ClusterListResponse struct {
	ClusterDetailResponse

	SecurityGroupId string `json:"securityGroupId"`
	SubnetId        string `json:"subnetId"`
	VpcId           string `json:"vpcId"`
}

// ClusterDetailActionProgress
type ClusterDetailActionProgress struct {
	Creating string `json:"CREATING"`
}

type ClusterDetailDatastore struct {
	// Cluster type. The default value is Elasticsearch. Currently, the value can only be Elasticsearch.
	Type string `json:"type"`
	// Cluster version. The value can be 5.5.1, 6.2.3, 6.5.4, 7.1.1, 7.6.2, or 7.9.3.
	Version string `json:"version"`
}

// ClusterDetailFailedReasons
type ClusterDetailFailedReasons struct {
	// Error code.
	// CSS.6000: indicates that a cluster fails to be created.
	// CSS.6001: indicates that capacity expansion of a cluster fails.
	// CSS.6002: indicates that a cluster fails to be restarted.
	// CSS.6004: indicates that a node fails to be created in a cluster.
	// CSS.6005: indicates that the service fails to be initialized.
	ErrorCode string `json:"error_code"`
	// Detailed error information
	ErrorMsg string `json:"error_msg"`
}

type ClusterDetailInstances struct {
	Type     string `json:"type"` // Supported type: ess (indicating the Elasticsearch node)
	Id       string `json:"id"`
	Name     string `json:"name"`
	SpecCode string `json:"specCode"` // Node specifications.
	AzCode   string `json:"azCode"`   // AZ to which a node belongs.

	// Instance status.
	// 100: The operation, such as instance creation, is in progress.
	// 200: The instance is available.
	// 303: The instance is unavailable.
	Status string `json:"status"`
}

type EsFlavorsResp struct {
	// List of engine versions
	Versions []EsflavorsVersionsResp `json:"versions"`
}

type EsflavorsVersionsResp struct {
	// Engine version. Versions 5.5.1, 6.2.3, 6.5.4, 7.1.1, 7.6.2, and 7.9.3 are supported.
	Version string     `json:"version"`
	Type    string     `json:"type"` // Instance type. The options are ess, ess-cold, ess-master, and ess-client.
	Flavors []EsFlavor `json:"flavors"`
}

type EsFlavor struct {
	Ram       int    `json:"ram"`       // Memory size of an instance. Unit: GB
	Cpu       int    `json:"cpu"`       // Number of vCPUs of an instance.
	Name      string `json:"name"`      // Flavor name.
	Region    string `json:"region"`    // AZ
	Diskrange string `json:"diskrange"` // Disk capacity range of an instance.
	FlavorId  string `json:"flavor_id"` // ID of a flavor.
}

// RestartClusterResponse This is a auto create Response Object
type RestartClusterResponse struct {
	JobId string `json:"jobId"`
}
