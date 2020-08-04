package volumes

import (
	"github.com/huaweicloud/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVolumeCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Volume. This object is passed to
// the volumes.Create function. For more information about these parameters,
// see the Volume object.
type CreateOpts struct {
	// The availability zone
	AvailabilityZone string `json:"availability_zone" required:"true"`
	// The associated volume type
	VolumeType string `json:"volume_type" required:"true"`
	// The volume name
	Name string `json:"name,omitempty"`
	// The volume description
	Description string `json:"description,omitempty"`
	// The size of the volume, in GB
	Size int `json:"size,omitempty"`
	// The number to be created in a batch
	Count int `json:"count,omitempty"`
	// The backup_id
	BackupID string `json:"backup_id,omitempty"`
	// the ID of the existing volume snapshot
	SnapshotID string `json:"snapshot_id,omitempty"`
	// the ID of the image in IMS
	ImageRef string `json:"imageRef,omitempty"`
	// This field is no longer used. Use multiattach.
	Shareable string `json:"shareable,omitempty"`
	// Shared disk
	Multiattach bool `json:"multiattach,omitempty"`
	// One or more metadata key and value pairs to associate with the volume
	Metadata map[string]string `json:"metadata,omitempty"`
	// One or more tag key and value pairs to associate with the volume
	Tags map[string]string `json:"tags,omitempty"`
	// the enterprise project id
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

// ToVolumeCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToVolumeCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "volume")
}

// Create will create a new Volume based on the values in CreateOpts.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToVolumeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get retrieves the Volume with the provided ID. To extract the Volume object
// from the response, call the Extract method on the GetResult.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
