package cloudimages

import (
	"fmt"
	"net/url"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

type ListOpts struct {
	Isregistered           string `q:"__isregistered"`
	Imagetype              string `q:"__imagetype"`
	WholeImage             bool   `q:"__whole_image"`
	SystemCmkid            string `q:"__system__cmkid"`
	Protected              bool   `q:"protected"`
	Visibility             string `q:"visibility"`
	Owner                  string `q:"owner"`
	ID                     string `q:"id"`
	Status                 string `q:"status"`
	Name                   string `q:"name"`
	FlavorId               string `q:"flavor_id"`
	ContainerFormat        string `q:"container_format"`
	DiskFormat             string `q:"disk_format"`
	MinRam                 int    `q:"min_ram"`
	MinDisk                int    `q:"min_disk"`
	Marker                 string `q:"marker"`
	Limit                  int    `q:"limit"`
	SortKey                string `q:"sort_key"`
	SortDir                string `q:"sort_dir"`
	OsType                 string `q:"__os_type"`
	Platform               string `q:"__platform"`
	OsVersion              string `q:"__os_version"`
	OsBit                  string `q:"__os_bit"`
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
	Architecture           string `q:"architecture"`
	EnterpriseProjectID    string `q:"enterprise_project_id"`
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
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

// CreateOpts represents options used to create an image.
type CreateByOBSOpts struct {
	// the name of the system disk image
	Name string `json:"name" required:"true"`
	// Description of image
	Description string `json:"description,omitempty"`
	// The OS type of the image
	OsType string `json:"os_type,omitempty"`
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
	// Whether to use the image file quick import method to create an image
	IsQuickImport bool `json:"is_quick_import,omitempty"`
	// The schema type of the image
	Architecture string `json:"architecture,omitempty"`
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

// CreateWholeImageOpts represents options used to create an image.
type CreateWholeImageOpts struct {
	// the name of the system disk image
	Name string `json:"name" required:"true"`
	// Description of image
	Description string `json:"description,omitempty"`
	// the ID of the instance
	InstanceId string `json:"instance_id,omitempty"`
	// the ID of the CBR backup
	BackupId string `json:"backup_id,omitempty"`
	// image label "key.value"
	Tags []string `json:"tags,omitempty"`
	// One or more tag key and value pairs to associate with the image
	ImageTags []ImageTag `json:"image_tags,omitempty"`
	// the maximum memory of the image in the unit of MB
	MaxRam int `json:"max_ram,omitempty"`
	// the minimum memory of the image in the unit of MB
	MinRam int `json:"min_ram,omitempty"`
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
	// the ID of the vault to which an ECS is to be added or has been added
	VaultId string `json:"vault_id,omitempty"`
	// the method of creating a full-ECS image
	WholeImageType string `json:"whole_image_type,omitempty"`
}

// CreateOpts represents options used to create an image.
type CreateDataImageByServerOpts struct {
	// the data disks to be converted
	DataImages []DataImage `json:"data_images" required:"true"`
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
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
	// One or more tag key and value pairs to associate with the image
	ImageTags []ImageTag `json:"image_tags,omitempty"`
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

// CreateSystemImageByVolumeOpts is the structure used to create a system image from EVS volume.
type CreateSystemImageByVolumeOpts struct {
	// The name of the system image.
	Name string `json:"name" required:"true"`
	// The data disk ID.
	VolumeId string `json:"volume_id" required:"true"`
	// The operating system version.
	OsVersion string `json:"os_version,omitempty"`
	// The image type, the value can be **ECS**, **FusionCompute**, **BMS**, or **Ironic**.
	Type string `json:"type,omitempty"`
	// The description of the image.
	Description string `json:"description,omitempty"`
	// The minimum memory of the image, in MB unit.
	MinRam int `json:"min_ram,omitempty"`
	// The maximum memory of the image, in MB unit.
	MaxRam int `json:"max_ram,omitempty"`
	// The image label list, **key.value** format.
	Tags []string `json:"tags,omitempty"`
	// One or more tag key and value pairs to associate with the image.
	ImageTags []ImageTag `json:"image_tags,omitempty"`
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
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

func (opts CreateWholeImageOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ToImageCreateMap assembles a request body based on the contents of the CreateSystemImageByVolumeOpts.
func (opts CreateSystemImageByVolumeOpts) ToImageCreateMap() (map[string]interface{}, error) {
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

// Create implements create whole image request.
func CreateWholeImageByServer(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createWholeImageURL(client), b, &r.Body, nil)
	return
}

// Create implements create image request.
func CreateWholeImageByBackup(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createWholeImageURL(client), b, &r.Body, nil)
	return
}

// Update implements image updated request.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToImageUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}
	_, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/openstack-images-v2.1-json-patch"},
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToImageUpdateMap() ([]interface{}, error)
}

// UpdateOpts implements UpdateOpts
type UpdateOpts []Patch

// ToImageUpdateMap assembles a request body based on the contents of UpdateOpts.
func (opts UpdateOpts) ToImageUpdateMap() ([]interface{}, error) {
	m := make([]interface{}, len(opts))
	for i, patch := range opts {
		patchJSON := patch.ToImagePatchMap()
		m[i] = patchJSON
	}
	return m, nil
}

// Patch represents a single update to an existing image. Multiple updates
// to an image can be submitted at the same time.
type Patch interface {
	ToImagePatchMap() map[string]interface{}
}

// UpdateOp represents a valid update operation.
type UpdateOp string

const (
	AddOp     UpdateOp = "add"
	ReplaceOp UpdateOp = "replace"
	RemoveOp  UpdateOp = "remove"
)

// UpdateImageProperty represents an update property request.
type UpdateImageProperty struct {
	Op    UpdateOp
	Name  string
	Value interface{}
}

// ToImagePatchMap assembles a request body based on UpdateImageProperty.
func (r UpdateImageProperty) ToImagePatchMap() map[string]interface{} {
	updateMap := map[string]interface{}{
		"op":    r.Op,
		"path":  fmt.Sprintf("/%s", r.Name),
		"value": r.Value,
	}

	return updateMap
}
