package backups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Get is a method to obtain an specified backup by its ID.
func Get(client *golangsdk.ServiceClient, backupId string) (*BackupResp, error) {
	var r getResp
	_, err := client.Get(resourceURL(client, backupId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r.Backup, err
}

// ListOpts is the structure that used to query backup list using given parameters.
type ListOpts struct {
	// Restore point ID.
	CheckpointId string `q:"checkpoint_id"`
	// Dedicated cloud tag, which only takes effect in dedicated cloud scenarios.
	Dec bool `q:"dec"`
	// Time when the backup ends, in %YYYY-%mm-%ddT%HH:%MM:%SSZ format. For example, 2018-02-01T12:00:00Z.
	EndTime string `q:"end_time"`
	// Enterprise project ID or all_granted_eps. all_granted_eps indicates querying the IDs of all enterprise projects
	// on which the user has permissions.
	EnterpriseProjectId string `q:"enterprise_project_id"`
	// Backup type, which can be backup or replication.
	ImageType string `q:"image_type"`
	// Whether incremental backup is used.
	// Default: false
	Incremental bool `q:"incremental"`
	// Number of records displayed per page.
	// The value must be a positive integer.
	Limit int `q:"limit"`
	// ID of the last record displayed on the previous page.
	Marker string `q:"marker"`
	// Backup sharing status
	// Enumeration values:
	// + pending
	// + accepted
	// + rejected
	MemberStatus string `q:"member_status"`
	// Backup name.
	Name string `q:"name"`
	// Offset value. The value must be a positive integer.
	Offset int `q:"offset"`
	// Owning type of a backup. private backups are queried by default.
	// Enumeration values:
	// + all_granted
	// + private
	// + shared
	// Default: private
	OwnType string `q:"own_type"`
	// Parent backup ID.
	ParentId string `q:"parent_id"`
	// AZ-based filtering is supported.
	ResourceAZ string `q:"resource_az"`
	// Resource ID.
	ResourceId string `q:"resource_id"`
	// Resource name.
	ResourceName string `q:"resource_name"`
	// Resource type, which can be:
	// + OS::Nova::Server
	// + OS::Cinder::Volume
	// + OS::Ironic::BareMetalServer
	// + OS::Native::Server
	// + OS::Sfs::Turbo
	// + OS::Workspace::DesktopV2
	ResourceType string `q:"resource_type"`
	// Whether to show replication records.
	// Default: false
	ShowReplication bool `q:"show_replication"`
	// A group of properties separated by commas (,) and sorting directions.
	// The value is in the format of <key1>[:<direction>],<key2>[:<direction>], where the value of direction is
	// asc (ascending order) or desc (descending order).
	// If a direction is not specified, the default sorting direction is desc.
	// The value of sort can contain a maximum of 255 characters.
	// The key can be as follows: created_at, updated_at, name, status, protected_at, id
	Sort string `q:"sort"`
	// Time when the backup starts, in %YYYY-%mm-%ddT%HH:%MM:%SSZ format.
	// For example, 2018-02-01T12:00:00Z.
	StartTime string `q:"start_time"`
	// Status When the API is called, multiple statuses can be transferred for filtering.
	// for example, status=available&status=error.
	// Enumeration values:
	// + available
	// + protecting
	// + deleting
	// + restoring
	// + error
	// + waiting_protect
	// + waiting_delete
	// + waiting_restore
	Status string `q:"status"`
	// Backups are filtered based on the occupied vault capacity.
	// The value ranges from 1 to 100.
	// For example, if used_percent is set to 80, all backups who occupied 80% or more of the vault capacity are
	// displayed.
	UsedPercent string `q:"used_percent"`
	// Vault ID.
	VaultId string `q:"vault_id"`
}

// List is a method used to query backup list under specified checkpoint using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]BackupResp, error) {
	url := rootURL(client)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := BackupPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractBackups(pages)
}

// Delete is a method to delete an specified backup by its ID.
func Delete(client *golangsdk.ServiceClient, backupId string) error {
	_, err := client.Delete(resourceURL(client, backupId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
