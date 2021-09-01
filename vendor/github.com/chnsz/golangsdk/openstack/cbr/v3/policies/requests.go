package policies

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	Name                string          `json:"name" required:"true"`
	OperationDefinition *PolicyODCreate `json:"operation_definition" required:"true"`
	OperationType       string          `json:"operation_type" required:"true"`
	Trigger             *Trigger        `json:"trigger" required:"true"`
	Enabled             *bool           `json:"enabled,omitempty"`
}

// PolicyODCreate is policy operation definition
type PolicyODCreate struct {
	DailyBackups          int    `json:"day_backups,omitempty"`
	WeekBackups           int    `json:"week_backups,omitempty"`
	YearBackups           int    `json:"year_backups,omitempty"`
	MonthBackups          int    `json:"month_backups,omitempty"`
	MaxBackups            int    `json:"max_backups,omitempty"`
	RetentionDurationDays int    `json:"retention_duration_days,omitempty"`
	Timezone              string `json:"timezone,omitempty"`
	EnableAcceleration    bool   `json:"enable_acceleration,omitempty"`
	DestinationProjectID  string `json:"destination_project_id,omitempty"`
	DestinationRegion     string `json:"destination_region,omitempty"`
}

type TriggerProperties struct {
	Pattern []string `json:"pattern"  required:"true"`
}

type Trigger struct {
	Properties TriggerProperties `json:"properties"  required:"true"`
}

type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "policy")
}

//Create is a method by which to create function that create a CBR policy
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, err = client.Post(rootURL(client), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	r.Err = err
	return
}

type ListOpts struct {
	OperationType string `q:"operation_type"`
	VaultID       string `q:"vault_id"`
}

type ListOptsBuilder interface {
	ToPolicyListQuery() (string, error)
}

func (opts ListOpts) ToPolicyListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

//List is a method to obtain the specified CBR policy according to the vault ID or operation type.
//This method can also obtain all the CBR policies through the default parameter settings.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.SinglePageBase(r)}
	})
}

//Get is a method to obtain the specified CBR policy according to the policy ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

type UpdateOpts struct {
	Enabled             *bool           `json:"enabled,omitempty"`
	Name                string          `json:"name,omitempty"`
	OperationDefinition *PolicyODCreate `json:"operation_definition,omitempty"`
	Trigger             *Trigger        `json:"trigger,omitempty"`
}

type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "policy")
}

//Delete is a method to update an existing CBR policy
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToPolicyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, err = client.Put(resourceURL(client, id), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	r.Err = err
	return
}

//Delete is a method to delete an existing CBR policy
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}
