package resourcetags

import (
	"github.com/chnsz/golangsdk"
)

// BatchOpts is the structure that used to manage tags.
type BatchOpts struct {
	// Specifies the resource list.
	Resources []Resource `json:"resources" required:"true"`
	// Tags list.
	Tags []ResourceTag `json:"tags" required:"true"`
	// Specifies the project ID. This parameter is mandatory when resource_type is a region-specific service.
	ProjectId string `json:"project_id,omitempty"`
}

// Resource is the object that represents the managed resource configuration.
type Resource struct {
	// Specifies the resource type.
	ResourceType string `json:"resource_type" required:"true"`
	// Specifies the resource ID.
	ResourceId string `json:"resource_id" required:"true"`
}

// ResourceTag is the object that represents the tags configuration for batch management.
type ResourceTag struct {
	// Specifies the tag key.
	// The value can contain up to 36 characters including letters, digits, hyphens (-), and underscores (_).
	Key string `json:"key" required:"true"`
	// Specifies the tag value.
	// The value can contain up to 43 characters including letters, digits, periods (.), hyphens (-) and
	// underscores (_). It can be an empty string.
	Value *string `json:"value,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create tags in batch using given parameters.
func Create(client *golangsdk.ServiceClient, opts BatchOpts) ([]FailedResource, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r struct {
		FailedResources []FailedResource `json:"failed_resources"`
	}
	_, err = client.Post(createURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.FailedResources, err
}

// QueryOpts is the structure that used to query tags detail for specified resource.
type QueryOpts struct {
	// Resource ID to be queried.
	ResourceId string `json:"resource_id" required:"true"`
	// Resource type to be queried.
	ResourceType string `q:"resource_type" required:"true"`
	// The project ID to which the managed resources belong.
	ProjectId string `q:"project_id"`
}

// Get is a method to obtain the tags details using given parameters.
func Get(client *golangsdk.ServiceClient, opts QueryOpts) ([]ResourceTag, error) {
	url := queryURL(client, opts.ResourceId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()
	var r struct {
		Tags []ResourceTag `json:"tags"`
	}
	_, err = client.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Tags, err
}

// Delete is a method to delete tags in batch using given parameters.
func Delete(client *golangsdk.ServiceClient, opts BatchOpts) ([]FailedResource, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r struct {
		FailedResources []FailedResource `json:"failed_resources"`
	}
	_, err = client.Post(deleteURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.FailedResources, err
}
