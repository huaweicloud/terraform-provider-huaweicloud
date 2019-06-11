package agency

import (
	"github.com/huaweicloud/golangsdk"
)

type CreateOpts struct {
	Name            string `json:"name" required:"true"`
	DomainID        string `json:"domain_id" required:"true"`
	DelegatedDomain string `json:"trust_domain_name" required:"true"`
	Description     string `json:"description,omitempty"`
}

type CreateOptsBuilder interface {
	ToAgencyCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToAgencyCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "agency")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAgencyCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

type UpdateOpts struct {
	DelegatedDomain string `json:"trust_domain_name,omitempty"`
	Description     string `json:"description,omitempty"`
}

type UpdateOptsBuilder interface {
	ToAgencyUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToAgencyUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "agency")
}

func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAgencyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r ErrResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}

func AttachRoleByProject(c *golangsdk.ServiceClient, agencyID, projectID, roleID string) (r ErrResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Put(roleURL(c, "projects", projectID, agencyID, roleID), nil, nil, reqOpt)
	return
}

func AttachRoleByDomain(c *golangsdk.ServiceClient, agencyID, domainID, roleID string) (r ErrResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Put(roleURL(c, "domains", domainID, agencyID, roleID), nil, nil, reqOpt)
	return
}

func DetachRoleByProject(c *golangsdk.ServiceClient, agencyID, projectID, roleID string) (r ErrResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Delete(roleURL(c, "projects", projectID, agencyID, roleID), reqOpt)
	return
}

func DetachRoleByDomain(c *golangsdk.ServiceClient, agencyID, domainID, roleID string) (r ErrResult) {
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{204}}
	_, r.Err = c.Delete(roleURL(c, "domains", domainID, agencyID, roleID), reqOpt)
	return
}

func ListRolesAttachedOnProject(c *golangsdk.ServiceClient, agencyID, projectID string) (r ListRolesResult) {
	_, r.Err = c.Get(listRolesURL(c, "projects", projectID, agencyID), &r.Body, nil)
	return
}

func ListRolesAttachedOnDomain(c *golangsdk.ServiceClient, agencyID, domainID string) (r ListRolesResult) {
	_, r.Err = c.Get(listRolesURL(c, "domains", domainID, agencyID), &r.Body, nil)
	return
}
