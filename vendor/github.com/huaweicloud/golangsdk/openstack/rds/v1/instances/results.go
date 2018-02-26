package instances

import "github.com/huaweicloud/golangsdk"

type Instance struct {
	ID               string              `json:"id"`
	Status           string              `json:"status"`
	Name             string              `json:"name"`
	Created          string              `json:"created"`
	HostName         string              `json:"hostname"`
	Type             string              `json:"type"`
	Region           string              `json:"region"`
	Updated          string              `json:"updated"`
	AvailabilityZone string              `json:"availabilityZone"`
	Vpc              string              `json:"vpc"`
	Nics             NicsInfor           `json:"nics"`
	SecurityGroup    SecurityGroupInfor  `json:"securityGroup"`
	Flavor           FlavorInfo          `json:"flavor"`
	Volume           VolumeInfor         `json:"volume"`
	DbPort           int                 `json:"dbPort"`
	DataStore        DataStoreInfo       `json:"dataStoreInfo"`
	ExtendParameters ExtendParamInfo     `json:"extendparam"`
	BackupStrategy   BackupStrategyInfor `json:"backupStrategy"`
	Ha               HaInfor             `json:"ha"`
	SlaveId          string              `json:"slaveId"`
}

type ExtendParamInfo struct {
	Jobs []Job `json:"jobs"`
}

type FlavorInfo struct {
	Id string `json:"id"`
}

type DataStoreInfo struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type VolumeInfor struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

type NicsInfor struct {
	SubnetId string `json:"subnetId"`
}

type SecurityGroupInfor struct {
	Id string `json:"id"`
}

type BackupStrategyInfor struct {
	StartTime string `json:"startTime"`
	KeepDays  int    `json:"keepDays"`
}

type HaInfor struct {
	Enable          bool   `json:"enable"`
	ReplicationMode string `json:"replicationMode"`
}

type Job struct {
	ID string `json:"id"`
}

// Extract will get the Instance object out of the commonResult object.
func (r commonResult) Extract() (*Instance, error) {
	var s Instance
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "instance")
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type ListResult struct {
	golangsdk.Result
}

func (lr ListResult) Extract() ([]Instance, error) {
	var a struct {
		Instances []Instance `json:"instances"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Instances, err
}
