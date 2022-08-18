package roles

import (
	"encoding/json"

	"github.com/chnsz/golangsdk/pagination"
)

// ListResp is the structure that represents the API response of 'List' method.
type ListResp struct {
	// Total number of query results.
	TotalCount int `json:"total_count"`
	// List of roles to query.
	Roles string `json:"roles"`
}

// RoleResp is the structure that represents the detail of the database role.
type RoleResp struct {
	// Whether role is built-in.
	IsBuiltin bool `json:"isBuiltin"`
	// Role name.
	Name string `json:"role"`
	// Database name.
	DbName string `json:"db"`
	// The list of privileges inherited by the newly created role.
	Privileges []Privilege `json:"privileges"`
	// The list of privileges inherited by the newly created role, includes all privileges inherited by inherited roles.
	InheritedPrivileges []Privilege `json:"inheritedPrivileges"`
	// The list of roles inherited by the newly created role.
	Roles []RoleDetail `json:"roles"`
	// The list of roles inherited by the newly created role, includes all roles inherited by inherited roles.
	InheritedRoles []RoleDetail `json:"inheritedRoles"`
}

// Privilege is the structure that represents the privilege detail for database.
type Privilege struct {
	// The details of the resource to which the privilege belongs.
	Resource Resource `json:"resource"`
	// The operation permission list.
	Actions []string `json:"actions"`
}

// Resource is the structure that represents the database details to which the role and user belongs.
type Resource struct {
	// The database to which the privilege belongs.
	Collection string `json:"collection"`
	// The database name.
	DbName string `json:"db"`
}

// RoleDetail is the structure that represents the inherited role details.
type RoleDetail struct {
	// Role name.
	Name string `json:"role"`
	// The database name to which the role belongs.
	DbName string `json:"db"`
}

// RolePage is a single page maximum result representing a query by offset page.
type RolePage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a RolePage struct is empty.
func (p RolePage) IsEmpty() (bool, error) {
	arr, err := ExtractRoles(p)
	return len(arr) == 0, err
}

// ExtractRoles is a method to extract the list of database role for DDS service.
func ExtractRoles(p pagination.Page) ([]RoleResp, error) {
	var r ListResp
	err := (p.(RolePage)).ExtractInto(&r)
	if err != nil {
		return nil, err
	}
	var ur []RoleResp
	if r.Roles != "" {
		err = json.Unmarshal([]byte(r.Roles), &ur)
	}
	return ur, err
}
