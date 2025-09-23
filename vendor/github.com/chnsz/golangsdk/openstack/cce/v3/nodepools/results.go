package nodepools

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
)

// Describes the Node Pool Structure of cluster
type ListNodePool struct {
	// API type, fixed value "List"
	Kind string `json:"kind"`
	// API version, fixed value "v3"
	Apiversion string `json:"apiVersion"`
	// all Node Pools
	NodePools []NodePool `json:"items"`
}

// Individual node pools of the cluster
type NodePool struct {
	//  API type, fixed value " Host "
	Kind string `json:"kind"`
	// API version, fixed value v3
	Apiversion string `json:"apiVersion"`
	// Node Pool metadata
	Metadata Metadata `json:"metadata"`
	// Node Pool detailed parameters
	Spec Spec `json:"spec"`
	// Node Pool status information
	Status Status `json:"status"`
}

// Metadata of the node pool
type Metadata struct {
	//Node Pool name
	Name string `json:"name"`
	//Node Pool ID
	Id string `json:"uid"`
}

// Gives the current status of the node pool
type Status struct {
	// The state of the node pool
	Phase string `json:"phase"`
	// Number of nodes in the node pool
	CurrentNode int `json:"currentNode"`
}

// Spec describes Node pools specification
type Spec struct {
	// Node type. Currently, only VM nodes are supported.
	Type string `json:"type" required:"true"`
	// Node Pool template
	NodeTemplate nodes.Spec `json:"nodeTemplate" required:"true"`
	// Initial number of expected node pools
	InitialNodeCount int `json:"initialNodeCount" required:"true"`
	// Auto scaling parameters
	Autoscaling AutoscalingSpec `json:"autoscaling"`
	// Node pool management parameters
	NodeManagement NodeManagementSpec `json:"nodeManagement"`
	// Pod security group configurations
	PodSecurityGroups []PodSecurityGroupSpec `json:"podSecurityGroups"`
	// Node security group configurations
	CustomSecurityGroups []string `json:"customSecurityGroups"`
	// label (k8s tag) policy on existing nodes
	LabelPolicyOnExistingNodes string `json:"labelPolicyOnExistingNodes"`
	// tag policy on existing nodes
	UserTagPolicyOnExistingNodes string `json:"userTagsPolicyOnExistingNodes"`
	// taint policy on existing nodes
	TaintPolicyOnExistingNodes string `json:"taintPolicyOnExistingNodes"`
	// The list of extension scale groups
	ExtensionScaleGroups []ExtensionScaleGroups `json:"extensionScaleGroups"`
}

type AutoscalingSpec struct {
	// Whether to enable auto scaling
	Enable bool `json:"enable"`
	// Minimum number of nodes allowed if auto scaling is enabled
	MinNodeCount int `json:"minNodeCount"`
	// This value must be greater than or equal to the value of minNodeCount
	MaxNodeCount int `json:"maxNodeCount"`
	// Interval between two scaling operations, in minutes
	ScaleDownCooldownTime int `json:"scaleDownCooldownTime"`
	// Weight of a node pool
	Priority int `json:"priority"`
}

type NodeManagementSpec struct {
	// ECS group ID
	ServerGroupReference string `json:"serverGroupReference"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a node pool.
func (r commonResult) Extract() (*NodePool, error) {
	var s NodePool
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractNodePool is a function that accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) ExtractNodePool() ([]NodePool, error) {
	var s ListNodePool
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.NodePools, nil
}

// ListResult represents the result of a list operation. Call its ExtractNode
// method to interpret it as a Node Pool.
type ListResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Node Pool.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Node Pool.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Node Pool.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
