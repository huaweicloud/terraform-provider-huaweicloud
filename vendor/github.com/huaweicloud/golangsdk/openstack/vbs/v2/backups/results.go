package backups

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type Backup struct {
	//Backup ID
	Id string `json:"id"`
	//Backup name
	Name string `json:"name"`
	//Backup URL
	Links []golangsdk.Link `json:"links"`
	//Backup status
	Status string `json:"status"`
	//Backup description
	Description string `json:"description"`
	//AvailabilityZone where the backup resides
	AvailabilityZone string `json:"availability_zone"`
	//Source volume ID of the backup
	VolumeId string `json:"volume_id"`
	//Cause of the backup failure
	FailReason string `json:"fail_reason"`
	//Backup size
	Size int `json:"size"`
	//Number of objects on OBS for the disk data
	ObjectCount int `json:"object_count"`
	//Container of the backup
	Container string `json:"container"`
	//Backup creation time
	CreatedAt time.Time `json:"-"`
	//ID of the tenant to which the backup belongs
	TenantId string `json:"os-bak-tenant-attr:tenant_id"`
	//Backup metadata
	ServiceMetadata string `json:"service_metadata"`
	//Time when the backup was updated
	UpdatedAt time.Time `json:"-"`
	//Current time
	DataTimeStamp time.Time `json:"-"`
	//Whether a dependent backup exists
	DependentBackups bool `json:"has_dependent_backups"`
	//ID of the snapshot associated with the backup
	SnapshotId string `json:"snapshot_id"`
	//Whether the backup is an incremental backup
	Incremental bool `json:"is_incremental"`
}

type BackupRestoreInfo struct {
	//Backup ID
	BackupId string `json:"backup_id"`
	//Volume ID
	VolumeId string `json:"volume_id"`
	//Volume name
	VolumeName string `json:"volume_name"`
}

// BackupPage is the page returned by a pager when traversing over a
// collection of Backups.
type BackupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of Backups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BackupPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"backups_links"`
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
// and extracts the elements into a slice of Backups struct. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBackups(r pagination.Page) ([]Backup, error) {
	var s struct {
		Backups []Backup `json:"backups"`
	}
	err := (r.(BackupPage)).ExtractInto(&s)
	return s.Backups, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a backup.
func (r commonResult) Extract() (*Backup, error) {
	var s struct {
		Backup *Backup `json:"backup"`
	}
	err := r.ExtractInto(&s)
	return s.Backup, err
}

// ExtractBackupRestore is a function that accepts a result and extracts a backup
func (r commonResult) ExtractBackupRestore() (*BackupRestoreInfo, error) {
	var s struct {
		Restore *BackupRestoreInfo `json:"restore"`
	}
	err := r.ExtractInto(&s)
	return s.Restore, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Backup.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Backup.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}

// UnmarshalJSON overrides the default, to convert the JSON API response into our Backup struct
func (r *Backup) UnmarshalJSON(b []byte) error {
	type tmp Backup
	var s struct {
		tmp
		CreatedAt     golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt     golangsdk.JSONRFC3339MilliNoZ `json:"updated_at"`
		DataTimeStamp golangsdk.JSONRFC3339MilliNoZ `json:"data_timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Backup(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	r.DataTimeStamp = time.Time(s.DataTimeStamp)

	return nil
}
