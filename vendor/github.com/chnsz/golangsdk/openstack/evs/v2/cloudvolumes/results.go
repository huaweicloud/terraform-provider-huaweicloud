package cloudvolumes

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// Attachment contains the disk attachment information
type Attachment struct {
	// Specifies the ID of the attachment information
	AttachmentID string `json:"attachment_id"`
	// Specifies the disk ID
	VolumeID string `json:"volume_id"`
	// Specifies the ID of the attached resource, equals to volume_id
	ResourceID string `json:"id"`
	// Specifies the ID of the server to which the disk is attached
	ServerID string `json:"server_id"`
	// Specifies the name of the host accommodating the server to which the disk is attached
	HostName string `json:"host_name"`
	// Specifies the device name
	Device string `json:"device"`
	// Specifies the time when the disk was attached. Time format: UTC YYYY-MM-DDTHH:MM:SS.XXXXXX
	AttachedAt string `json:"attached_at"`
}

// Volume contains all the information associated with a Volume.
type Volume struct {
	// Unique identifier for the volume.
	ID string `json:"id"`
	// Human-readable display name for the volume.
	Name string `json:"name"`
	// Current status of the volume.
	Status string `json:"status"`
	// Size of the volume in GB.
	Size int `json:"size"`
	// Human-readable description for the volume.
	Description string `json:"description"`
	// The type of volume to create, either SATA or SSD.
	VolumeType string `json:"volume_type"`
	// AvailabilityZone is which availability zone the volume is in.
	AvailabilityZone string `json:"availability_zone"`
	// Instances onto which the volume is attached.
	Attachments []Attachment `json:"attachments"`

	// The metadata of the disk image.
	ImageMetadata map[string]string `json:"volume_image_metadata"`
	// The ID of the snapshot from which the volume was created
	SnapshotID string `json:"snapshot_id"`
	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `json:"source_volid"`

	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable"`
	// Multiattach denotes if the volume is multi-attach capable.
	Multiattach bool `json:"multiattach"`
	// Encrypted denotes if the volume is encrypted.
	Encrypted bool `json:"encrypted"`
	// wwn of the volume.
	WWN string `json:"wwn"`
	// enterprise project ID bound to the volume
	EnterpriseProjectID string `json:"enterprise_project_id"`
	// ReplicationStatus is the status of replication.
	ReplicationStatus string `json:"replication_status"`
	// ConsistencyGroupID is the consistency group ID.
	ConsistencyGroupID string `json:"consistencygroup_id"`
	// Arbitrary key-value pairs defined by the metadata field table.
	Metadata map[string]string `json:"metadata"`
	// Arbitrary key-value pairs defined by the user.
	Tags map[string]string `json:"tags"`
	// UserID is the id of the user who created the volume.
	UserID string `json:"user_id"`
	// The date when this volume was created.
	CreatedAt string `json:"created_at"`
	// The date when this volume was last updated
	UpdatedAt string `json:"updated_at"`
}

// VolumePage is a pagination.pager that is returned from a call to the List function.
type VolumePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a ListResult contains no Volumes.
func (r VolumePage) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(r)
	return len(volumes) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r VolumePage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"volumes_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// ExtractVolumes extracts and returns Volumes. It is used while iterating over a cloudvolumes.List call.
func ExtractVolumes(r pagination.Page) ([]Volume, error) {
	var s []Volume
	err := extractVolumesInto(r, &s)
	return s, err
}

func extractVolumesInto(r pagination.Page, v interface{}) error {
	return r.(VolumePage).Result.ExtractIntoSlicePtr(v, "volumes")
}

type commonResult struct {
	golangsdk.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (*Volume, error) {
	var s Volume
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a volume struct
func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "volume")
}

// GetResult contains the response body from a Get request.
type GetResult struct {
	commonResult
}

// UpdateResult contains the response body from a Update request.
type UpdateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}

// ErrorInfo contains the error message returned when an error occurs
type ErrorInfo struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// JobResponse contains all the information from Create and ExtendSize response
type JobResponse struct {
	JobID     string    `json:"job_id"`
	OrderID   string    `json:"order_id"`
	VolumeIDs []string  `json:"volume_ids"`
	Error     ErrorInfo `json:"error"`
}

// JobResult contains the response body and error from Create and ExtendSize requests
type JobResult struct {
	golangsdk.Result
}

// Extract will get the JobResponse object out of the JobResult
func (r JobResult) Extract() (*JobResponse, error) {
	job := new(JobResponse)
	err := r.ExtractInto(job)
	return job, err
}
