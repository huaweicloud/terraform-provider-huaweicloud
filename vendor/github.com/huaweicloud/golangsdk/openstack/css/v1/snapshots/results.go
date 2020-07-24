package snapshots

import (
	"time"

	"github.com/huaweicloud/golangsdk"
)

// Policy contains all the information associated with a snapshot policy.
type Policy struct {
	KeepDay  int    `json:"keepday"`
	Period   string `json:"period"`
	Prefix   string `json:"prefix"`
	Bucket   string `json:"bucket"`
	BasePath string `json:"basePath"`
	Agency   string `json:"agency"`
	Enable   string `json:"enable"`
}

// Snapshot contains all the information associated with a Cluster Snapshot.
type Snapshot struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"backupType"`
	Method        string `json:"backupMethod"`
	Description   string `json:"description"`
	ClusterID     string `json:"clusterId"`
	ClusterName   string `json:"clusterName"`
	Indices       string `json:"indices"`
	TotalShards   int    `json:"totalShards"`
	FailedShards  int    `json:"failedShards"`
	KeepDays      int    `json:"backupKeepDay"`
	Period        string `json:"backupPeriod"`
	Bucket        string `json:"bucketName"`
	Version       string `json:"version"`
	Status        string `json:"status"`
	RestoreStatus string `json:"restoreStatus"`

	// type of the data search engine
	DataStore DataStore `json:"datastore"`

	// the information about times
	ExpectedStartTime time.Time `json:"-"`
	StartTime         time.Time `json:"-"`
	EndTime           time.Time `json:"-"`
	Created           string    `json:"created"`
	Updated           string    `json:"updated"`
}

type DataStore struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type commonResult struct {
	golangsdk.Result
}

// PolicyResult contains the response body and error from a policy request.
type PolicyResult struct {
	commonResult
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// ListResult contains the response body and error from a List request.
type ListResult struct {
	commonResult
}

// ErrorResult contains the response body and error from a request.
type ErrorResult struct {
	golangsdk.ErrResult
}

// Extract will get the Policy object out of the PolicyResult object.
func (r PolicyResult) Extract() (*Policy, error) {
	var pol Policy
	err := r.ExtractInto(&pol)
	return &pol, err
}

// Extract will get the Snapshot object out of the CreateResult object.
func (r CreateResult) Extract() (*Snapshot, error) {
	var s struct {
		Snapshot *Snapshot `json:"backup"`
	}
	err := r.ExtractInto(&s)
	return s.Snapshot, err
}

// Extract will get all Snapshot objects out of the ListResult object.
func (r ListResult) Extract() ([]Snapshot, error) {
	var s struct {
		Snapshots []Snapshot `json:"backups"`
	}
	err := r.ExtractInto(&s)
	return s.Snapshots, err
}
