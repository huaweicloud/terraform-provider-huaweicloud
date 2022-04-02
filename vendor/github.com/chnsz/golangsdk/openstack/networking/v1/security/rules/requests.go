package rules

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is a struct which will be used to create a new security group rule.
type CreateOpts struct {
	// Specifies the security group ID.
	SecurityGroupId string `json:"security_group_id" required:"true"`
	// Provides supplementary information about the security group rule.
	// The value can contain no more than 255 characters, including letters and digits.
	Description string `json:"description,omitempty"`
	// Specifies the direction of access control.
	// Possible values are as follows:
	//   egress
	//   ingress
	Direction string `json:"direction" required:"true"`
	// Specifies the IP protocol version. The value can be IPv4 or IPv6. The default value is IPv4.
	Ethertype string `json:"ethertype,omitempty"`
	// Specifies the protocol type. The value can be icmp, tcp, or udp.
	// If the parameter is left blank, all protocols are supported.
	Protocol string `json:"protocol,omitempty"`
	// Specifies the start port number.
	// The value ranges from 1 to 65535.
	// The value cannot be greater than the port_range_max value. An empty value indicates all ports.
	PortRangeMin int `json:"port_range_min,omitempty"`
	// Specifies the end port number.
	// The value ranges from 1 to 65535.
	// The value cannot be smaller than the port_range_min value. An empty value indicates all ports.
	PortRangeMax int `json:"port_range_max,omitempty"`
	// Specifies the remote IP address.
	// If the access control direction is set to egress, the parameter specifies the source IP address.
	// If the access control direction is set to ingress, the parameter specifies the destination IP address.
	// The value can be in the CIDR format or IP addresses.
	// The parameter is exclusive with parameter remote_group_id.
	RemoteIpPrefix string `json:"remote_ip_prefix,omitempty"`
	// Specifies the ID of the peer security group.
	// The value is exclusive with parameter remote_ip_prefix.
	RemoteGroupId string `json:"remote_group_id,omitempty"`
}

// Create is a method to create a new security group rule.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) (*SecurityGroupRule, error) {
	b, err := golangsdk.BuildRequestBody(opts, "security_group_rule")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, nil)
	if err == nil {
		var r SecurityGroupRule
		rst.ExtractIntoStructPtr(&r, "security_group_rule")
		return &r, nil
	}
	return nil, err
}

// Get is a method to obtain the security group rule detail.
func Get(c *golangsdk.ServiceClient, ruleId string) (*SecurityGroupRule, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, ruleId), &rst.Body, nil)
	if err == nil {
		var r SecurityGroupRule
		rst.ExtractIntoStructPtr(&r, "security_group_rule")
		return &r, nil
	}
	return nil, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Specifies a resource ID for pagination query, indicating that the query starts from the next record of the
	// specified resource ID. This parameter can work together with the parameter limit.
	//   1. If parameters marker and limit are not passed, all resource records will be returned.
	//   2. If the parameter marker is not passed and the value of parameter limit is set to 10, the first 10 resource
	//     records will be returned.
	//   3. If the value of the parameter marker is set to the resource ID of the 10th record and the value of parameter
	//     limit is set to 10, the 11th to 20th resource records will be returned.
	//   4. If the value of the parameter marker is set to the resource ID of the 10th record and the parameter limit is
	//     not passed, resource records starting from the 11th records (including 11th) will be returned.
	Marker string `q:"marker"`
	// Specifies the number of records that will be returned on each page. The value is from 0 to intmax.
	// limit can be used together with marker. For details, see the parameter description of marker.
	Limit int `q:"limit"`
	// Specifies the security group ID.
	SecurityGroupId string `q:"security_group_id"`
	// Specifies the project ID.
	ProjectId string `q:"project_id"`
}

// List is a method to obtain the list of the security group rules.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]SecurityGroupRule, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := SecurityGroupRulePage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractSecurityGroupRules(pages)
}

// Delete is a method to delete an existing security group rule.
func Delete(c *golangsdk.ServiceClient, securityGroupRuleId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(resourceURL(c, securityGroupRuleId), nil)
	return &r
}
