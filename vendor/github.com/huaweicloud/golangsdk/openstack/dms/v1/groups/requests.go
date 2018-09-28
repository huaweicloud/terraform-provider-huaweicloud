package groups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CreateOpsBuilder is used for creating group parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToGroupCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	// Indicates the informations of a consumer group.
	Groups []GroupOps `json:"groups" required:"true"`
}

// GroupOps is referred by CreateOps
type GroupOps struct {
	// Indicates the name of a consumer group.
	// A string of 1 to 32 characters that contain
	// a-z, A-Z, 0-9, hyphens (-), and underscores (_).
	Name string `json:"name" required:"true"`
}

// ToGroupCreateMap is used for type convert
func (ops CreateOps) ToGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a group with given parameters.
func Create(client *golangsdk.ServiceClient, queueID string, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client, queueID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})

	return
}

// Delete a group by id
func Delete(client *golangsdk.ServiceClient, queueID string, groupID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, queueID, groupID), nil)
	return
}

// List all the groups
func List(client *golangsdk.ServiceClient, queueID string, includeDeadLetter bool) pagination.Pager {
	url := listURL(client, queueID, includeDeadLetter)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return GroupPage{pagination.SinglePageBase(r)}
	})
}
