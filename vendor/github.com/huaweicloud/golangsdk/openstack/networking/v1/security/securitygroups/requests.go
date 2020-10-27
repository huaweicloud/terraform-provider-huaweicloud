package securitygroups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type CreateOpts struct {

	// Specifies the security group name.
	Name string `json:"name" required:"true"`

	// Specifies the enterprise project ID. This field can be used to
	// filter out the VPCs associated with a specified enterprise project.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`

	// Specifies the resource ID of the VPC to which the security
	// group belongs.
	VpcId string `json:"vpc_id,omitempty"`

	// Specifies the default security group rule, which ensures that
	//// hosts in the security group can communicate with one another.
	//SecurityGroupRules []SecurityGroupRule `json:"security_group_rules"`

	Description string `json:"description,omitempty"`
}

type CreateOptsBuilder interface {
	ToSecuritygroupsCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToSecuritygroupsCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "security_group")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecuritygroupsCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(CreateURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Delete(client *golangsdk.ServiceClient, securityGroupId string) (r DeleteResult) {
	url := DeleteURL(client, securityGroupId)
	_, r.Err = client.Delete(url, nil)
	return
}

func Get(client *golangsdk.ServiceClient, securityGroupId string) (r GetResult) {
	url := GetURL(client, securityGroupId)
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

type ListOpts struct {

	// Specifies the resource ID of pagination query. If the parameter
	// is left blank, only resources on the first page are queried.
	Marker string `q:"marker"`

	// Specifies the number of records returned on each page.
	Limit int `q:"limit"`

	// Specifies the VPC ID used as the query filter.
	VpcId string `q:"vpc_id"`

	// enterprise_project_id
	// Specifies the enterprise_project_id used as the query filter.
	EnterpriseProjectId string `q:"enterprise_project_id"`
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := ListURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url,
		func(r pagination.PageResult) pagination.Page {
			return SecurityGroupPage{pagination.LinkedPageBase{PageResult: r}}

		})
}
