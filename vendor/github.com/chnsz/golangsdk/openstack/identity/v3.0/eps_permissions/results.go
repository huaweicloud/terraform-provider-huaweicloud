package eps_permissions

import (
	"github.com/chnsz/golangsdk"
)

// CommonResult is used to accept error.
type CommonResult struct {
	golangsdk.ErrResult
}

type RoleResult struct {
	golangsdk.Result
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"display_name"`
	Catalog     string `json:"catalog"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Policy      Policy `json:"policy"`
	DomainId    string `json:"domain_id"`
	References  int    `json:"references"`
}

type Policy struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Action    []string               `json:"Action"`
	Effect    string                 `json:"Effect"`
	Condition map[string]interface{} `json:"Condition"`
	Resource  interface{}            `json:"Resource"`
}

func (r RoleResult) Extract() ([]Role, error) {
	var s struct {
		Roles []Role `json:"roles"`
	}
	err := r.ExtractInto(&s)
	return s.Roles, err
}
