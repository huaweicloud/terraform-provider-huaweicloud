package logtank

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
}

// Opts is a struct that contains all the parameters.
type Opts struct {
	LogGroupID string `json:"log_group_id" required:"true"`

	LogStreamID string `json:"log_stream_id" required:"true"`
}

// Create a logtank with given parameters.
func Create(client *golangsdk.ServiceClient, topicUrn string, opts Opts) (r CreateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(baseURL(client, topicUrn), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})

	return
}

// List all the logtanks under the topicUrn
func List(client *golangsdk.ServiceClient, topicUrn string) (r ListResult) {
	_, r.Err = client.Get(baseURL(client, topicUrn), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Update a logtank with given parameters.
func Update(client *golangsdk.ServiceClient, topicUrn string, logTankID string, opts Opts) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(resourceURL(client, topicUrn, logTankID), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Delete a logtank by topicUrn and LogTankID.
func Delete(client *golangsdk.ServiceClient, topicUrn string, LogTankID string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, topicUrn, LogTankID), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}
