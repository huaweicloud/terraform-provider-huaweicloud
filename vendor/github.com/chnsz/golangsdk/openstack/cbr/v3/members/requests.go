package members

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure that used to batch add shared members.
type CreateOpts struct {
	// Backup ID.
	BackupId string `json:"-" required:"true"`
	// The list of sharing members configuration.
	Members []string `json:"members" required:"true"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new checkpoint using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) ([]Member, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r struct {
		Members []Member `json:"members"`
	}
	_, err = client.Post(rootURL(client, opts.BackupId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Members, err
}

// ListOpts is the structure that used to query backup shared list.
type ListOpts struct {
	// Backup ID.
	BackupId string `json:"-" required:"true"`
	// The ID of the project that accepts the backup share.
	DestProjectId string `q:"dest_project_id"`
	// The ID of the image registered with the shared backup copy.
	ImageId string `q:"image_id"`
	// Number of records displayed per page.
	// The value must be a positive integer.
	Limit int `q:"limit"`
	// ID of the last record displayed on the previous page.
	Marker string `q:"marker"`
	// Offset value. The value must be a positive integer.
	Offset int `q:"offset"`
	// A group of properties separated by commas (,) and sorting directions.
	// The value is in the format of <key1>[:<direction>],<key2>[:<direction>], where the value of direction is
	// asc (ascending order) or desc (descending order).
	// If a direction is not specified, the default sorting direction is desc.
	// The value of sort can contain a maximum of 255 characters.
	// The key can be as follows: created_at, updated_at, name, status, protected_at, id
	Sort string `q:"sort"`
	// The status of the backup share.
	Status string `q:"status"`
	// Vault ID.
	VaultId string `q:"vault_id"`
}

// List is a method used to query shared member list for specified backup using given parameters.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Member, error) {
	url := rootURL(client, opts.BackupId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := MemberPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractMembers(pages)
}

// Delete is a method to remove a specified member from the specified backup.
func Delete(client *golangsdk.ServiceClient, backupId, memberId string) error {
	_, err := client.Delete(resourceURL(client, backupId, memberId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
