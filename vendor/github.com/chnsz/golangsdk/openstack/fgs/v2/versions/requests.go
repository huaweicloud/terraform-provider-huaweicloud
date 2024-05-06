package versions

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure that used to create a release version under specified function.
type CreateOpts struct {
	// Function URN to which the version belongs.
	FunctionUrn string `json:"-" required:"true"`
	// The MD5 value.
	Digest string `json:"digest,omitempty"`
	// The name of the release version.
	Version string `json:"version,omitempty"`
	// The description of the release version.
	Description string `json:"description,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create version using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*Version, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Version
	_, err = client.Post(rootURL(client, opts.FunctionUrn), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts is the structure that used to query function version list.
type ListOpts struct {
	// Function URN.
	FunctionUrn string `json:"-" required:"true"`
	// The current query index.
	Marker int `q:"marker"`
	// Maximum number of functions to obtain in a request.
	MaxItems int `q:"maxitems"`
}

// List is a method to query the list of the function versions using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Version, error) {
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url := rootURL(client, opts.FunctionUrn) + query.String()
	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := VersionPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	pageInfo, err := extractPageInfo(pages)
	if err != nil {
		return nil, err
	}
	return pageInfo.Versions, nil
}
