package clusters

import (
	"github.com/huaweicloud/golangsdk"
)

type ListCluster struct {
	// API type, fixed value Cluster
	Kind string `json:"kind"`
	//API version, fixed value v3
	ApiVersion string `json:"apiVersion"`
	//all Clusters
	Clusters []Clusters `json:"items"`
}

type Clusters struct {
	// API type, fixed value Cluster
	Kind string `json:"kind" required:"true"`
	//API version, fixed value v3
	ApiVersion string `json:"apiversion" required:"true"`
	//Metadata of a Cluster
	Metadata MetaData `json:"metadata" required:"true"`
	//specifications of a Cluster
	Spec Spec `json:"spec" required:"true"`
	//status of a Cluster
	Status Status `json:"status"`
}

//Metadata required to create a cluster
type MetaData struct {
	//Cluster unique name
	Name string `json:"name"`
	//Cluster unique Id
	Id string `json:"uid"`
	// Cluster tag, key/value pair format
	Labels map[string]string `json:"labels,omitempty"`
	//Cluster annotation, key/value pair format
	Annotations map[string]string `json:"annotations,omitempty"`
}

//Specifications to create a cluster
type Spec struct {
	//Cluster Type: VirtualMachine, BareMetal, or Windows
	Type string `json:"type" required:"true"`
	// Cluster specifications
	Flavor string `json:"flavor" required:"true"`
	// For the cluster version, please fill in v1.7.3-r10 or v1.9.2-r1. Currently only Kubernetes 1.7 and 1.9 clusters are supported.
	Version string `json:"version,omitempty"`
	//Cluster description
	Description string `json:"description,omitempty"`
	// Node network parameters
	HostNetwork HostNetworkSpec `json:"hostNetwork" required:"true"`
	//Container network parameters
	ContainerNetwork ContainerNetworkSpec `json:"containerNetwork" required:"true"`
	// Charging mode of the cluster, which is 0 (on demand)
	BillingMode int `json:"billingMode,omitempty"`
	//Extended parameter for a cluster
	ExtendParam map[string]string `json:"extendParam,omitempty"`
}

// Node network parameters
type HostNetworkSpec struct {
	//The ID of the VPC used to create the node
	VpcId string `json:"vpc" required:"true"`
	//The ID of the subnet used to create the node
	SubnetId string `json:"subnet" required:"true"`
	// The ID of the high speed network used to create bare metal nodes.
	// This parameter is required when creating a bare metal cluster.
	HighwaySubnet string `json:"highwaySubnet,omitempty"`
}

//Container network parameters
type ContainerNetworkSpec struct {
	//Container network type: overlay_l2 , underlay_ipvlan or vpc-router
	Mode string `json:"mode" required:"true"`
	//Container network segment: 172.16.0.0/16 ~ 172.31.0.0/16. If there is a network segment conflict, it will be automatically reselected.
	Cidr string `json:"cidr,omitempty"`
}

type Status struct {
	//The state of the cluster
	Phase string `json:"phase"`
	//The ID of the Job that is operating asynchronously in the cluster
	JobID string `json:"jobID"`
	//Reasons for the cluster to become current
	Reason string `json:"reason"`
	//The status of each component in the cluster
	Conditions Conditions `json:"conditions"`
	//Kube-apiserver access address in the cluster
	Endpoints []Endpoints `json:"endpoints"`
}

type Conditions struct {
	//The type of component
	Type string `json:"type"`
	//The state of the component
	Status string `json:"status"`
	//The reason that the component becomes current
	Reason string `json:"reason"`
}

type Endpoints struct {
	//The address accessed within the user's subnet
	Url string `json:"url"`
	//Public network access address
	Type string `json:"type"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a cluster.
func (r commonResult) Extract() (*Clusters, error) {
	var s Clusters
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractCluster is a function that accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) ExtractClusters() ([]Clusters, error) {
	var s ListCluster
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s.Clusters, nil

}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Cluster.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Cluster.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Cluster.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// ListResult represents the result of a list operation. Call its ExtractCluster
// method to interpret it as a Cluster.
type ListResult struct {
	commonResult
}
