package replications

import (
	"github.com/chnsz/golangsdk"
)

type Replication struct {
	//Replication ID
	Id string `json:"id"`
	//Replication Name
	Name string `json:"name"`
	//Replication Description
	Description string `json:"description"`
	//Replication Model
	ReplicaModel string `json:"replication_model"`
	//Replication Status
	Status string `json:"status"`
	//Replication Attachment
	Attachment []Attachment `json:"attachment"`
	//Replication Group ID
	GroupID string `json:"server_group_id"`
	//Replication Volume IDs
	VolumeIDs string `json:"volume_ids"`
	//Replication Priority Station
	PriorityStation string `json:"priority_station"`
	//Replication Fault Level
	FaultLevel string `json:"fault_level"`
	//Replication Record Metadata
	RecordMetadata RecordMetadata `json:"record_metadata"`
}

type Attachment struct {
	//Device Name
	Device string `json:"device"`
	//Protected Instance ID
	ProtectedInstance string `json:"protected_instance"`
}

type RecordMetadata struct {
	//Whether Multiattach
	Multiattach bool `json:"multiattach"`
	//Whether Bootable
	Bootable bool `json:"bootable"`
	//Volume Size
	VolumeSize int `json:"volume_size"`
	//Volume Type
	VolumeType string `json:"volume_type"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a replication.
func (r commonResult) Extract() (*Replication, error) {
	var response Replication
	err := r.ExtractInto(&response)
	return &response, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "replication")
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Replication.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Replication.
type GetResult struct {
	commonResult
}
