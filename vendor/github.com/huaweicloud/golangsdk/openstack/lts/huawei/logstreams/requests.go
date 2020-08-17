package logstreams

import "github.com/huaweicloud/golangsdk"

// CreateOptsBuilder is used for creating log stream parameters.
type CreateOptsBuilder interface {
	ToLogStreamsCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct that contains all the parameters.
type CreateOpts struct {
	// Specifies the log stream name.
	LogStreamName string `json:"log_stream_name" required:"true"`
}

// ToLogStreamsCreateMap is used for type convert
func (ops CreateOpts) ToLogStreamsCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a log stream with given parameters.
func Create(client *golangsdk.ServiceClient, groupId string, ops CreateOptsBuilder) (r CreateResult) {
	b, err := ops.ToLogStreamsCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client, groupId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})

	return
}

// Delete a log stream by id
func Delete(client *golangsdk.ServiceClient, groupId string, id string) (r DeleteResult) {
	opts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=utf8"},
	}
	_, r.Err = client.Delete(deleteURL(client, groupId, id), &opts)
	return
}

// Get log stream list
func List(client *golangsdk.ServiceClient, groupId string) (r ListResults) {
	opts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=utf8"},
	}
	_, r.Err = client.Get(listURL(client, groupId), &r.Body, &opts)
	return
}
