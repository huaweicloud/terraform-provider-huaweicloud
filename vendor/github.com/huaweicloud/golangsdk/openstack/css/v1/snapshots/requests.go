package snapshots

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSnapshotCreateMap() (map[string]interface{}, error)
}

// PolicyCreateOpts contains options for creating a snapshot policy.
// This object is passed to the snapshots.PolicyCreate function.
type PolicyCreateOpts struct {
	Prefix     string `json:"prefix" required:"true"`
	Period     string `json:"period" required:"true"`
	KeepDay    int    `json:"keepday" required:"true"`
	Enable     string `json:"enable" required:"true"`
	DeleteAuto string `json:"deleteAuto,omitempty"`
}

// ToSnapshotCreateMap assembles a request body based on the contents of a
// PolicyCreateOpts.
func (opts PolicyCreateOpts) ToSnapshotCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// PolicyCreate will create a new snapshot policy based on the values in PolicyCreateOpts.
func PolicyCreate(client *golangsdk.ServiceClient, opts CreateOptsBuilder, clusterId string) (r ErrorResult) {
	b, err := opts.ToSnapshotCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(policyURL(client, clusterId), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// PolicyGet retrieves the snapshot policy with the provided cluster ID.
// To extract the snapshot policy object from the response, call the Extract method on the GetResult.
func PolicyGet(client *golangsdk.ServiceClient, clusterId string) (r PolicyResult) {
	_, r.Err = client.Get(policyURL(client, clusterId), &r.Body, nil)
	return
}

// Enable will enable the Snapshot function with the provided ID.
func Enable(client *golangsdk.ServiceClient, clusterId string) (r ErrorResult) {
	_, r.Err = client.Post(enableURL(client, clusterId), nil, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

// Disable will disable the Snapshot function with the provided ID.
func Disable(client *golangsdk.ServiceClient, clusterId string) (r ErrorResult) {
	_, r.Err = client.Delete(disableURL(client, clusterId), nil)
	return
}

// CreateOpts contains options for creating a snapshot.
// This object is passed to the snapshots.Create function.
type CreateOpts struct {
	Name        string `json:"name" required:"true"`
	Description string `json:"description,omitempty"`
	Indices     string `json:"indices,omitempty"`
}

// ToSnapshotCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToSnapshotCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new snapshot based on the values in CreateOpts.
// To extract the result from the response, call the Extract method on the CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder, clusterId string) (r CreateResult) {
	b, err := opts.ToSnapshotCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client, clusterId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// List retrieves the Snapshots with the provided ID. To extract the Snapshot
// objects from the response, call the Extract method on the GetResult.
func List(client *golangsdk.ServiceClient, clusterId string) (r ListResult) {
	_, r.Err = client.Get(listURL(client, clusterId), &r.Body, nil)
	return
}

// Delete will delete the existing Snapshot ID with the provided ID.
func Delete(client *golangsdk.ServiceClient, clusterId, id string) (r ErrorResult) {
	_, r.Err = client.Delete(deleteURL(client, clusterId, id), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	})
	return
}
