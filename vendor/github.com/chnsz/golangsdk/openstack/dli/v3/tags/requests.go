package tags

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// UpdateTagsOpts is the structure that used to update tags.
type UpdateTagsOpts struct {
	// Specifies the action to update the tags. The value is create or delete.
	Action string `json:"-" required:"true"`
	// Specifies resource type.
	// The valid values are as follows:
	// + dli_package_resource
	// + dli_flink_job
	ResourceType string `json:"-" required:"true"`
	// Specifies the key/value pairs of the DLI package.
	Tags []tags.ResourceTag `json:"tags" required:"true"`
}

// UpdateTags is a method that used to update tags to specifies resource.
func UpdateTagsToResource(c *golangsdk.ServiceClient, id string, opts UpdateTagsOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(tagsURL(c, id, opts.ResourceType, opts.Action), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return err
}
