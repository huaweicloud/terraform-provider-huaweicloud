package custom

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure used to create a new custom event source.
type CreateOpts struct {
	// The ID of the event channel to which the custom event source belongs.
	ChannelId string `json:"channel_id" required:"true"`
	// The name of the custom event source.
	Name string `json:"name" required:"true"`
	// The type of custom event source.
	// The valid values are as follows:
	// + APPLICATION (default)
	// + RABBITMQ
	// + ROCKETMQ
	Type string `json:"type,omitempty"`
	// The description of the custom event source.
	Description string `json:"description,omitempty"`
	// The configuration detail of the event source, in JSON format.
	Detail interface{} `json:"detail,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new custom event source using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*Source, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Source
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to query an existing event source by its ID.
func Get(client *golangsdk.ServiceClient, sourceId string) (*Source, error) {
	var r Source
	_, err := client.Get(resourceURL(client, sourceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts is the structure used to query event source list.
type ListOpts struct {
	// The ID of the event channel to which the event source belongs.
	ChannelId string `q:"channel_id"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0.
	// The valid value ranges from 0 to 100, defaults to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page.
	// The valid value ranges from 1 to 1000, defaults to 15.
	Limit int `q:"limit"`
	// The query sorting.
	// The default value is 'created_time:DESC'.
	Sort string `q:"sort"`
	// The type of the event source provider.
	// + OFFICIAL: official cloud service event source.
	// + CUSTOM: the user-defined event source.
	// + PARTNER: partner event source.
	ProviderType string `q:"provider_type"`
	// The name of the event source.
	Name string `q:"name"`
	// The fuzzy name of the event source.
	FuzzyName string `q:"fuzzy_name"`
	// The fuzzy label of the event source.
	FuzzyLabel string `q:"fuzzy_label"`
}

// List is a method to query the event source list using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Source, error) {
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := SourcePage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractSources(pages)
}

// UpdateOpts is the structure used to update an existing custom event source.
type UpdateOpts struct {
	// The ID of the event source.
	SourceId string `json:"-" required:"true"`
	// The description of the event source.
	Description *string `json:"description,omitempty"`
	// The configuration detail of the event source, in JSON format.
	Detail interface{} `json:"detail,omitempty"`
}

// Update is a method used to modify an existing custom event source using given parameters.
func Update(client *golangsdk.ServiceClient, opts UpdateOpts) (*Source, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Source
	_, err = client.Put(resourceURL(client, opts.SourceId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing custom event source using its ID.
func Delete(client *golangsdk.ServiceClient, sourceId string) error {
	_, err := client.Delete(resourceURL(client, sourceId), nil)
	return err
}
