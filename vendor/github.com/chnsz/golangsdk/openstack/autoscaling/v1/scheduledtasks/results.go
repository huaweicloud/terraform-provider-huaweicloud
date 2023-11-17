package scheduledtasks

import (
	"github.com/chnsz/golangsdk/pagination"
)

// createResp is the structure that represents the response to creation planned task.
type createResp struct {
	// The ID of planned task.
	TaskID string `json:"task_id"`
}

// ScheduledTask is the structure that represents the planned task detail.
type ScheduledTask struct {
	// The ID of created planned task.
	ID string `json:"task_id"`
	// The ID of scaling group.
	GroupID string `json:"scaling_group_id"`
	// The name of created planned task.
	Name string `json:"name"`
	// The policy of created planned task.
	ScheduledPolicy ScheduledPolicy `json:"scheduled_policy"`
	// The numbers of scaling group instance of created planned task.
	InstanceNumber InstanceNumber `json:"instance_number"`
	// The creation time of planned task.
	CreateTime string `json:"create_time"`
	// The ID of the tenant to which this planned task belongs.
	TenantID string `json:"tenant_id"`
	// The ID of the domain to which this planned task belongs.
	DomainID string `json:"domain_id"`
}

// listResp is the structure that represents the planned tasks list and page detail.
type listResp struct {
	// The list of the planned tasks.
	ScheduledTasks []ScheduledTask `json:"scheduled_tasks"`
	// The page information.
	PageInfo pageInfo `json:"page_info"`
}

// pageInfo is the structure that represents the page information.
type pageInfo struct {
	// The next marker information.
	NextMarker string `json:"next_marker"`
}

// ScheduledTaskPage represents the response pages of the List method.
type ScheduledTaskPage struct {
	pagination.MarkerPageBase
}

// IsEmpty method checks whether the current planned task page is empty.
func (r ScheduledTaskPage) IsEmpty() (bool, error) {
	tasks, err := ExtractScheduledTasks(r)
	return len(tasks) == 0, err
}

// LastMarker method returns the last planned task ID in a planned task page.
func (p ScheduledTaskPage) LastMarker() (string, error) {
	tasks, err := ExtractScheduledTasks(p)
	if err != nil {
		return "", err
	}
	if len(tasks) == 0 {
		return "", nil
	}
	return tasks[len(tasks)-1].ID, nil
}

// ExtractScheduledTasks is a method to extract the list of planned task details.
func ExtractScheduledTasks(r pagination.Page) ([]ScheduledTask, error) {
	var s []ScheduledTask
	err := r.(ScheduledTaskPage).Result.ExtractIntoSlicePtr(&s, "scheduled_tasks")
	return s, err
}
