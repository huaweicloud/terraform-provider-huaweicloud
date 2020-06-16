package nodes

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

//Describes the Node Structure of cluster
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
	// Elastic IP parameters of the node
	PublicIP PublicIPSpec `json:"publicIP,omitempty"`
	// The billing mode of the node: the value is 0 (on demand)
	BillingMode int `json:"billingMode,omitempty"`
	// Number of nodes when creating in batch
	Count int `json:"count" required:"true"`
	// The node nic spec
	NodeNicSpec NodeNicSpec `json:"nodeNicSpec,omitempty"`
	// Extended parameter
	ExtendParam ExtendParam `json:"extendParam,omitempty"`
	// UUID of an ECS group
	EcsGroupID string `json:"ecsGroupId,omitempty"`
	// Tag of a VM, key value pair format
	UserTags []tags.ResourceTag `json:"userTags,omitempty"`
	// Tag of a Kubernetes node, key value pair format
	K8sTags map[string]string `json:"k8sTags,omitempty"`
	// taints to created nodes to configure anti-affinity
	Taints []TaintSpec `json:"taints,omitempty"`
}

// Gives the Nic spec of the node
type NodeNicSpec struct {
	// The primary Nic of the Node
	PrimaryNic PrimaryNic `json:"primaryNic,omitempty"`
}

// Gives the Primary Nic of the node
type PrimaryNic struct {
	// The Subnet ID of the primary Nic
	SubnetId string `json:"subnetId,omitempty"`
}

// TaintSpec to created nodes to configure anti-affinity
type TaintSpec struct {
	Key   string `json:"key" required:"true"`
	Value string `json:"value" required:"true"`
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
	//The status of each component in the Node
	Conditions Conditions `json:"conditions"`
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
	// Disk extension parameter
	ExtendParam string `json:"extendParam,omitempty"`
}

type ExtendParam struct {
	// Node charging mode, 0 is on-demand charging.
	ChargingMode int `json:"chargingMode,omitempty"`
	// Classification of cloud server specifications.
	EcsPerformanceType string `json:"ecs:performancetype,omitempty"`
	// Order ID, mandatory when the node payment type is the automatic payment package period type.
	OrderID string `json:"orderID,omitempty"`
	// The Product ID.
	ProductID string `json:"productID,omitempty"`
	// The Public Key.
	PublicKey string `json:"publicKey,omitempty"`
	// The maximum number of instances a node is allowed to create.
	MaxPods int `json:"maxPods,omitempty"`
	// Script required before the installation.
	PreInstall string `json:"alpha.cce/preInstall,omitempty"`
	// Script required after the installation.
	PostInstall string `json:"alpha.cce/postInstall,omitempty"`
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

type BandwidthOpts struct {
	ChargeMode string `json:"chargemode,omitempty"`
	Size       int    `json:"size,omitempty"`
	ShareType  string `json:"sharetype,omitempty"`
}

type Conditions struct {
	//The type of component
	Type string `json:"type"`
	//The state of the component
	Status string `json:"status"`
	//The reason that the component becomes current
	Reason string `json:"reason"`
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
