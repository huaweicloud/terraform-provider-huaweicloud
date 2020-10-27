package securitygroups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type SecurityGroup struct {
	// Specifies the security group name.
	Name string `json:"name"`

	// Provides supplementary information about the security group.
	Description string `json:"description"`

	// Specifies the security group ID, which uniquely identifies the
	// security group.
	ID string `json:"id"`

	// Specifies the resource ID of the VPC to which the security
	// group belongs.
	VpcId string `json:"vpc_id"`

	// Specifies the default security group rule, which ensures that
	// hosts in the security group can communicate with one another.
	SecurityGroupRules []SecurityGroupRule `json:"security_group_rules"`

	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
}

type SecurityGroupRule struct {
	// Specifies the security group rule ID.
	ID string `json:"id,omitempty"`

	// Specifies the description.
	Description string `json:"description,omitempty"`

	// Specifies the security group ID.
	SecurityGroupId string `json:"security_group_id,omitempty"`

	// Specifies the direction of access control. The value can
	// be?egress?or?ingress.
	Direction string `json:"direction,omitempty"`

	// Specifies the version of the Internet Protocol. The value can
	// be?IPv4?or?IPv6.
	Ethertype string `json:"ethertype,omitempty"`

	// Specifies the protocol type. If the parameter is left blank,
	// the security group supports all types of protocols. The value can be?icmp,?tcp,
	// or?udp.
	Protocol string `json:"protocol,omitempty"`

	// Specifies the start port. The value ranges from 1 to 65,535.
	// The value must be less than or equal to the value of?port_range_max. An empty value
	// indicates all ports. If?protocol?is?icmp, the value range is determined by the
	// ICMP-port range relationship table provided in Appendix A.2.
	PortRangeMin *int `json:"port_range_min,omitempty"`

	// Specifies the end port. The value ranges from 1 to 65,535. The
	// value must be greater than or equal to the value of?port_range_min. An empty value
	// indicates all ports. If?protocol?is?icmp, the value range is determined by the
	// ICMP-port range relationship table provided in Appendix A.2.
	PortRangeMax *int `json:"port_range_max,omitempty"`

	// Specifies the remote IP address. If the access control
	// direction is set to?egress, the parameter specifies the source IP address. If the
	// access control direction is set to?ingress, the parameter specifies the destination
	// IP address. The parameter is exclusive with parameter?remote_group_id. The value can
	// be in the CIDR format or IP addresses.
	RemoteIpPrefix string `json:"remote_ip_prefix,omitempty"`

	// Specifies the ID of the peer security group. The value is
	// exclusive with parameter?remote_ip_prefix.
	RemoteGroupId string `json:"remote_group_id,omitempty"`
}

type CreateResult struct {
	commonResult
}

func (r CreateResult) Extract() (*SecurityGroup, error) {
	var entity SecurityGroup
	err := r.ExtractIntoStructPtr(&entity, "security_group")
	return &entity, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	commonResult
}

func (r GetResult) Extract() (*SecurityGroup, error) {
	var entity SecurityGroup
	err := r.ExtractIntoStructPtr(&entity, "security_group")
	return &entity, err
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() (*[]SecurityGroup, error) {
	var list []SecurityGroup
	err := r.ExtractIntoSlicePtr(&list, "security_groups")
	return &list, err
}

func (r SecurityGroupPage) IsEmpty() (bool, error) {
	list, err := ExtractSecurityGroups(r)
	return len(list) == 0, err
}

type SecurityGroupPage struct {
	pagination.LinkedPageBase
}

func ExtractSecurityGroups(r pagination.Page) ([]SecurityGroup, error) {
	var s struct {
		SecurityGroups []SecurityGroup `json:"security_groups"`
	}
	err := r.(SecurityGroupPage).ExtractInto(&s)
	return s.SecurityGroups, err
}

func (r SecurityGroupPage) NextPageURL() (string, error) {
	s, err := ExtractSecurityGroups(r)
	if err != nil {
		return "", err
	}
	return r.WrapNextPageURL(s[len(s)-1].ID)
}
