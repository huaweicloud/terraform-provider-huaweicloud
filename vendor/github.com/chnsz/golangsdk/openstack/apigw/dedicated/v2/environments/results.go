package environments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// UpdateResult represents a result of the Update method.
type UpdateResult struct {
	commonResult
}

type Environment struct {
	// Environment ID.
	Id string `json:"id"`
	// Environment name.
	Name string `json:"name"`
	// Description.
	Description string `json:"remark"`
	// Create time, in RFC-3339 format.
	CreateTime string `json:"create_time"`
}

func (r commonResult) Extract() (*Environment, error) {
	var s Environment
	err := r.ExtractInto(&s)
	return &s, err
}

// EnvironmentPage represents the response pages of the List method.
type EnvironmentPage struct {
	pagination.SinglePageBase
}

func ExtractEnvironments(r pagination.Page) ([]Environment, error) {
	var s []Environment
	err := r.(EnvironmentPage).Result.ExtractIntoSlicePtr(&s, "envs")
	return s, err
}

// DeleteResult represents a result of the Delete and DeleteVariable method.
type DeleteResult struct {
	golangsdk.ErrResult
}

type VariableResult struct {
	golangsdk.Result
}

// VariableCreateResult represents a result of the CreateVariable method.
type VariableCreateResult struct {
	VariableResult
}

// VariableGetResult represents a result of the GetVariable operation.
type VariableGetResult struct {
	VariableResult
}

type Variable struct {
	// Environment variable ID.
	Id string `json:"id"`
	// Variable name.
	Name string `json:"variable_name"`
	// Variable value.
	Value string `json:"variable_value"`
	// API group ID.
	GroupId string `json:"group_id"`
	// Environment ID.
	EnvId string `json:"env_id"`
}

func (r VariableResult) Extract() (*Variable, error) {
	var s Variable
	err := r.ExtractInto(&s)
	return &s, err
}

// VariablePage represents the response pages of the List operation.
type VariablePage struct {
	pagination.SinglePageBase
}

func ExtractVariables(r pagination.Page) ([]Variable, error) {
	var s []Variable
	err := r.(VariablePage).Result.ExtractIntoSlicePtr(&s, "variables")
	return s, err
}
