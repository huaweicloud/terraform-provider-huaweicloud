package drill

import (
	"github.com/chnsz/golangsdk"
)

type Drill struct {
	// DR-Drill ID
	Id string `json:"id"`
	// DR-Drill Name
	Name string `json:"name"`
	// DR-Drill Status
	Status string `json:"status"`
	// DR-Drill VPC ID
	DrillVpcID string `json:"drill_vpc_id"`
	// DR-Drill Group ID
	GroupID string `json:"server_group_id"`
	// DR-Drill Volume IDs
	Servers []Servers `json:"drill_servers"`
}

type Servers struct {
	// Protected Instance ID
	ProtectedInstance string `json:"protected_instance"`
	// Drill server ID
	ServerID string `json:"drill_server_id"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a dr-drill.
func (r commonResult) Extract() (*Drill, error) {
	var response Drill
	err := r.ExtractInto(&response)
	return &response, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "disaster_recovery_drill")
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a dr-drill.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a dr-drill.
type GetResult struct {
	commonResult
}
