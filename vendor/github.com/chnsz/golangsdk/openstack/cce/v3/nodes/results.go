package nodes

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// Describes the Node Structure of cluster
type ListNode struct {
	// API type, fixed value "List"
	Kind string `json:"kind"`
	// API version, fixed value "v3"
	Apiversion string `json:"apiVersion"`
	// all Clusters
	Nodes []Nodes `json:"items"`
}

// Individual nodes of the cluster
type Nodes struct {
	//  API type, fixed value " Host "
	Kind string `json:"kind"`
	// API version, fixed value v3
	Apiversion string `json:"apiVersion"`
	// Node metadata
	Metadata Metadata `json:"metadata"`
	// Node detailed parameters
	Spec Spec `json:"spec"`
	// Node status information
	Status Status `json:"status"`
}

// Metadata required to create a node
type Metadata struct {
	//Node name
	Name string `json:"name"`
	//Node ID
	Id string `json:"uid"`
	// Node tag, key value pair format
	Labels map[string]string `json:"labels,omitempty"`
	//Node annotation, keyvalue pair format
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Spec describes Nodes specification
type Spec struct {
	// Node specifications
	Flavor string `json:"flavor" required:"true"`
	// The value of the available partition name
	Az string `json:"az" required:"true"`
	// The OS of the node
	Os string `json:"os,omitempty"`
	// ID of the dedicated host to which nodes will be scheduled
	DedicatedHostID string `json:"dedicatedHostId,omitempty"`
	// Node login parameters
	Login LoginSpec `json:"login" required:"true"`
	// System disk parameter of the node
	RootVolume VolumeSpec `json:"rootVolume" required:"true"`
	// The data disk parameter of the node must currently be a disk
	DataVolumes []VolumeSpec `json:"dataVolumes" required:"true"`
	// Disk initialization configuration management parameters
	// If omit, disk management is performed according to the DockerLVMConfigOverride parameter in extendParam
	Storage *StorageSpec `json:"storage,omitempty"`
	// Elastic IP parameters of the node
	PublicIP PublicIPSpec `json:"publicIP,omitempty"`
	// The billing mode of the node: the value is 0 (on demand)
	BillingMode int `json:"billingMode,omitempty"`
	// Number of nodes when creating in batch
	Count int `json:"count" required:"true"`
	// The node nic spec
	NodeNicSpec NodeNicSpec `json:"nodeNicSpec,omitempty"`
	// Extended parameter
	ExtendParam map[string]interface{} `json:"extendParam,omitempty"`
	// UUID of an ECS group
	EcsGroupID string `json:"ecsGroupId,omitempty"`
	// Tag of a VM, key value pair format
	UserTags []tags.ResourceTag `json:"userTags,omitempty"`
	// Tag of a Kubernetes node, key value pair format
	K8sTags map[string]string `json:"k8sTags,omitempty"`
	// The runtime spec
	RunTime *RunTimeSpec `json:"runtime,omitempty"`
	// taints to created nodes to configure anti-affinity
	Taints []TaintSpec `json:"taints,omitempty"`
	// The name of the created partition
	Partition string `json:"partition,omitempty"`
	// The initialized conditions
	InitializedConditions []string `json:"initializedConditions,omitempty"`
}

// Gives the Nic spec of the node
type NodeNicSpec struct {
	// The primary Nic of the Node
	PrimaryNic PrimaryNic `json:"primaryNic,omitempty"`
	// The extension Nics of the Node
	ExtNics []ExtNic `json:"extNics,omitempty"`
}

// Gives the Primary Nic of the node
type PrimaryNic struct {
	// The Subnet ID of the primary Nic
	SubnetId string `json:"subnetId,omitempty"`
	// Fixed ips of the primary Nic
	FixedIps []string `json:"fixedIps,omitempty"`
}

type ExtNic struct {
	// The Subnet ID of the extension Nic
	SubnetId string `json:"subnetId,omitempty"`
	// Fixed ips of the extension Nic
	FixedIps []string `json:"fixedIps,omitempty"`
	// IP block of the extension Nic
	IPBlock string `json:"ipBlock,omitempty"`
}

// TaintSpec to created nodes to configure anti-affinity
type TaintSpec struct {
	Key   string `json:"key" required:"true"`
	Value string `json:"value,omitempty"`
	// Available options are NoSchedule, PreferNoSchedule, and NoExecute
	Effect string `json:"effect" required:"true"`
}

// Gives the current status of the node
type Status struct {
	// The state of the Node
	Phase string `json:"phase"`
	// The virtual machine ID of the node in the ECS
	ServerID string `json:"ServerID"`
	// Elastic IP of the node
	PublicIP string `json:"PublicIP"`
	//Private IP of the node
	PrivateIP string `json:"privateIP"`
	// The ID of the Job that is operating asynchronously in the Node
	JobID string `json:"jobID"`
	// Reasons for the Node to become current
	Reason string `json:"reason"`
	// Details of the node transitioning to the current state
	Message string `json:"message"`
}

type LoginSpec struct {
	// Select the key pair name when logging in by key pair mode
	SshKey string `json:"sshKey,omitempty"`
	// Select the user/password when logging in
	UserPassword UserPassword `json:"userPassword,omitempty"`
}

type UserPassword struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type VolumeSpec struct {
	// Disk size in GB
	Size int `json:"size" required:"true"`
	// Disk type
	VolumeType string `json:"volumetype" required:"true"`
	//hw:passthrough
	HwPassthrough bool `json:"hw:passthrough,omitempty"`
	// Disk extension parameter
	ExtendParam map[string]interface{} `json:"extendParam,omitempty"`
	// Disk encryption information.
	Metadata *VolumeMetadata `json:"metadata,omitempty"`
	// DSS pool ID
	ClusterID string `json:"cluster_id,omitempty"`
	// DSS pool type, fixed to dss
	ClusterType string `json:"cluster_type,omitempty"`
}

type VolumeMetadata struct {
	// Whether the EVS disk is encrypted.
	// The value 0 indicates that the EVS disk is not encrypted,
	// and the value 1 indicates that the EVS disk is encrypted.
	SystemEncrypted string `json:"__system__encrypted,omitempty"`
	// CMK ID, which indicates encryption in metadata.
	SystemCmkid string `json:"__system__cmkid,omitempty"`
}

type PublicIPSpec struct {
	// List of existing elastic IP IDs
	Ids []string `json:"ids,omitempty"`
	// The number of elastic IPs to be dynamically created
	Count int `json:"count,omitempty"`
	// Elastic IP parameters
	Eip EipSpec `json:"eip,omitempty"`
}

type EipSpec struct {
	// The value of the iptype keyword
	IpType string `json:"iptype,omitempty"`
	// Elastic IP bandwidth parameters
	Bandwidth BandwidthOpts `json:"bandwidth,omitempty"`
}

type RunTimeSpec struct {
	// the name of runtime: docker or containerd
	Name string `json:"name,omitempty"`
}

type BandwidthOpts struct {
	ChargeMode string `json:"chargemode,omitempty"`
	Size       int    `json:"size,omitempty"`
	ShareType  string `json:"sharetype,omitempty"`
}

type StorageSpec struct {
	// Disk selection. Matched disks are managed according to matchLabels and storageType
	StorageSelectors []StorageSelectorsSpec `json:"storageSelectors" required:"true"`
	// A storage group consists of multiple storage devices. It is used to divide storage space
	StorageGroups []StorageGroupsSpec `json:"storageGroups" required:"true"`
}

type StorageSelectorsSpec struct {
	// Selector name, used as the index of selectorNames in storageGroup, the name of each selector must be unique
	Name string `json:"name" required:"true"`
	// Specifies the storage type. Currently, only evs and local are supported
	// The local storage does not support disk selection. All local disks will form a VG
	// Therefore, only one storageSelector of the local type is allowed
	StorageType string `json:"storageType" required:"true"`
	// Matching field of an EVS volume
	MatchLabels MatchLabelsSpec `json:"matchLabels,omitempty"`
}

type MatchLabelsSpec struct {
	// Matched disk size if left unspecified, the disk size is not limited
	Size string `json:"size,omitempty"`
	// EVS disk type
	VolumeType string `json:"volumeType,omitempty"`
	// Disk encryption identifier
	// 0 indicates that the disk is not encrypted, and 1 indicates that the disk is encrypted
	MetadataEncrypted string `json:"metadataEncrypted,omitempty"`
	// Customer master key ID of an encrypted disk
	MetadataCmkid string `json:"metadataCmkid,omitempty"`
	// Number of disks to be selected, if left blank, all disks of this type are selected
	Count string `json:"count,omitempty"`
}

type StorageGroupsSpec struct {
	// Name of a virtual storage group, each group name must be unique
	Name string `json:"name" required:"true"`
	// Storage space for Kubernetes and runtime components
	// Only one group can be set to true, default value is false
	CceManaged bool `json:"cceManaged,omitempty"`
	// This parameter corresponds to name in storageSelectors
	// A group can match multiple selectors, but a selector can match only one group
	SelectorNames []string `json:"selectorNames" required:"true"`
	// Detailed management of space configuration in a group
	VirtualSpaces []VirtualSpacesSpec `json:"virtualSpaces" required:"true"`
}

type VirtualSpacesSpec struct {
	// virtualSpace name, currently, only kubernetes, runtime, and user are supported
	// kubernetes and user require lvmConfig to be configured, runtime requires runtimeConfig to be configured
	Name string `json:"name" required:"true"`
	// Size of a virtual space, only an integer percentage is supported, example: 90%
	// Note that the total percentage of all virtual spaces in a group cannot exceed 100%
	Size string `json:"size" required:"true"`
	// LVM configurations, applicable to kubernetes and user spaces
	// One virtual space supports only one config
	LVMConfig *LVMConfigSpec `json:"lvmConfig,omitempty"`
	// runtime configurations, applicable to the runtime space
	// One virtual space supports only one config
	RuntimeConfig *RuntimeConfigSpec `json:"runtimeConfig,omitempty"`
}

type LVMConfigSpec struct {
	// LVM write mode, values can be linear and striped
	LvType string `json:"lvType" required:"true"`
	// Path to which the disk is attached, this parameter takes effect only in user configuration
	// The value is an absolute path
	Path string `json:"path,omitempty"`
}

type RuntimeConfigSpec struct {
	// LVM write mode, values can be linear and striped
	LvType string `json:"lvType" required:"true"`
}

// Describes the Job Structure
type Job struct {
	// API type, fixed value "Job"
	Kind string `json:"kind"`
	// API version, fixed value "v3"
	Apiversion string `json:"apiVersion"`
	// Node metadata
	Metadata JobMetadata `json:"metadata"`
	// Node detailed parameters
	Spec JobSpec `json:"spec"`
	//Node status information
	Status JobStatus `json:"status"`
}

type JobMetadata struct {
	// ID of the job
	ID string `json:"uid"`
}

type JobSpec struct {
	// Type of job
	Type string `json:"type"`
	// ID of the cluster where the job is located
	ClusterID string `json:"clusterUID"`
	// ID of the IaaS resource for the job operation
	ResourceID string `json:"resourceID"`
	// The name of the IaaS resource for the job operation
	ResourceName string `json:"resourceName"`
	// List of child jobs
	SubJobs []Job `json:"subJobs"`
	// ID of the parent job
	OwnerJob string `json:"ownerJob"`
}

type JobStatus struct {
	// Job status
	Phase string `json:"phase"`
	// The reason why the job becomes the current state
	Reason string `json:"reason"`
	// The job becomes the current state details
	Message string `json:"message"`
}

type AddNodeResponse struct {
	JobID string `json:"jobid"`
}

func (r commonResult) ExtractAddNode() (*AddNodeResponse, error) {
	var s AddNodeResponse
	err := r.ExtractInto(&s)
	return &s, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a node.
func (r commonResult) Extract() (*Nodes, error) {
	var s Nodes
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractNode is a function that accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) ExtractNode() ([]Nodes, error) {
	var s ListNode
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Nodes, nil
}

// ExtractJob is a function that accepts a result and extracts a job.
func (r commonResult) ExtractJob() (*Job, error) {
	var s Job
	err := r.ExtractInto(&s)
	return &s, err
}

// ListResult represents the result of a list operation. Call its ExtractNode
// method to interpret it as a Nodes.
type ListResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Node.
type CreateResult struct {
	commonResult
}

type AddResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Node.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Node.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
