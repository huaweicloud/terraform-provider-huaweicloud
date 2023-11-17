package scheduledtasks

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the structure that used to create planned task.
type CreateOpts struct {
	// The name of planned task to be created.
	Name string `json:"name" required:"true"`
	// The policy of planned task to be created.
	ScheduledPolicy ScheduledPolicy `json:"scheduled_policy" required:"true"`
	// The numbers of scaling group instance for planned task to be created.
	InstanceNumber InstanceNumber `json:"instance_number" required:"true"`
}

// ScheduledPolicy is the structure that used to create the policy of planned task.
type ScheduledPolicy struct {
	// The execution time of the planned task.
	LaunchTime string `json:"launch_time" required:"true"`
	// The type of the planned task.
	RecurrenceType string `json:"recurrence_type,omitempty"`
	// The specific date when the planned task is executed according to the cycle.
	RecurrenceValue string `json:"recurrence_value,omitempty"`
	// The effective start time of the planned task.
	StartTime string `json:"start_time,omitempty"`
	// The effective end time of the planned task.
	EndTime string `json:"end_time,omitempty"`
}

// InstanceNumber is the structure that used to create the numbers of scaling group instance for planned task.
type InstanceNumber struct {
	// The max number of instances of the scaling group to which the planned task belongs.
	Max *int `json:"max,omitempty"`
	// The min number of instances of the scaling group to which the planned task belongs.
	Min *int `json:"min,omitempty"`
	// The desire number of instances of the scaling group to which the planned task belongs.
	Desire *int `json:"desire,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create planned task using given parameters.
func Create(c *golangsdk.ServiceClient, groupID string, opts CreateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var r createResp
	_, err = c.Post(rootURL(c, groupID), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.TaskID, err
}

// ListOpts is the structure that used to query the planned tasks.
type ListOpts struct {
	// The ID of scaling group.
	GroupID string `q:"scaling_group_id"`
	// Number of records displayed per page.
	// The value must be a positive integer.
	Limit int `q:"limit"`
	// The ID of the last record displayed on the previous page.
	Marker string `q:"marker"`
}

// List is a method used to query the planned tasks with given parameters.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]ScheduledTask, error) {
	url := rootURL(c, opts.GroupID)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ScheduledTaskPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractScheduledTasks(pages)
}

// UpdateOpts is the structure that used to update the planned task.
type UpdateOpts struct {
	// The name of planned task to be updated.
	Name string `json:"name,omitempty"`
	// The policy of planned task to be updated.
	ScheduledPolicy *ScheduledPolicy `json:"scheduled_policy,omitempty"`
	// The numbers of scaling group instance for planned task to be updated.
	InstanceNumber *InstanceNumber `json:"instance_number,omitempty"`
}

// Update is a method used to update planned task using given parameters.
func Update(c *golangsdk.ServiceClient, groupID, taskID string, opts UpdateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	var r UpdateOpts
	_, err = c.Put(resourceURL(c, groupID, taskID), b, &r, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return err
}

// Delete is a method used to delete planned task using given parameters.
func Delete(client *golangsdk.ServiceClient, groupID, taskID string) error {
	_, err := client.Delete(resourceURL(client, groupID, taskID), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
