package policyexecutelogs

import (
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

// ExecuteLog is the structure that represents the execution log detail.
type ExecuteLog struct {
	// The policy execution log ID.
	ID string `json:"id"`
	// The policy execution status. The value can be SUCCESS, FAIL or EXECUTING.
	Status string `json:"status"`
	// The reason of policy execution failure.
	FailedReason string `json:"failed_reason"`
	// The policy execution type. The value can be SCHEDULED, RECURRENCE, ALARM or MANUAL.
	ExecuteType string `json:"execute_type"`
	// The policy execution time.
	ExecuteTime string `json:"execute_time"`
	// The scaling policy ID.
	PolicyID string `json:"scaling_policy_id"`
	// The scaling resource type. The value can be SCALING_GROUP or BANDWIDTH.
	ResourceType string `json:"scaling_resource_type"`
	// The scaling resource ID.
	ResourceID string `json:"scaling_resource_id"`
	// The scaling original value, indicates the number of instances or bandwidth size.
	OldValue string `json:"old_value"`
	// The scaling target value, indicates the number of instances or bandwidth size.
	DesireValue string `json:"desire_value"`
	// The operational limitations.
	LimitValue string `json:"limit_value"`
	// The policy execution task type. The value can be REMOVE, ADD or SET.
	Type string `json:"type"`
	// The concrete tasks included in executing actions.
	JobRecords []JobRecord `json:"job_records"`
	// The additional information.
	MetaData map[string]interface{} `json:"meta_data"`
	// The project ID.
	TenantID string `json:"tenant_id"`
}

// JobRecord is the structure that represents the concrete tasks included in executing actions.
type JobRecord struct {
	// The job name.
	JobName string `json:"job_name"`
	// The record type. The value can be API or MEG.
	RecordType string `json:"record_type"`
	// The record time.
	RecordTime string `json:"record_time"`
	// The request information.
	Request string `json:"request"`
	// The response information.
	Response string `json:"response"`
	// The response code.
	Code string `json:"code"`
	// The message content.
	Message string `json:"message"`
	// The job execution status. The value can be SUCCESS or FAIL.
	JobStatus string `json:"job_status"`
}

// ExecuteLogPage is a single page maximum result representing a query by start_number page.
type ExecuteLogPage struct {
	pagination.OffsetPageBase
}

// NextStartNumber returns startNumber of the next element of the page.
func (current ExecuteLogPage) NextStartNumber() int {
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
func (current ExecuteLogPage) NextPageURL() (string, error) {
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

// IsEmpty checks whether a ExecuteLogPage struct is empty.
func (b ExecuteLogPage) IsEmpty() (bool, error) {
	logs, err := extractExecuteLogs(b)
	return len(logs) == 0, err
}

// extractExecuteLogs is a method to extract the list of execution logs for specified scaling policy.
func extractExecuteLogs(r pagination.Page) ([]ExecuteLog, error) {
	var s []ExecuteLog
	err := r.(ExecuteLogPage).Result.ExtractIntoSlicePtr(&s, "scaling_policy_execute_log")
	return s, err
}
