package tags

import (
	"github.com/huaweicloud/golangsdk"
)

// Tag is a structure of key value pair.
type Tag struct {
	//tag key
	Key string `json:"key" required:"true"`
	//tag value
	Value string `json:"value" required:"true"`
}

// BatchOptsBuilder allows extensions to add additional parameters to the
// BatchAction request.
type BatchOptsBuilder interface {
	ToTagsBatchMap() (map[string]interface{}, error)
}

// BatchOpts contains all the values needed to perform BatchAction on the image tags.
type BatchOpts struct {
	//List of tags to perform batch operation
	Tags []Tag `json:"tags,omitempty"`
	//Operator , Possible values are:create or delete
	Action ActionType `json:"action" required:"true"`
}

//ActionType specifies the type of batch operation action to be performed
type ActionType string

var (
	// ActionCreate is used to set action operator to create
	ActionCreate ActionType = "create"
	// ActionDelete is used to set action operator to delete
	ActionDelete ActionType = "delete"
)

// ToTagsBatchMap builds a BatchAction request body from BatchOpts.
func (opts BatchOpts) ToTagsBatchMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//BatchAction is used to create ,update or delete the tags of a specified image.
func BatchAction(client *golangsdk.ServiceClient, imageID string, opts BatchOptsBuilder) (r ActionResults) {
	b, err := opts.ToTagsBatchMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, imageID), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Get retrieves the tags of a specific image.
func Get(client *golangsdk.ServiceClient, imageID string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, imageID), &r.Body, nil)
	return
}
