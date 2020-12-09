package cloudimages

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Image represents an image found in the IMS.
type Image struct {
	// the URL for uploading and downloading the image file
	File string `json:"file"`
	// the image owner
	Owner string `json:"owner"`
	// the image id
	ID string `json:"id"`
	// the image URL
	Self string `json:"self"`
	// the image schema
	Schema string `json:"schema"`
	// the image status, the value can be [queued, saving, deleted, killed,active]
	Status string `json:"status"`
	// the image tags
	Tags []string `json:"tags"`
	// whether the image can be seen by others
	Visibility string `json:"visibility"`
	// the image name
	Name string `json:"name"`
	// whether the image has been deleted
	Deleted bool `json:"deleted"`
	// whether the image is protected
	Protected bool `json:"protected"`
	// the container type
	ContainerFormat string `json:"container_format"`
	// the minimum memory size (MB) required for running the image
	MinRam int `json:"min_ram"`
	// the maximum memory of the image in the unit of MB, notice: string
	MaxRam string `json:"max_ram"`
	// the disk format, the value can be [vhd, raw, zvhd, qcow2]
	DiskFormat string `json:"disk_format"`
	// the minimum disk space (GB) required for running the image
	MinDisk int `json:"min_disk"`
	// the environment where the image is used
	VirtualEnvType string `json:"virtual_env_type"`
	// *size, virtual_size and checksum parameter are unavailable currently*
	Size        int64  `json:"size"`
	VirtualSize int    `json:"virtual_size"`
	Checksum    string `json:"checksum"`
	// created_at and updated_at are in UTC format
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt string    `json:"deleted_at"`
	// the OS architecture: 32 or 64
	OsBit                  string `json:"__os_bit"`
	OsVersion              string `json:"__os_version"`
	Description            string `json:"__description"`
	OsType                 string `json:"__os_type"`
	Isregistered           string `json:"__isregistered"`
	Platform               string `json:"__platform"`
	ImageSourceType        string `json:"__image_source_type"`
	Imagetype              string `json:"__imagetype"`
	Originalimagename      string `json:"__originalimagename"`
	BackupID               string `json:"__backup_id"`
	Productcode            string `json:"__productcode"`
	ImageSize              string `json:"__image_size"`
	DataOrigin             string `json:"__data_origin"`
	SupportKvm             string `json:"__support_kvm"`
	SupportXen             string `json:"__support_xen"`
	SupportLargeMemory     string `json:"__support_largememory"`
	SupportDiskintensive   string `json:"__support_diskintensive"`
	SupportHighperformance string `json:"__support_highperformance"`
	SupportXenGpuType      string `json:"__support_xen_gpu_type"`
	SupportKvmGpuType      string `json:"__support_kvm_gpu_type"`
	SupportXenHana         string `json:"__support_xen_hana"`
	SupportKvmInfiniband   string `json:"__support_kvm_infiniband"`
	SystemSupportMarket    bool   `json:"__system_support_market"`
	RootOrigin             string `json:"__root_origin"`
	SequenceNum            string `json:"__sequence_num"`
}

func (r *Image) UnmarshalJSON(b []byte) error {
	type tmp Image
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339Milli `json:"created_at"`
		UpdatedAt golangsdk.JSONRFC3339Milli `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Image(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

type commonResult struct {
	golangsdk.Result
}

// Extract will get the Image object out of the commonResult object.
func (r commonResult) Extract() (*Image, error) {
	var s Image
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a volume struct
func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "images")
}

// ImagePage represents the results of a List request.
type ImagePage struct {
	serviceURL string
	pagination.LinkedPageBase
}

// IsEmpty returns true if an ImagePage contains no Images results.
func (r ImagePage) IsEmpty() (bool, error) {
	images, err := ExtractImages(r)
	return len(images) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to
// the next page of results.
func (r ImagePage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	if s.Next == "" {
		return "", nil
	}

	return nextPageURL(r.serviceURL, s.Next)
}

// ExtractImages interprets the results of a single page from a List() call,
// producing a slice of Image entities.
func ExtractImages(r pagination.Page) ([]Image, error) {
	var s struct {
		Images []Image `json:"images"`
	}

	err := (r.(ImagePage)).ExtractInto(&s)
	return s.Images, err
}
