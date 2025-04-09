package configurations

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (string, error) {
	var a struct {
		ID string `json:"scaling_configuration_id"`
	}

	err := r.Result.ExtractInto(&a)
	return a.ID, err
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (Configuration, error) {
	var a Configuration
	err := r.Result.ExtractIntoStructPtr(&a, "scaling_configuration")
	return a, err
}

type Configuration struct {
	ID             string         `json:"scaling_configuration_id"`
	Name           string         `json:"scaling_configuration_name"`
	InstanceConfig InstanceConfig `json:"instance_config"`
	ScalingGroupID string         `json:"scaling_group_id"`
	Tenant         string         `json:"tenant"`
	CreateTime     string         `json:"create_time"`
}

type InstanceConfig struct {
	FlavorRef            string `json:"flavorRef"`
	ImageRef             string `json:"imageRef"`
	SSHKey               string `json:"key_name"`
	KeyFingerprint       string `json:"key_fingerprint"`
	InstanceName         string `json:"instance_name"`
	InstanceID           string `json:"instance_id"`
	AdminPass            string `json:"adminPass"`
	ServerGroupID        string `json:"server_group_id"`
	Tenancy              string `json:"tenancy"`
	DedicatedHostID      string `json:"dedicated_host_id"`
	MarketType           string `json:"market_type"`
	FlavorPriorityPolicy string `json:"multi_flavor_priority_policy"`

	Disk           []Disk          `json:"disk"`
	PublicIp       PublicIp        `json:"public_ip"`
	SecurityGroups []SecurityGroup `json:"security_groups"`
	Personality    []Personality   `json:"personality"`

	UserData string                 `json:"user_data"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Disk struct {
	Size               int               `json:"size"`
	VolumeType         string            `json:"volume_type"`
	DiskType           string            `json:"disk_type"`
	DedicatedStorageID string            `json:"dedicated_storage_id"`
	DataDiskImageID    string            `json:"data_disk_image_id"`
	SnapshotID         string            `json:"snapshot_id"`
	Iops               int               `json:"iops"`
	Throughput         int               `json:"throughput"`
	Metadata           map[string]string `json:"metadata"`
}

type Personality struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type PublicIp struct {
	Eip Eip `json:"eip"`
}

type Eip struct {
	Type      string    `json:"ip_type"`
	Bandwidth Bandwidth `json:"bandwidth"`
}

type Bandwidth struct {
	Size         int    `json:"size"`
	ShareType    string `json:"share_type"`
	ChargingMode string `json:"charging_mode"`
	ID           string `json:"id"`
}

type SecurityGroup struct {
	ID string `json:"id"`
}
type DeleteResult struct {
	golangsdk.ErrResult
}

type ConfigurationPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Volumes.
func (r ConfigurationPage) IsEmpty() (bool, error) {
	configs, err := r.Extract()
	return len(configs) == 0, err
}

func (r ConfigurationPage) Extract() ([]Configuration, error) {
	var cs []Configuration
	err := r.Result.ExtractIntoSlicePtr(&cs, "scaling_configurations")
	return cs, err
}
