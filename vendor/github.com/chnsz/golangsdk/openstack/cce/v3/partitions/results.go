package partitions

import (
	"github.com/chnsz/golangsdk"
)

// Describes the Partition Structure of cluster
type ListPartition struct {
	// API type, fixed value "List"
	Kind string `json:"kind"`
	// API version, fixed value "v3"
	Apiversion string `json:"apiVersion"`
	// all Clusters
	Partitions []Partitions `json:"items"`
}

// Individual partitions of the cluster
type Partitions struct {
	// API type, fixed value " Host "
	Kind string `json:"kind"`
	// API version, fixed value v3
	Apiversion string `json:"apiVersion"`
	// Partition metadata
	Metadata Metadata `json:"metadata"`
	// Partition detailed parameters
	Spec Spec `json:"spec"`
}

// Metadata required to create a partition
type Metadata struct {
	// Partition name
	Name string `json:"name"`
}

// Spec describes Partitions specification
type Spec struct {
	// The category of partition
	Category string `json:"category,omitempty"`
	// The availability zone name of the partition
	PublicBorderGroup string `json:"publicBorderGroup,omitempty"`
	// The default host network for the partition
	HostNetwork HostNetwork `json:"hostNetwork,omitempty"`
	// The default host network for the partition container
	ContainerNetwork []ContainerNetwork `json:"containerNetwork,omitempty"`
}

// HostNetwork the default host network for the partition
type HostNetwork struct {
	// The default SubnetID for the partition
	SubnetID string `json:"subnetID,omitempty"`
}

// ContainerNetwork the default host network for the partition container
type ContainerNetwork struct {
	// The default SubnetID for the partition container
	SubnetID string `json:"subnetID,omitempty"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a partition.
func (r commonResult) Extract() (*Partitions, error) {
	var s Partitions
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractPartitions is a function that accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func (r commonResult) ExtractPartitions() ([]Partitions, error) {
	var s ListPartition
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s.Partitions, nil
}

// ListResult represents the result of a list operation. Call its ExtractCluster
// method to interpret it as a Cluster.
type ListResult struct {
	commonResult
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
