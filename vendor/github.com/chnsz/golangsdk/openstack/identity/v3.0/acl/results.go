package acl

import (
	"github.com/chnsz/golangsdk"
)

// ACLResult is response of the ACL policy for consloe or api access.
type ACLResult struct {
	golangsdk.Result
}

// ConsoleExtract interprets any acl results as a acl policy for console access.
func (r ACLResult) ConsoleExtract() (*ACLPolicy, error) {
	var s struct {
		ACLPolicy *ACLPolicy `json:"console_acl_policy"`
	}
	err := r.ExtractInto(&s)
	return s.ACLPolicy, err
}

// APIExtract interprets any acl results as a acl policy for api access.
func (r ACLResult) APIExtract() (*ACLPolicy, error) {
	var s struct {
		ACLPolicy *ACLPolicy `json:"api_acl_policy"`
	}
	err := r.ExtractInto(&s)
	return s.ACLPolicy, err
}
