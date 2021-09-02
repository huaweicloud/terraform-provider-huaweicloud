package topics

import (
	"github.com/chnsz/golangsdk"
)

// CreateResponse is a struct that contains the create response
type CreateResponse struct {
	Name string `json:"name"`
}

// Topic includes the parameters of an topic
type Topic struct {
	Name             string      `json:"name"`
	Partition        int         `json:"partition"`
	Replication      int         `json:"replication"`
	RetentionTime    int         `json:"retention_time"`
	SyncReplication  bool        `json:"sync_replication"`
	SyncMessageFlush bool        `json:"sync_message_flush"`
	TopicType        int         `json:"topic_type"`
	PoliciesOnly     bool        `json:"policiesOnly"`
	ExternalConfigs  interface{} `json:"external_configs"`
}

// ListResponse is a struct that contains the list response
type ListResponse struct {
	Total            int     `json:"total"`
	Size             int     `json:"size"`
	RemainPartitions int     `json:"remain_partitions"`
	MaxPartitions    int     `json:"max_partitions"`
	Topics           []Topic `json:"topics"`
}

// TopicDetail includes the detail parameters of an topic
type TopicDetail struct {
	Name            string      `json:"topic"`
	Partitions      []Partition `json:"partitions"`
	GroupSubscribed []string    `json:"group_subscribed"`
}

// Partition represents the details of a partition
type Partition struct {
	Partition int       `json:"partition"`
	Replicas  []Replica `json:"replicas"`
	// Node ID
	Leader int `json:"leader"`
	// Log End Offset
	Leo int `json:"leo"`
	// High Watermark
	Hw int `json:"hw"`
	// Log Start Offset
	Lso int `json:"lso"`
	// time stamp
	UpdateTimestamp int64 `json:"last_update_timestamp"`
}

// Replica represents the details of a replica
type Replica struct {
	Broker int  `json:"broker"`
	Leader bool `json:"leader"`
	InSync bool `json:"in_sync"`
	Size   int  `json:"size"`
	Lag    int  `json:"lag"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*CreateResponse, error) {
	var s CreateResponse
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// GetResult is a struct which contains the result of query
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*TopicDetail, error) {
	var s TopicDetail
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// ListResult contains the body of getting detailed
type ListResult struct {
	golangsdk.Result
}

// Extract from ListResult
func (r ListResult) Extract() ([]Topic, error) {
	var s ListResponse
	err := r.Result.ExtractInto(&s)
	return s.Topics, err
}

// UpdateResult is a struct from which can get the result of update method
type UpdateResult struct {
	golangsdk.ErrResult
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.Result
}

// DeleteResponse is a struct that contains the deletion response
type DeleteResponse struct {
	Name    string `json:"id"`
	Success bool   `json:"success"`
}

// Extract from DeleteResult
func (r DeleteResult) Extract() ([]DeleteResponse, error) {
	var s struct {
		Topics []DeleteResponse `json:"topics"`
	}
	err := r.Result.ExtractInto(&s)
	return s.Topics, err
}
