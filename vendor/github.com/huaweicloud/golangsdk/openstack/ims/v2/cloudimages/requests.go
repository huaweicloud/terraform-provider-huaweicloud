package cloudimages

import (
	"fmt"
	"net/url"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

type ListOpts struct {
	Isregistered           string `q:"__isregistered"`
	Imagetype              string `q:"__imagetype"`
	Protected              bool   `q:"protected"`
	Visibility             string `q:"visibility"`
	Owner                  string `q:"owner"`
	ID                     string `q:"id"`
	Status                 string `q:"status"`
	Name                   string `q:"name"`
	ContainerFormat        string `q:"container_format"`
	DiskFormat             string `q:"disk_format"`
	MinRam                 int    `q:"min_ram"`
	MinDisk                int    `q:"min_disk"`
	OsBit                  string `q:"__os_bit"`
	Platform               string `q:"__platform"`
	Marker                 string `q:"marker"`
	Limit                  int    `q:"limit"`
	SortKey                string `q:"sort_key"`
	SortDir                string `q:"sort_dir"`
	OsType                 string `q:"__os_type"`
	Tag                    string `q:"tag"`
	MemberStatus           string `q:"member_status"`
	SupportKvm             string `q:"__support_kvm"`
	SupportXen             string `q:"__support_xen"`
	SupportLargeMemory     string `q:"__support_largememory"`
	SupportDiskintensive   string `q:"__support_diskintensive"`
	SupportHighperformance string `q:"__support_highperformance"`
	SupportXenGpuType      string `q:"__support_xen_gpu_type"`
	SupportKvmGpuType      string `q:"__support_kvm_gpu_type"`
	SupportXenHana         string `q:"__support_xen_hana"`
	SupportKvmInfiniband   string `q:"__support_kvm_infiniband"`
	VirtualEnvType         string `q:"virtual_env_type"`
	// CreatedAtQuery filters images based on their creation date.
	CreatedAtQuery *ImageDateQuery
	// UpdatedAtQuery filters images based on their updated date.
	UpdatedAtQuery *ImageDateQuery
}

// ImageDateFilter represents a valid filter to use for filtering
// images by their date during a List.
type ImageDateFilter string

const (
	FilterGT  ImageDateFilter = "gt"
	FilterGTE ImageDateFilter = "gte"
	FilterLT  ImageDateFilter = "lt"
	FilterLTE ImageDateFilter = "lte"
	FilterNEQ ImageDateFilter = "neq"
	FilterEQ  ImageDateFilter = "eq"
)

// ImageDateQuery represents a date field to be used for listing images.
// If no filter is specified, the query will act as though FilterEQ was
// set.
type ImageDateQuery struct {
	Date   time.Time
	Filter ImageDateFilter
}

// ToImageListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImageListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	params := q.Query()

	if opts.CreatedAtQuery != nil {
		createdAt := opts.CreatedAtQuery.Date.Format(time.RFC3339)
		if v := opts.CreatedAtQuery.Filter; v != "" {
			createdAt = fmt.Sprintf("%s:%s", v, createdAt)
		}
		params.Add("created_at", createdAt)
	}

	if opts.UpdatedAtQuery != nil {
		updatedAt := opts.UpdatedAtQuery.Date.Format(time.RFC3339)
		if v := opts.UpdatedAtQuery.Filter; v != "" {
			updatedAt = fmt.Sprintf("%s:%s", v, updatedAt)
		}
		params.Add("updated_at", updatedAt)
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

// List implements images list request
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		imagePage := ImagePage{
			serviceURL:     client.ServiceURL(),
			LinkedPageBase: pagination.LinkedPageBase{PageResult: r},
		}

		return imagePage
	})
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	// Returns value that can be passed to json.Marshal
	ToImageCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create an image.
type CreateByServerOpts struct {
	// the name of the system disk image
	Name string `json:"name" required:"true"`
	// Description of the image
	Description string `json:"description,omitempty"`
	// server id to be converted
	InstanceId string `json:"instance_id" required:"true"`
	// the data disks to be converted
	DataImages []DataImage `json:"data_images,omitempty"`
	// image label "key.value"
	Tags []string `json:"tags,omitempty"`
	// One or more tag key and value pairs to associate with the image
	ImageTags []ImageTag `json:"image_tags,omitempty"`
	// the maximum memory of the image in the unit of MB
	MaxRam int `json:"max_ram,omitempty"`
	// the minimum memory of the image in the unit of MB
	MinRam int `json:"min_ram,omitempty"`
}

// CreateOpts represents options used to create an image.
type CreateByOBSOpts struct {
	// the name of the system disk image
	Name string `json:"name" required:"true"`
	// Description of image
	Description string `json:"description,omitempty"`
	// the OS version
	OsVersion string `json:"os_version,omitempty"`
	// the URL of the external image file in the OBS bucket
	ImageUrl string `json:"image_url" required:"true"`
	// the minimum size of the system disk in the unit of GB
	MinDisk int `json:"min_disk" required:"true"`
	//whether automatic configuration is enabled，the value can be true or false
	IsConfig bool `json:"is_config,omitempty"`
	// the master key used for encrypting an image
	CmkId string `json:"cmk_id,omitempty"`
	// image label "key.value"
	Tags []string `json:"tags,omitempty"`
	// One or more tag key and value pairs to associate with the image
	ImageTags []ImageTag `json:"image_tags,omitempty"`
	// the image type, the value can be ECS,BMS,FusionCompute, or Ironic
	Type string `json:"type,omitempty"`
	// the maximum memory of the image in the unit of MB
	MaxRam int `json:"max_ram,omitempty"`
	// the minimum memory of the image in the unit of MB
	MinRam int `json:"min_ram,omitempty"`
}

// CreateOpts represents options used to create an image.
type CreateDataImageByServerOpts struct {
	// the data disks to be converted
	DataImages []DataImage `json:"data_images" required:"true"`
}

// CreateOpts represents options used to create an image.
type CreateDataImageByOBSOpts struct {
	// the name of the data disk image
	Name string `json:"name" required:"true"`
	// Description of image
	Description string `json:"description,omitempty"`
	// the OS type
	OsType string `json:"os_type" required:"true"`
	// the URL of the external image file in the OBS bucket
	ImageUrl string `json:"image_url" required:"true"`
	// the minimum size of the system disk in the unit of GB
	MinDisk int `json:"min_disk" required:"true"`
	// the master key used for encrypting an image
	CmkId string `json:"cmk_id,omitempty"`
}

type DataImage struct {
	// the data disk image name
	Name string `json:"name" required:"true"`
	// the data disk ID
	VolumeId string `json:"volume_id" required:"true"`
	// information about the data disk
	Description string `json:"description,omitempty"`
	// the data disk image tags
	Tags []string `json:"tags,omitempty"`
}

type ImageTag struct {
	Key   string `json:"key" required:"true"`
	Value string `json:"value,omitempty"`
}

// ToImageCreateMap assembles a request body based on the contents of
// a CreateByServerOpts.
func (opts CreateByServerOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts CreateByOBSOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts CreateDataImageByServerOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts CreateDataImageByOBSOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create implements create image request.
func CreateImageByServer(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Create implements create image request.
func CreateImageByOBS(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Create implements create image request.
func CreateDataImageByServer(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}

// Create implements create image request.
func CreateDataImageByOBS(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createDataImageURL(client), b, &r.Body, &golangsdk.RequestOpts{OkCodes: []int{200}})
	return
}
