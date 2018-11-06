package backups

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the backup attributes you want to see returned.
type ListOpts struct {
	Id         string
	SnapshotId string
	Name       string `q:"name"`
	Status     string `q:"status"`
	Limit      int    `q:"limit"`
	Offset     int    `q:"offset"`
	VolumeId   string `q:"volume_id"`
}

// List returns collection of
// Backup. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those Backup that are owned by the
// tenant who submits the request, unless an admin user submits the request.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Backup, error) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return nil, err
	}
	u := listURL(c) + q.String()
	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	if err != nil {
		return nil, err
	}

	allBackups, err := ExtractBackups(pages)
	if err != nil {
		return nil, err
	}

	return FilterBackups(allBackups, opts)
}

func FilterBackups(backup []Backup, opts ListOpts) ([]Backup, error) {

	var refinedBackup []Backup
	var matched bool
	m := map[string]interface{}{}

	if opts.Id != "" {
		m["Id"] = opts.Id
	}
	if opts.SnapshotId != "" {
		m["SnapshotId"] = opts.SnapshotId
	}

	if len(m) > 0 && len(backup) > 0 {
		for _, backup := range backup {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&backup, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedBackup = append(refinedBackup, backup)
			}
		}
	} else {
		refinedBackup = backup
	}
	return refinedBackup, nil
}

func getStructField(v *Backup, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new backup.
type CreateOpts struct {
	//ID of the disk to be backed up
	VolumeId string `json:"volume_id" required:"true"`
	//Snapshot ID of the disk to be backed up
	SnapshotId string `json:"snapshot_id,omitempty" `
	//Backup name, which cannot start with autobk
	Name string `json:"name" required:"true"`
	//Backup description
	Description string `json:"description,omitempty"`
	//List of tags to be configured for the backup resources
	Tags []Tag `json:"tags,omitempty"`
}

type Tag struct {
	//Tag key
	Key string `json:"key" required:"true"`
	//Tag value
	Value string `json:"value" required:"true"`
}

// ToBackupCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToBackupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "backup")
}

// Create will create a new Backup based on the values in CreateOpts. To extract
// the Backup object from the response, call the ExtractJobResponse method on the
// JobResult.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r JobResult) {
	b, err := opts.ToBackupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// RestoreOptsBuilder allows extensions to add additional parameters to the
// Create request.
type RestoreOptsBuilder interface {
	ToRestoreCreateMap() (map[string]interface{}, error)
}

// BackupRestoreOpts contains all the values needed to create a new backup.
type BackupRestoreOpts struct {
	//ID of the disk to be backed up
	VolumeId string `json:"volume_id" required:"true"`
}

// ToRestoreCreateMap builds a create request body from BackupRestoreOpts.
func (opts BackupRestoreOpts) ToRestoreCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "restore")
}

// CreateBackupRestore will create a new Restore based on the values in BackupRestoreOpts. To extract
// the BackupRestoreInfo object from the response, call the ExtractBackupRestore method on the
// CreateResult.
func CreateBackupRestore(c *golangsdk.ServiceClient, id string, opts RestoreOptsBuilder) (r CreateResult) {
	b, err := opts.ToRestoreCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(restoreURL(c, id), b, &r.Body, nil)
	return
}

// Get retrieves a particular backup based on its unique ID. To extract
//// the Backup object from the response, call the Extract method on the
//// GetResult.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// Delete will permanently delete a particular backup based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
