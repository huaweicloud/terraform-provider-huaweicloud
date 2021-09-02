package securities

import "github.com/chnsz/golangsdk"

type commonResult struct {
	golangsdk.Result
}

// SSLUpdateResult represents a result of the ConfigureSSL method.
type SSLUpdateResult struct {
	golangsdk.ErrResult
}

// DBPortUpdateResult represents a result of the UpdateDBPort method.
type DBPortUpdateResult struct {
	commonResult
}

// SecGroupUpdateResult represents a result of the UpdateSecGroup method.
type SecGroupUpdateResult struct {
	commonResult
}

// WorkFlow is a struct that represents the result of database updation.
type WorkFlow struct {
	// Indicates the workflow ID.
	WorkflowId string `json:"workflowId"`
}

func (r commonResult) Extract() (*WorkFlow, error) {
	var s WorkFlow
	err := r.ExtractInto(&s)
	return &s, err
}
