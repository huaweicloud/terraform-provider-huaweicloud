package policies

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains the options for create a Policy. This object is
// passed to Create(). For more information about these parameters,
// please refer to the Policy object, or the volume backup service API v2
// documentation
type CreateOpts struct {
	//Backup policy name.It cannot start with default.
	Name string `json:"backup_policy_name" required:"true"`
	//Details about the scheduling policy
	ScheduledPolicy ScheduledPolicy `json:"scheduled_policy" required:"true"`
	// Tags to be configured for the backup policy
	Tags []Tag `json:"tags,omitempty"`
}

// ScheduledPolicy defines the details about scheduling policy for create
type ScheduledPolicy struct {
	//Start time of the backup job.
	StartTime string `json:"start_time" required:"true"`
	//Backup interval (1 to 14 days)
	Frequency int `json:"frequency,omitempty"`
	//Specifies on which days of each week backup jobs are ececuted.
	WeekFrequency []string `json:"week_frequency,omitempty"`
	//Number of retained backups, minimum 2.
	RententionNum int `json:"rentention_num,omitempty"`
	//Days of retained backups, minimum 2.
	RententionDay int `json:"rentention_day,omitempty"`
	//Whether to retain the first backup in the current month, possible values Y or N
	RemainFirstBackup string `json:"remain_first_backup_of_curMonth" required:"true"`
	//Backup policy status, ON or OFF
	Status string `json:"status" required:"true"`
}

type Tag struct {
	//Tag key. A tag key consists of up to 36 characters
	Key string `json:"key" required:"true"`
	//Tag value. A tag value consists of 0 to 43 characters
	Value string `json:"value" required:"true"`
}

// ToPolicyCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new Policy based on the values in CreateOpts. To extract
// the Policy object from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(commonURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the options for Update a Policy.
type UpdateOpts struct {
	//Backup policy name.It cannot start with default.
	Name string `json:"backup_policy_name,omitempty"`
	//Details about the scheduling policy
	ScheduledPolicy UpdateSchedule `json:"scheduled_policy,omitempty"`
}

// UpdateSchedule defiens the details about scheduling policy for update.
type UpdateSchedule struct {
	//Start time of the backup job.
	StartTime string `json:"start_time,omitempty"`
	//Backup interval (1 to 14 days)
	Frequency int `json:"frequency,omitempty"`
	//Specifies on which days of each week backup jobs are ececuted.
	WeekFrequency []string `json:"week_frequency,omitempty"`
	//Number of retained backups, minimum 2.
	RententionNum int `json:"rentention_num,omitempty"`
	//Days of retained backups, minimum 2.
	RententionDay int `json:"rentention_day,omitempty"`
	//Number of retained backups, minimum 2.
	RemainFirstBackup string `json:"remain_first_backup_of_curMonth,omitempty"`
	//Backup policy status, ON or OFF
	Status string `json:"status,omitempty"`
}

// ToPolicyUpdateMap assembles a request body based on the contents of a
// UpdateOpts.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

//Update will Update an existing backup Policy based on the values in UpdateOpts.To extract
// the Policy object from the response, call the Extract method on the
// UpdateResult.
func Update(c *golangsdk.ServiceClient, policyID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID), b, &r.Body, reqOpt)
	return
}

//Delete will delete the specified backup policy.
func Delete(c *golangsdk.ServiceClient, policyID string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, policyID), nil)
	return
}

//ListOpts allows filtering policies
type ListOpts struct {
	//Backup policy ID
	ID string
	//Backup policy name
	Name string
	//Backup policy status
	Status string
}

// List returns a Pager which allows you to iterate over a collection of
// Policies. It accepts a ListOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]Policy, error) {

	pages, err := pagination.NewPager(c, commonURL(c), func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allPolicies, err := ExtractPolicies(pages)
	if err != nil {
		return nil, err
	}

	return FilterPolicies(allPolicies, opts)
}

func FilterPolicies(policies []Policy, opts ListOpts) ([]Policy, error) {

	var refinedPolicies []Policy
	var matched bool

	m := map[string]FilterStruct{}

	if opts.ID != "" {
		m["ID"] = FilterStruct{Value: opts.ID}
	}
	if opts.Name != "" {
		m["Name"] = FilterStruct{Value: opts.Name}
	}

	if opts.Status != "" {
		m["Status"] = FilterStruct{Value: opts.Status, Driller: []string{"ScheduledPolicy"}}
	}

	if len(m) > 0 && len(policies) > 0 {
		for _, policies := range policies {
			matched = true

			for key, value := range m {
				if sVal := GetStructNestedField(&policies, key, value.Driller); !(sVal == value.Value) {
					matched = false
				}
			}
			if matched {
				refinedPolicies = append(refinedPolicies, policies)
			}
		}
	} else {
		refinedPolicies = policies
	}
	return refinedPolicies, nil
}

func GetStructNestedField(v *Policy, field string, structDriller []string) string {
	r := reflect.ValueOf(v)
	for _, drillField := range structDriller {
		f := reflect.Indirect(r).FieldByName(drillField).Interface()
		r = reflect.ValueOf(f)
	}
	f1 := reflect.Indirect(r).FieldByName(field)
	return string(f1.String())
}

type FilterStruct struct {
	Value   string
	Driller []string
}

// AssociateOptsBuilder allows extensions to add additional parameters to the
// Associate request.
type AssociateOptsBuilder interface {
	ToPolicyAssociateMap() (map[string]interface{}, error)
}

// AssociateOpts contains the options to associate a resource to a Policy.
type AssociateOpts struct {
	//Backup policy ID, to which the resource is to be associated.
	PolicyID string `json:"backup_policy_id" required:"true"`
	//Details about the resources to associate with the policy.
	Resources []AssociateResource `json:"resources" required:"true"`
}

type AssociateResource struct {
	//The ID of the resource to associate with policy.
	ResourceID string `json:"resource_id" required:"true"`
	//Type of the resource , e.g. volume.
	ResourceType string `json:"resource_type" required:"true"`
}

// ToPolicyAssociateMap assembles a request body based on the contents of a
// AssociateOpts.
func (opts AssociateOpts) ToPolicyAssociateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Associate will associate a resource tp a backup policy based on the values in AssociateOpts. To extract
// the associated resources from the response, call the ExtractResource method on the
// ResourceResult.
func Associate(client *golangsdk.ServiceClient, opts AssociateOpts) (r ResourceResult) {
	b, err := opts.ToPolicyAssociateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(associateURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// DisassociateOptsBuilder allows extensions to add additional parameters to the
// Disassociate request.
type DisassociateOptsBuilder interface {
	ToPolicyDisassociateMap() (map[string]interface{}, error)
}

// DisassociateOpts contains the options disassociate a resource from a Policy.
type DisassociateOpts struct {
	//Disassociate Resources
	Resources []DisassociateResource `json:"resources" required:"true"`
}

type DisassociateResource struct {
	//ResourceID
	ResourceID string `json:"resource_id" required:"true"`
}

// ToPolicyDisassociateMap assembles a request body based on the contents of a
// DisassociateOpts.
func (opts DisassociateOpts) ToPolicyDisassociateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Disassociate will disassociate a resource from a backup policy based on the values in DisassociateOpts. To extract
// the disassociated resources from the response, call the ExtractResource method on the
// ResourceResult.
func Disassociate(client *golangsdk.ServiceClient, policyID string, opts DisassociateOpts) (r ResourceResult) {
	b, err := opts.ToPolicyDisassociateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(disassociateURL(client, policyID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
