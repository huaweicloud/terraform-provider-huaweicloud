package activitylogs

import (
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

// ActivityLog is the structure that represents the activity log detail.
type ActivityLog struct {
	ID                  string `json:"id"`
	Status              string `json:"status"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	InstanceRemovedList string `json:"instance_removed_list"`
	InstanceDeletedList string `json:"instance_deleted_list"`
	InstanceAddedList   string `json:"instance_added_list"`
	InstanceValue       int    `json:"instance_value"`
	DesireValue         int    `json:"desire_value"`
	ScalingValue        int    `json:"scaling_value"`
	Description         string `json:"description"`
}

// ActivityLogPage is a single page maximum result representing a query by start_number page.
type ActivityLogPage struct {
	pagination.OffsetPageBase
}

// NextStartNumber returns startNumber of the next element of the page.
func (current ActivityLogPage) NextStartNumber() int {
	q := current.URL.Query()
	// get `start_number` and `limit` from query path.
	// If the limit is not set, it will be set to `20` by default for querying.
	startNumber, _ := strconv.Atoi(q.Get("start_number"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit == 0 {
		limit = 20
	}

	return startNumber + limit
}

// NextPageURL generates the URL for the page of results after this one.
func (current ActivityLogPage) NextPageURL() (string, error) {
	next := current.NextStartNumber()
	if next == 0 {
		return "", nil
	}

	currentURL := current.URL
	q := currentURL.Query()
	q.Set("start_number", strconv.Itoa(next))
	currentURL.RawQuery = q.Encode()

	return currentURL.String(), nil
}

// IsEmpty checks whether a ActivityLogPage struct is empty.
func (b ActivityLogPage) IsEmpty() (bool, error) {
	groups, err := ExtractActivityLogs(b)
	return len(groups) == 0, err
}

// ExtractActivityLogs is a method to extract the list of activity logs for specified scaling group.
func ExtractActivityLogs(r pagination.Page) ([]ActivityLog, error) {
	var s []ActivityLog
	err := r.(ActivityLogPage).Result.ExtractIntoSlicePtr(&s, "scaling_activity_log")
	return s, err
}
