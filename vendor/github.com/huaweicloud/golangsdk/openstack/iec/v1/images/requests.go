package images

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
//
// http://developer.openstack.org/api-ref-image-v2.html
type ListOpts struct {
	//Image type
	ImageType string `q:"__imagetype"`

	Protected string `q:"protected"`

	// Visibility filters on the visibility of the image.
	Visibility string `q:"visibility"`

	// Status filters on the status of the image.
	Status string `q:"status"`

	// Name filters on the name of the image.
	Name string `q:"name"`

	//Indicates the image OS type. The value can be Linux, Windows, or Other.
	OsType string `q:"__os_type"`

	//
	VirtualEnvType string `q:"virtual_env_type"`

	// Specifies whether the image is available.
	IsRegistered string `q:"__isregistered"`

	// Integer value for the limit of values to return.
	Limit int `q:"limit"`

	Offset string `q:"offset"`

	// SortKey will sort the results based on a specified image property.
	SortKey string `q:"sort_key"`

	// SortDir will sort the list results either ascending or decending.
	SortDir string `q:"sort_dir"`

	// Owner filters on the project ID of the image.
	Owner string `q:"owner"`

	//Image ID
	ID string `q:"id"`

	//Specifies whether the image supports KVM.
	//If yes, the value is true. Otherwise, this attribute is not required.
	SupportKvm string `q:"__support_kvm"`

	//If the image supports the GPU type on the KVM virtualization platform,
	//the value is V100_vGPU or RTX5000. Otherwise, this attribute is unavailable.
	SupportKvmGpuType string `q:"__support_kvm_gpu_type"`

	//Specifies whether the image supports AI acceleration.
	SupportKvmAscend310 string `q:"__support_kvm_ascend_310"`

	//Specifies whether the image supports Computing enhancement
	SupportKvmHi1822Hiovs string `q:"__support_kvm_hi1822_hiovs"`

	//Specifies whether the image supports ARM architecture
	SupportArm string `q:"__support_arm"`

	//HwFirmwareType firmware type
	HwFirmwareType string `q:"hw_firmware_type"`
}

// ToImageListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImageListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List implements image list request.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := ListURL(c)
	if opts != nil {
		query, err := opts.ToImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}
