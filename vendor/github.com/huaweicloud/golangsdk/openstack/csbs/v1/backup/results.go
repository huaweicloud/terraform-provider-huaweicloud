package backup

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Checkpoint struct {
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"-"`
	Id             string         `json:"id"`
	ResourceGraph  string         `json:"resource_graph"`
	ProjectId      string         `json:"project_id"`
	ProtectionPlan ProtectionPlan `json:"protection_plan"`
}

type ProtectionPlan struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	BackupResources []BackupResource `json:"resources"`
}

type BackupResource struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	ExtraInfo string `json:"-"`
}

type ResourceCapability struct {
	Result       bool   `json:"result"`
	ResourceType string `json:"resource_type"`
	ErrorCode    string `json:"error_code"`
	ErrorMsg     string `json:"error_msg"`
	ResourceId   string `json:"resource_id"`
}

// UnmarshalJSON helps to unmarshal Checkpoint fields into needed values.
func (r *Checkpoint) UnmarshalJSON(b []byte) error {
	type tmp Checkpoint
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Checkpoint(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)

	return err
}

// UnmarshalJSON helps to unmarshal BackupResource fields into needed values.
func (r *BackupResource) UnmarshalJSON(b []byte) error {
	type tmp BackupResource
	var s struct {
		tmp
		ExtraInfo interface{} `json:"extra_info"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = BackupResource(s.tmp)

	switch t := s.ExtraInfo.(type) {
	case float64:
		r.ID = strconv.FormatFloat(t, 'f', -1, 64)
	case string:
		r.ID = t
	}

	return err
}

func (r commonResult) ExtractQueryResponse() ([]ResourceCapability, error) {
	var s struct {
		ResourcesCaps []ResourceCapability `json:"protectable"`
	}
	err := r.ExtractInto(&s)
	return s.ResourcesCaps, err
}

type Backup struct {
	CheckpointId string        `json:"checkpoint_id"`
	CreatedAt    time.Time     `json:"-"`
	ExtendInfo   ExtendInfo    `json:"extend_info"`
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	ResourceId   string        `json:"resource_id"`
	Status       string        `json:"status"`
	UpdatedAt    time.Time     `json:"-"`
	VMMetadata   VMMetadata    `json:"backup_data"`
	Description  string        `json:"description"`
	Tags         []ResourceTag `json:"tags"`
	ResourceType string        `json:"resource_type"`
}

type ExtendInfo struct {
	AutoTrigger          bool           `json:"auto_trigger"`
	AverageSpeed         int            `json:"average_speed"`
	CopyFrom             string         `json:"copy_from"`
	CopyStatus           string         `json:"copy_status"`
	FailCode             FailCode       `json:"fail_code"`
	FailOp               string         `json:"fail_op"`
	FailReason           string         `json:"fail_reason"`
	ImageType            string         `json:"image_type"`
	Incremental          bool           `json:"incremental"`
	Progress             int            `json:"progress"`
	ResourceAz           string         `json:"resource_az"`
	ResourceName         string         `json:"resource_name"`
	ResourceType         string         `json:"resource_type"`
	Size                 int            `json:"size"`
	SpaceSavingRatio     int            `json:"space_saving_ratio"`
	VolumeBackups        []VolumeBackup `json:"volume_backups"`
	FinishedAt           time.Time      `json:"-"`
	TaskId               string         `json:"taskid"`
	HypervisorType       string         `json:"hypervisor_type"`
	SupportedRestoreMode string         `json:"supported_restore_mode"`
	Supportlld           bool           `json:"support_lld"`
}

type VMMetadata struct {
	RegionName       string `json:"__openstack_region_name"`
	CloudServiceType string `json:"cloudservicetype"`
	Disk             int    `json:"disk"`
	ImageType        string `json:"imagetype"`
	Ram              int    `json:"ram"`
	Vcpus            int    `json:"vcpus"`
	Eip              string `json:"eip"`
	PrivateIp        string `json:"private_ip"`
}

type FailCode struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type VolumeBackup struct {
	AverageSpeed     int    `json:"average_speed"`
	Bootable         bool   `json:"bootable"`
	Id               string `json:"id"`
	ImageType        string `json:"image_type"`
	Incremental      bool   `json:"incremental"`
	SnapshotID       string `json:"snapshot_id"`
	Name             string `json:"name"`
	Size             int    `json:"size"`
	SourceVolumeId   string `json:"source_volume_id"`
	SourceVolumeSize int    `json:"source_volume_size"`
	SpaceSavingRatio int    `json:"space_saving_ratio"`
	Status           string `json:"status"`
	SourceVolumeName string `json:"source_volume_name"`
}

// UnmarshalJSON helps to unmarshal Backup fields into needed values.
func (r *Backup) UnmarshalJSON(b []byte) error {
	type tmp Backup
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt golangsdk.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Backup(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

// UnmarshalJSON helps to unmarshal ExtendInfo fields into needed values.
func (r *ExtendInfo) UnmarshalJSON(b []byte) error {
	type tmp ExtendInfo
	var s struct {
		tmp
		FinishedAt golangsdk.JSONRFC3339MilliNoZ `json:"finished_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ExtendInfo(s.tmp)

	r.FinishedAt = time.Time(s.FinishedAt)

	return err
}

// Extract will get the checkpoint object from the commonResult
func (r commonResult) Extract() (*Checkpoint, error) {
	var s struct {
		Checkpoint *Checkpoint `json:"checkpoint"`
	}

	err := r.ExtractInto(&s)
	return s.Checkpoint, err
}

// ExtractBackup will get the backup object from the commonResult
func (r commonResult) ExtractBackup() (*Backup, error) {
	var s struct {
		Backup *Backup `json:"checkpoint_item"`
	}

	err := r.ExtractInto(&s)
	return s.Backup, err
}

// BackupPage is the page returned by a pager when traversing over a
// collection of backups.
type BackupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of backups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BackupPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"checkpoint_items_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a BackupPage struct is empty.
func (r BackupPage) IsEmpty() (bool, error) {
	is, err := ExtractBackups(r)
	return len(is) == 0, err
}

// ExtractBackups accepts a Page struct, specifically a BackupPage struct,
// and extracts the elements into a slice of Backup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBackups(r pagination.Page) ([]Backup, error) {
	var s struct {
		Backups []Backup `json:"checkpoint_items"`
	}
	err := (r.(BackupPage)).ExtractInto(&s)
	return s.Backups, err
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type QueryResult struct {
	commonResult
}
