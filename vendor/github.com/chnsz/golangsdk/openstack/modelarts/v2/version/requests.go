package version

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	ClearHardProperty        *bool  `json:"clear_hard_property,omitempty"`
	Description              string `json:"description,omitempty"`
	ExportImages             *bool  `json:"export_images,omitempty"`
	RemoveSampleUsage        *bool  `json:"remove_sample_usage,omitempty"`
	TrainEvaluateSampleRatio string `json:"train_evaluate_sample_ratio,omitempty"`
	VersionFormat            string `json:"version_format,omitempty"`
	VersionName              string `json:"version_name,omitempty"`
	WithColumnHeader         *bool  `json:"with_column_header,omitempty"`
}

type ListOpts struct {
	Status             int    `q:"status"`
	TrainEvaluateRatio string `q:"train_evaluate_ratio"`
	VersionFormat      int    `q:"version_format"`
	Offset             int    `q:"offset"`
	Limit              int    `q:"limit"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, datesetId string, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CreateResp
	_, err = c.Post(createURL(c, datesetId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst, err
}

func Get(c *golangsdk.ServiceClient, datesetId string, versionId string) (*DatasetVersion, error) {
	var rst DatasetVersion
	_, err := c.Get(getURL(c, datesetId, versionId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Delete(c *golangsdk.ServiceClient, datesetId string, versionId string) *golangsdk.ErrResult {
	url := deleteURL(c, datesetId, versionId)
	var rst golangsdk.ErrResult
	_, rst.Err = c.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return &rst
}

func List(c *golangsdk.ServiceClient, datesetId string, opts ListOpts) (*pagination.Pager, error) {
	url := listURL(c, datesetId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	page := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := DatasetVersionPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})

	return &page, nil
}
