package checkpoints

import "github.com/chnsz/golangsdk"

// CreateOpts is the structure that used to create a checkpoint and backup some resources.
type CreateOpts struct {
	// The vault ID.
	VaultId string `json:"vault_id" required:"true"`
	// Parameters.
	Parameters CheckpointParameter `json:"parameters,omitempty"`
}

// CheckpointParameter is the structure that represents the backup configuration.
type CheckpointParameter struct {
	// Whether automatic triggering is enabled.
	// Possible values are true (yes) and false (no).
	// Defaults to false.
	AutoTrigger *bool `json:"auto_trigger,omitempty"`
	// Backup description.
	Description string `json:"description,omitempty"`
	// Whether the backup is an incremental backup.
	// Possible values are true (yes) and false (no).
	// Defaults to true.
	Incremental *bool `json:"incremental,omitempty"`
	// Backup name, which can contain only digits, letters, underscores (_), and hyphens (-).
	Name string `json:"name,omitempty"`
	// The UUID list of resources to be backed up.
	Resources []string `json:"resources,omitempty"`
	// Resource details.
	ResourceDetails []Resource `json:"resource_details,omitempty"`
}

// Resource is the structure that represents the backup resource detail.
type Resource struct {
	// ID of the vault resource type.
	ID string `json:"id" required:"true"`
	// The type of the resource to be backed up, which can be:
	// + OS::Nova::Server
	// + OS::Cinder::Volume
	// + OS::Ironic::BareMetalServer
	// + OS::Native::Server
	// + OS::Sfs::Turbo
	// + OS::Workspace::DesktopV2
	Type string `json:"type" required:"true"`
	// The vaule
	Name string `json:"name,omitempty"`
	// The extra volume configuration.
	ExtraInfo ResourceExtraInfo `json:"extra_info,omitempty"`
}

// ResourceExtraInfo is the structure that represents the configuration of the backup volume.
type ResourceExtraInfo struct {
	// The IDs of the disks that will not be backed up. This parameter is used when servers are added to a vault,
	// which include all server disks. But some disks do not need to be backed up. Or in case that a server was
	// previously added and some disks on this server do not need to be backed up.
	ExcludeVolumes []string `json:"exclude_volumes,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new checkpoint using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*Checkpoint, error) {
	b, err := golangsdk.BuildRequestBody(opts, "checkpoint")
	if err != nil {
		return nil, err
	}

	var r struct {
		Checkpoint Checkpoint `json:"checkpoint"`
	}
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Checkpoint, err
}

// Get is a method to obtain an existing checkpoint by its ID.
func Get(client *golangsdk.ServiceClient, checkpointId string) (*Checkpoint, error) {
	var r struct {
		Checkpoint Checkpoint `json:"checkpoint"`
	}
	_, err := client.Get(resourceURL(client, checkpointId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Checkpoint, err
}
