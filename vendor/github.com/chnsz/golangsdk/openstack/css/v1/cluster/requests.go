package cluster

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

const (
	// Instance type. The options are ess, ess-cold, ess-master, and ess-client.
	InstanceTypeEss       = "ess"
	InstanceTypeEssCode   = "ess-cold"
	InstanceTypeEssMaster = "ess-master"
	InstanceTypeEssClient = "ess-client"

	ClusterStatusInProcess   = "100" //The operation, such as instance creation, is in progress.
	ClusterStatusAvailable   = "200"
	ClusterStatusUnavailable = "303"
)

// Opts
type CreateOpts struct {
	// Instance. For details about related parameters, see Table 4.
	Instance *InstanceBody `json:"instance" required:"true"`
	// Type of the data search engine. For details about related parameters, see Table 7.
	Datastore *DatastoreBody `json:"datastore,omitempty"`
	// Cluster name. It contains 4 to 32 characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// The value must start with a letter.
	Name string `json:"name" required:"true"`
	// Number of clusters. The value range is 1 to 32.
	InstanceNum int `json:"instanceNum" required:"true"`
	// Automatic snapshot creation. This function is enabled by default. For details about related parameters.
	BackupStrategy *BackupStrategyBody `json:"backupStrategy,omitempty"`
	// Whether disks are encrypted. For details about related parameters, see Table 9.
	DiskEncryption *EncryptionBody `json:"diskEncryption,omitempty"`
	// Whether communication encryption is performed on the cluster. Available values include true and false.
	// By default, communication encryption is disabled. When httpsEnable is set to true,
	// authorityEnable must be set to true.
	// Value true indicates that communication encryption is performed on the cluster.
	// Value false indicates that communication encryption is not performed on the cluster.
	// NOTE:
	// This parameter is supported in clusters 6.5.4 or later.
	HttpsEnable bool `json:"httpsEnable,omitempty"`
	// Whether to enable authentication. Available values include true and false. Authentication is disabled by default.
	// When authentication is enabled, httpsEnable must be set to true.
	// Value true indicates that authentication is enabled for the cluster.
	// Value false indicates that authentication is disabled for the cluster.
	// NOTE:
	// This parameter is supported in clusters 6.5.4 or later.
	AuthorityEnable bool `json:"authorityEnable,omitempty"`
	// Password of the cluster user admin in security mode. This parameter is mandatory only when authorityEnable
	// is set to true.
	// NOTE:
	// The administrator password must meet the following requirements:
	// The password can contain 8 to 32 characters.
	// Passwords must contain at least 3 of the following character types: uppercase letters, lowercase letters,
	// numbers, and special characters (~!@#$%^&*()-_=+\\|[{}];:,<.>/?).
	// Weak password verification is required for a security cluster. You are advised to set a strong password.
	AdminPwd string `json:"adminPwd,omitempty"`
	// Enterprise project ID. When creating a cluster, associate the enterprise project ID with the cluster. The value
	// can contain a maximum of 36 characters. It is string 0 or in UUID format with hyphens (-). Value 0 indicates
	// the default enterprise project.
	// NOTE:
	// For details about how to obtain enterprise project IDs and features, see the Enterprise Management Service User
	// Guide.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// Tags in a cluster.
	// NOTE:
	// For details about the tag feature, see the Tag Management Service Overview.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

type BackupStrategyBody struct {
	// Time when a snapshot is created every day. Snapshots can only be created on the hour. The time format is the time
	// followed by the time zone, specifically, HH:mm z. In the format, HH:mm refers to the hour time and z
	// refers to the time zone, for example, 00:00 GMT+08:00 and 01:00 GMT+08:00.
	Period string `json:"period" required:"true"`
	// Prefix of the name of the snapshot that is automatically created.
	Prefix string `json:"prefix" required:"true"`
	// Number of days for which automatically created snapshots are reserved.
	// Value range: 1 to 90
	Keepday int `json:"keepday" required:"true"`
	// The name of the OBS bucket used for backup. If the bucket already stores snapshot data, it cannot be changed.
	Bucket string `json:"bucket,omitempty"`
	// Storage path of the snapshot in the OBS bucket.
	BasePath string `json:"basePath,omitempty"`
	// The name of the IAM delegate used to access OBS.
	// illustrate:
	// If the three parameters bucket, basePath, and agency are empty at the same time, the system will automatically
	// create an OBS bucket and IAM agent, otherwise the configured parameter values will be used.
	Agency string `json:"agency,omitempty"`
}

// DatastoreBody
type DatastoreBody struct {
	// Cluster type. The default value is Elasticsearch. Currently, the value can only be Elasticsearch.
	Type string `json:"type,omitempty"`
	// Cluster version. The value can be 5.5.1, 6.2.3, 6.5.4, 7.1.1, 7.6.2, or 7.9.3. The default value is 5.5.1.
	Version string `json:"version" required:"true"`
}

// InstanceBody
type InstanceBody struct {
	// Instance flavor name. For example:
	// Value range of flavor ess.spec-2u16g: 40 GB to 1,280 GB
	// Value range of flavor ess.spec-4u32g: 40 GB to 2,560 GB
	// Value range of flavor ess.spec-8u64g: 80 GB to 5,120 GB
	// Value range of flavor ess.spec-16u128g: 160 GB to 10,240 GB
	FlavorRef string `json:"flavorRef" required:"true"`
	// If flavorRef is set to a local disk flavor, you do not need to set this parameter. You can obtain the local disk
	// flavor by calling the API for obtaining the instance flavor list. Currently,
	// the following local disk flavors are supported:
	// ess.spec-i3small
	// ess.spec-i3medium
	// ess.spec-i3.8xlarge.8
	// ess.spec-ds.xlarge.8
	// ess.spec-ds.2xlarge.8
	// ess.spec-ds.4xlarge.8
	// Information about the volume. For details about related parameters, see Table 5.
	Volume InstanceVolumeBody `json:"volume"`
	// Subnet information. For details about related parameters, see Table 6.
	Nics InstanceNicsBody `json:"nics" required:"true"`
	// Availability zone (AZ). A single AZ is created when this parameter is not specified.
	// Separate multiple AZs with commas (,), for example, az1,az2. AZs must be unique and ensure that the number of
	// nodes be at least the number of AZs.
	// If the number of nodes is a multiple of the number of AZs, the nodes are evenly distributed to each AZ. If the
	// number of nodes is not a multiple of the number of AZs, the absolute difference between node quantity in any two
	// AZs is 1 at most.
	AvailabilityZone string `json:"availability_zone"`
}

// InstanceNicsBody
type InstanceNicsBody struct {
	// Subnet ID. All instances in a cluster must have the same subnets and security groups.
	NetId string `json:"netId" required:"true"`
	// Security group ID. All instances in a cluster must have the same subnets and security groups.
	SecurityGroupId string `json:"securityGroupId" required:"true"`
	// VPC ID, which is used for configuring cluster network.
	VpcId string `json:"vpcId" required:"true"`
}

// InstanceVolumeBody volume
type InstanceVolumeBody struct {
	// Volume size, which must be a multiple of 4 and 10.  Unit: GB
	Size int `json:"size" required:"true"`
	// COMMON: Common I/O
	// HIGH: High I/O
	// ULTRAHIGH: Ultra-high I/O
	VolumeType string `json:"volume_type" required:"true"`
}

type EncryptionBody struct {
	// Value 1 indicates encryption is performed, and value 0 indicates encryption is not performed.
	SystemEncrypted string `json:"systemEncrypted" required:"true"`
	// Key ID.
	// The Default Master Keys cannot be used to create grants. Specifically, you cannot use Default Master Keys whose
	// aliases end with /default in KMS to create clusters.
	// After a cluster is created, do not delete the key used by the cluster. Otherwise,
	// the cluster will become unavailable.
	SystemCmkid string `json:"systemCmkid" required:"true"`
}

type ListClustersDetailsOpts struct {
	// Start value of the query. The default value is 1, indicating that the query starts from the first cluster.
	Start int `q:"start"`
	// Number of clusters to be queried. The default value is 10, indicating that 10 clusters are queried at a time.
	Limit int `q:"limit"`
}

type RoleExtendReq struct {
	// Detailed description about the cluster scale-out request. For detai
	Grow []RoleExtendGrowReq `json:"grow" required:"true"`
}

// RoleExtendGrowReq
type RoleExtendGrowReq struct {
	// Storage capacity of the instance to be expanded. The total storage capacity of existing instances
	// and newly added instances in a cluster cannot exceed the maximum instance storage capacity allowed when a
	// cluster is being created. In addition, you can expand the instance storage capacity for a cluster
	// for up to six times.Unit: GB
	Disksize *int `json:"disksize" required:"true"`
	// Number of instances to be scaled out. The total number of existing instances and newly added instances
	// in a cluster cannot exceed 32.
	Nodesize *int `json:"nodesize" required:"true"`
	// Type of the instance to be scaled out. Select at least one from ess, ess-cold, ess-master, and ess-client.
	// You can only add instances rather than increase storage capacity on nodes of the ess-master and ess-client types.
	Type string `json:"type" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*CreateResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "cluster")
	if err != nil {
		return nil, err
	}

	var r CreateResponse

	_, err = c.Post(createURL(c), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	if err == nil {
		return &r, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, clusterId string) (*ClusterDetailResponse, error) {
	var rst ClusterDetailResponse
	_, err := c.Get(getURL(c, clusterId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	if err == nil {
		return &rst, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, clusterId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(deleteURL(c, clusterId), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}

// ListClustersDetails
func List(c *golangsdk.ServiceClient, opts ListClustersDetailsOpts) (*ClusterListResponse, error) {
	url := listURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r ClusterListResponse
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func ExtendInstanceStorage(c *golangsdk.ServiceClient, clusterId string, opts RoleExtendReq) (*ClusterResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst ClusterResponse
	_, err = c.Post(extendInstanceStorageURL(c, clusterId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}

// ListFlavors
func ListFlavors(c *golangsdk.ServiceClient) (*EsFlavorsResp, error) {
	var rst EsFlavorsResp
	_, err := c.Get(listFlavorsURL(c), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err

}

// RestartCluster
func Restart(c *golangsdk.ServiceClient, clusterId string) (*RestartClusterResponse, error) {

	var rst RestartClusterResponse
	_, err := c.Post(restartURL(c, clusterId), nil, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}
