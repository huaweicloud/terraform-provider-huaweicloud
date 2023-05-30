package backups

import (
	"github.com/chnsz/golangsdk"
)

// UpdateResult represents the result of a update operation.
type UpdateResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	golangsdk.Result
}

type BackupPolicy struct {
	KeepDays  int    `json:"keep_days"`
	StartTime string `json:"start_time"`
	Period    string `json:"period"`
}

func (r GetResult) Extract() (*BackupPolicy, error) {
	var s struct {
		BackupPolicy *BackupPolicy `json:"backup_policy"`
	}
	err := r.ExtractInto(&s)
	return s.BackupPolicy, err
}
