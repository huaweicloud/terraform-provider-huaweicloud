package configurations

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//CreateResult is a struct that contains all the return parameters of creation
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
	Tenant         string         `json:"tenant"`
	Name           string         `json:"scaling_configuration_name"`
	InstanceConfig InstanceConfig `json:"instance_config"`
	CreateTime     string         `json:"create_time"`
}

type InstanceConfig struct {
	FlavorRef    string                 `json:"flavorRef"`
	ImageRef     string                 `json:"imageRef"`
	Disk         []Disk                 `json:"disk"`
	SSHKey       string                 `json:"key_name"`
	InstanceName string                 `json:"instance_name"`
	InstanceID   string                 `json:"instance_id"`
	AdminPass    string                 `json:"adminPass"`
	Personality  []Personality          `json:"personality"`
	PublicIp     PublicIp               `json:"public_ip"`
	UserData     string                 `json:"user_data"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type Disk struct {
	Size       int    `json:"size"`
	VolumeType string `json:"volume_type"`
	DiskType   string `json:"disk_type"`
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
	err := r.Result.ExtractIntoSlicePtr(&cs, "scaling_groups")
	return cs, err
}
