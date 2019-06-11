package policies

import (
	"encoding/json"
	"strconv"

	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateBackupPolicy struct {
	CreatedAt           time.Time                      `json:"-"`
	Description         string                         `json:"description"`
	ID                  string                         `json:"id"`
	Name                string                         `json:"name"`
	Parameters          PolicyParam                    `json:"parameters"`
	ProjectId           string                         `json:"project_id"`
	ProviderId          string                         `json:"provider_id"`
	Resources           []Resource                     `json:"resources"`
	ScheduledOperations []CreateScheduledOperationResp `json:"scheduled_operations"`
	Status              string                         `json:"status"`
	Tags                []ResourceTag                  `json:"tags"`
}

type BackupPolicy struct {
	CreatedAt           time.Time                `json:"-"`
	Description         string                   `json:"description"`
	ID                  string                   `json:"id"`
	Name                string                   `json:"name"`
	Parameters          PolicyParam              `json:"parameters"`
	ProjectId           string                   `json:"project_id"`
	ProviderId          string                   `json:"provider_id"`
	Resources           []Resource               `json:"resources"`
	ScheduledOperations []ScheduledOperationResp `json:"scheduled_operations"`
	Status              string                   `json:"status"`
	Tags                []ResourceTag            `json:"tags"`
}

type ScheduledOperationResp struct {
	Description         string                  `json:"description"`
	Enabled             bool                    `json:"enabled"`
	Name                string                  `json:"name"`
	OperationType       string                  `json:"operation_type"`
	OperationDefinition OperationDefinitionResp `json:"operation_definition"`
	Trigger             TriggerResp             `json:"trigger"`
	ID                  string                  `json:"id"`
	TriggerID           string                  `json:"trigger_id"`
}

type CreateScheduledOperationResp struct {
	Description         string                        `json:"description"`
	Enabled             bool                          `json:"enabled"`
	Name                string                        `json:"name"`
	OperationType       string                        `json:"operation_type"`
	OperationDefinition CreateOperationDefinitionResp `json:"operation_definition"`
	Trigger             TriggerResp                   `json:"trigger"`
	ID                  string                        `json:"id"`
	TriggerID           string                        `json:"trigger_id"`
}

type OperationDefinitionResp struct {
	MaxBackups            int    `json:"max_backups"`
	RetentionDurationDays int    `json:"retention_duration_days"`
	Permanent             bool   `json:"permanent"`
	PlanId                string `json:"plan_id"`
	ProviderId            string `json:"provider_id"`
}

type CreateOperationDefinitionResp struct {
	MaxBackups            int    `json:"-"`
	RetentionDurationDays int    `json:"-"`
	Permanent             bool   `json:"-"`
	PlanId                string `json:"plan_id"`
	ProviderId            string `json:"provider_id"`
}

type TriggerResp struct {
	Properties TriggerPropertiesResp `json:"properties"`
	Name       string                `json:"name"`
	ID         string                `json:"id"`
	Type       string                `json:"type"`
}

type TriggerPropertiesResp struct {
	Pattern   string    `json:"pattern"`
	StartTime time.Time `json:"-"`
}

// UnmarshalJSON helps to unmarshal BackupPolicy fields into needed values.
func (r *BackupPolicy) UnmarshalJSON(b []byte) error {
	type tmp BackupPolicy
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = BackupPolicy(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)

	return err
}

// UnmarshalJSON helps to unmarshal TriggerPropertiesResp fields into needed values.
func (r *TriggerPropertiesResp) UnmarshalJSON(b []byte) error {
	type tmp TriggerPropertiesResp
	var s struct {
		tmp
		StartTime golangsdk.JSONRFC3339ZNoTNoZ `json:"start_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = TriggerPropertiesResp(s.tmp)

	r.StartTime = time.Time(s.StartTime)

	return err
}

// UnmarshalJSON helps to unmarshal OperationDefinitionResp fields into needed values.
func (r *CreateOperationDefinitionResp) UnmarshalJSON(b []byte) error {
	type tmp CreateOperationDefinitionResp
	var s struct {
		tmp
		MaxBackups            string `json:"max_backups"`
		RetentionDurationDays string `json:"retention_duration_days"`
		Permanent             string `json:"permanent"`
	}

	err := json.Unmarshal(b, &s)

	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError: //check if type error occurred (handles if no type conversion is required for cloud like Huawei)

			var s struct {
				tmp
				MaxBackups            int  `json:"max_backups"`
				RetentionDurationDays int  `json:"retention_duration_days"`
				Permanent             bool `json:"permanent"`
			}
			err := json.Unmarshal(b, &s)
			if err != nil {
				return err
			}
			*r = CreateOperationDefinitionResp(s.tmp)
			r.MaxBackups = s.MaxBackups
			r.RetentionDurationDays = s.RetentionDurationDays
			r.Permanent = s.Permanent
			return nil
		default:
			return err
		}
	}

	*r = CreateOperationDefinitionResp(s.tmp)

	switch s.MaxBackups {
	case "":
		r.MaxBackups = 0
	default:
		r.MaxBackups, err = strconv.Atoi(s.MaxBackups)
		if err != nil {
			return err
		}
	}

	switch s.RetentionDurationDays {
	case "":
		r.RetentionDurationDays = 0
	default:
		r.RetentionDurationDays, err = strconv.Atoi(s.RetentionDurationDays)
		if err != nil {
			return err
		}
	}

	switch s.Permanent {
	case "":
		r.Permanent = false
	default:
		r.Permanent, err = strconv.ParseBool(s.Permanent)
		if err != nil {
			return err
		}
	}

	return err
}

// Extract will get the backup policies object from the commonResult
func (r commonResult) Extract() (*BackupPolicy, error) {
	var s struct {
		BackupPolicy *BackupPolicy `json:"policy"`
	}

	err := r.ExtractInto(&s)
	return s.BackupPolicy, err
}

func (r cuResult) Extract() (*CreateBackupPolicy, error) {
	var s struct {
		BackupPolicy *CreateBackupPolicy `json:"policy"`
	}

	err := r.ExtractInto(&s)
	return s.BackupPolicy, err
}

// BackupPolicyPage is the page returned by a pager when traversing over a
// collection of backup policies.
type BackupPolicyPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of backup policies has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r BackupPolicyPage) NextPageURL() (string, error) {
	var s struct {
		Links []golangsdk.Link `json:"policies_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return golangsdk.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a BackupPolicyPage struct is empty.
func (r BackupPolicyPage) IsEmpty() (bool, error) {
	is, err := ExtractBackupPolicies(r)
	return len(is) == 0, err
}

// ExtractBackupPolicies accepts a Page struct, specifically a BackupPolicyPage struct,
// and extracts the elements into a slice of Policy structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBackupPolicies(r pagination.Page) ([]BackupPolicy, error) {
	var s struct {
		BackupPolicies []BackupPolicy `json:"policies"`
	}
	err := (r.(BackupPolicyPage)).ExtractInto(&s)
	return s.BackupPolicies, err
}

type commonResult struct {
	golangsdk.Result
}

type cuResult struct {
	golangsdk.Result
}

type CreateResult struct {
	cuResult
}

type GetResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type UpdateResult struct {
	cuResult
}

type ListResult struct {
	commonResult
}
