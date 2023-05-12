package backups

import "github.com/chnsz/golangsdk"

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
