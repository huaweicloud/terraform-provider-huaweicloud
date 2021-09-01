package clusters

import (
	"github.com/chnsz/golangsdk"
)

// Cluster is a struct that represents the result of Create methods.
type Cluster struct {
	ID string `json:"cluster_id"`
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	golangsdk.Result
}

// Extract is a method which to extract a cluster response.
func (r CreateResult) Extract() (*Cluster, error) {
	var s Cluster
	err := r.ExtractInto(&s)
	return &s, err
}
