package agency

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/roles"
)

type Agency struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	DomainID            string `json:"domain_id"`
	DelegatedDomainID   string `json:"trust_domain_id"`
	DelegatedDomainName string `json:"trust_domain_name"`
	Description         string `json:"description"`
	Duration            string `json:"duration"`
	ExpireTime          string `json:"expire_time"`
	CreateTime          string `json:"create_time"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Agency, error) {
	var s struct {
		Agency *Agency `json:"agency"`
	}
	err := r.ExtractInto(&s)
	return s.Agency, err
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type ErrResult struct {
	golangsdk.ErrResult
}

type ListRolesResult struct {
	golangsdk.Result
}

func (r ListRolesResult) ExtractRoles() ([]roles.Role, error) {
	var s struct {
		Roles []roles.Role `json:"roles"`
	}
	err := r.ExtractInto(&s)
	return s.Roles, err
}
