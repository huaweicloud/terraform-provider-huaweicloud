package groups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is a struct which will be used to create a new security group.
type CreateOpts struct {
	// Specifies the security group name. The value can contain 1 to 64 characters, including letters, digits,
	// underscores (_), hyphens (-), and periods (.).
	Name string `json:"name" required:"true"`
	// Specifies the resource ID of the VPC to which the security group belongs.
	// Note: This parameter has been discarded, it is not recommended to use it.
	VpcId string `json:"vpc_id,omitempty"`
	// Enterprise project ID. The default value is 0.
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

// Create is a method to create a new security group.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*SecurityGroup, error) {
	b, err := golangsdk.BuildRequestBody(opts, "security_group")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, nil)
	if err == nil {
		var r SecurityGroup
		err := rst.ExtractIntoStructPtr(&r, "security_group")
		return &r, err
	}
	return nil, err
}

// Get is a method to obtain the security group detail.
func Get(c *golangsdk.ServiceClient, securityGroupId string) (*SecurityGroup, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, securityGroupId), &rst.Body, nil)
	if err == nil {
		var r SecurityGroup
		err = rst.ExtractIntoStructPtr(&r, "security_group")
		return &r, err
	}
	return nil, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Specifies the number of records that will be returned on each page. The value is from 0 to intmax.
	// limit can be used together with marker. For details, see the parameter description of marker.
	Limit int `q:"limit"`
	// Specifies a resource ID for pagination query, indicating that the query starts from the next record of the
	// specified resource ID. This parameter can work together with the parameter limit.
	//   If parameters marker and limit are not passed, all resource records will be returned.
	//   If the parameter marker is not passed and the value of parameter limit is set to 10, the first 10 resource
	//     records will be returned.
	//   If the value of the parameter marker is set to the resource ID of the 10th record and the value of parameter
	//     limit is set to 10, the 11th to 20th resource records will be returned.
	//   If the value of the parameter marker is set to the resource ID of the 10th record and the parameter limit is
	//     not passed, resource records starting from the 11th records (including 11th) will be returned.
	Marker string `q:"marker"`
	// Security group ID. You can use this field to filter security groups precisely, supporting multiple IDs.
	ID string `q:"id"`
	// Security group name. You can use this field to accurately filter security groups that meet the conditions, and
	// support incoming multiple name filters.
	Name string `q:"name"`
	// Security group description. You can use this field to filter security groups precisely, and support multiple
	// descriptions for filtering.
	Description string `q:"description"`
	// Enterprise project ID. Maximum length 36 bytes, UUID format with "-" hyphen, or string "0" (default).
	// If the enterprise project sub-account needs to query all security groups under the enterprise project, please
	// input 'all_granted_eps'.
	EnterpriseProjectId string `q:"enterprise_project_id"`
}

// List is a method to obtain the list of the security groups.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]SecurityGroup, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := SecurityGroupPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractSecurityGroups(pages)
}

// UpdateOpts is a struct which will be used to update the existing security group using given parameters.
type UpdateOpts struct {
	// Specifies the security group name. The value can contain 1 to 64 characters, including letters, digits,
	// underscores (_), hyphens (-), and periods (.).
	Name string `json:"name,omitempty"`
	// Specifies the description of the security group.
	// The angle brackets (< and >) are not allowed for the description.
	Description *string `json:"description,omitempty"`
}

// Update is a method to update an existing security group.
func Update(c *golangsdk.ServiceClient, securityGroupId string, opts UpdateOpts) (*SecurityGroup, error) {
	b, err := golangsdk.BuildRequestBody(opts, "security_group")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, securityGroupId), b, &rst.Body, nil)
	if err == nil {
		var r SecurityGroup
		err = rst.ExtractIntoStructPtr(&r, "security_group")
		return &r, err
	}
	return nil, err
}

// Delete is a method to delete an existing security group.
func Delete(c *golangsdk.ServiceClient, securityGroupId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, securityGroupId), nil)
	return &r
}

// IDFromName is a convenience function that returns a securtiy group's ID, given its name.
func IDFromName(client *golangsdk.ServiceClient, name string) (string, error) {
	var count int
	var id string
	opt := ListOpts{
		Name:                name,
		EnterpriseProjectId: "all_granted_eps",
	}
	secgroupList, err := List(client, opt)
	if err != nil {
		return "", err
	}
	for _, sg := range secgroupList {
		if sg.Name == name {
			count++
			id = sg.ID
		}
	}

	switch count {
	case 0:
		return "", golangsdk.ErrResourceNotFound{Name: name, ResourceType: "Security Group"}
	case 1:
		return id, nil
	default:
		return "", golangsdk.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "Security Group"}
	}
}
