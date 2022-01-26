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

// VolumeMetadata is an object that represents the metadata about the disk.
type VolumeMetadata struct {
	// Specifies the parameter that describes the encryption CMK ID in metadata.
	// This parameter is used together with __system__encrypted for encryption.
	// The length of cmkid is fixed at 36 bytes.
	SystemCmkID string `json:"__system__cmkid"`
	// Specifies the parameter that describes the encryption function in metadata. The value can be 0 or 1.
	//   0: indicates the disk is not encrypted.
	//   1: indicates the disk is encrypted.
	//   If this parameter does not appear, the disk is not encrypted by default.
	SystemEncrypted string `json:"__system__encrypted"`
	// Specifies the clone method. When the disk is created from a snapshot,
	// the parameter value is 0, indicating the linked cloning method.
	FullClone string `json:"full_clone"`
	// Specifies the parameter that describes the disk device type in metadata. The value can be true or false.
	//   If this parameter is set to true, the disk device type is SCSI, that is, Small Computer System
	//     Interface (SCSI), which allows ECS OSs to directly access the underlying storage media and supports SCSI
	//     reservation commands.
	//   If this parameter is set to false, the disk device type is VBD (the default type),
	//     that is, Virtual Block Device (VBD), which supports only simple SCSI read/write commands.
	//   If this parameter does not appear, the disk device type is VBD.
	HwPassthrough string `json:"hw:passthrough"`
	// Specifies the parameter that describes the disk billing mode in metadata.
	// If this parameter is specified, the disk is billed on a yearly/monthly basis.
	// If this parameter is not specified, the disk is billed on a pay-per-use basis.
	OrderID string `json:"orderID"`
	// Specifies the resource type about the disk.
	ResourceType string `json:"resourceType"`
	// Specifies the special code about the disk.
	ResourceSpecCode string `json:"resourceSpecCode"`
	// Specifies whether disk is read-only.
	ReadOnly string `json:"readonly"`
	// Specifies the attached mode about the disk.
	AttachedMode string `json:"attached_mode"`
}

// Link is an object that represents a link to which the disk belongs.
type Link struct {
	// Specifies the corresponding shortcut link.
	Href string `json:"href"`
	// Specifies the shortcut link marker name.
	Rel string `json:"rel"`
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
	// Specifies the disk URI.
	Links []Link `json:"links"`
	// The metadata of the disk image.
	ImageMetadata map[string]string `json:"volume_image_metadata"`
	// The ID of the snapshot from which the volume was created
	SnapshotID string `json:"snapshot_id"`
	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `json:"source_volid"`
	// Specifies the ID of the tenant to which the disk belongs. The tenant ID is actually the project ID.
	OsVolTenantAttrTenantID string `json:"os-vol-tenant-attr:tenant_id"`
	// Specifies the service type. The value can be EVS, DSS or DESS.
	ServiceType string `json:"service_type"`
	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable"`
	// Multiattach denotes if the volume is multi-attach capable.
	Multiattach bool `json:"multiattach"`
	// Specifies the ID of the DSS storage pool accommodating the disk.
	DedicatedStorageID string `json:"dedicated_storage_id"`
	// Specifies the name of the DSS storage pool accommodating the disk.
	DedicatedStorageName string `json:"dedicated_storage_name"`
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
	Metadata VolumeMetadata `json:"metadata"`
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
	pagination.OffsetPageBase
}

// IsEmpty returns true if a ListResult contains no Volumes.
func (r VolumePage) IsEmpty() (bool, error) {
	volumes, err := ExtractVolumes(r)
	return len(volumes) == 0, err
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
