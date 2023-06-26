package attachreplication

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	ToReplicationAttachmentCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies replication attachment creation or import parameters.
type CreateOpts struct {
	// Device is the device that the volume will attach to the instance as.
	// Omit for "auto".
	Device string `json:"device" required:"true"`

	// ReplicationID is the ID of the volume to attach to the instance.
	ReplicationID string `json:"replication_id" required:"true"`
}

// ToReplicationAttachmentCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToReplicationAttachmentCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "replicationAttachment")
}

// Create requests the creation of a new replication attachment on the instance.
func Create(client *golangsdk.ServiceClient, instanceID string, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToReplicationAttachmentCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete requests the deletion of a previous stored ReplicationAttachment from
// the instance.
func Delete(client *golangsdk.ServiceClient, instanceID, replicationID string) (r JobResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = client.DeleteWithResponse(deleteURL(client, instanceID, replicationID), &r.Body, reqOpt)
	return
}
