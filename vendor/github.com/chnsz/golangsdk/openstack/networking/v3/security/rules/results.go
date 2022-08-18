package rules

import "github.com/chnsz/golangsdk/pagination"

// SecurityGroupRule is a struct that represents the detail of the security group rule.
type SecurityGroupRule struct {
	// Specifies the security group rule ID, which uniquely identifies the security group rule.
	ID string `json:"id"`
	// Provides supplementary information about the security group rule.
	// The value can contain no more than 255 characters, including letters and digits.
	Description string `json:"description"`
	// Specifies the security group ID, which uniquely identifies the security group.
	SecurityGroupId string `json:"security_group_id"`
	// Specifies the direction of access control.
	// Possible values are as follows:
	//   egress
	//   ingress
	Direction string `json:"direction"`
	// Specifies the protocol type. The value can be icmp, tcp, udp, icmpv6 or protocol number.
	// If the parameter is left blank, all protocols are supported.
	// When the protocol is icmpv6, the network type should be IPv6.
	// when the protocol is icmp, the network type should be IPv4.
	Protocol string `json:"protocol"`
	// Specifies the IP protocol version. The value can be IPv4 or IPv6.
	Ethertype string `json:"ethertype"`
	// Port value range.
	// Value range: support and single port (80), contiguous port (1-30) and discontinuous port (22, 3389, 80).
	MultiPort string `json:"multiport"`
	// Security group rules take effect policy.
	//   allow
	//   deny
	// Default is deny.
	Action string `json:"action"`
	// Priority
	// Value range: 1~100, 1 represents the highest priority.
	Priority int `json:"priority"`
	// Specifies the remote IP address.
	// If the access control direction is set to egress, the parameter specifies the source IP address.
	// If the access control direction is set to ingress, the parameter specifies the destination IP address.
	// The value can be in the CIDR format or IP addresses.
	// The parameter is exclusive with parameter remote_group_id.
	RemoteIpPrefix string `json:"remote_ip_prefix"`
	// Specifies the ID of the peer security group.
	// The value is exclusive with parameter remote_ip_prefix.
	RemoteGroupId string `json:"remote_group_id"`
	// Remote address group ID. Mutually exclusive with remote_ip_prefix, remote_group_id parameters.
	RemoteAddressGroupId string `json:"remote_address_group_id"`
	// Security group rule creation time, in UTC format: yyyy-MM-ddTHH:mm:ss.
	CreateAt string `json:"created_at"`
	// Security group rule udpate time, in UTC format: yyyy-MM-ddTHH:mm:ss.
	UpdateAt string `json:"updated_at"`
	// ID of the project to which the security group rule belongs.
	ProjectId string `json:"project_id"`
}

type SecurityGroupRulePage struct {
	pagination.MarkerPageBase
}

// LastMarker method returns the last security group rule ID in a SecurityGroupRulePage.
func (p SecurityGroupRulePage) LastMarker() (string, error) {
	secgroups, err := ExtractSecurityGroupRules(p)
	if err != nil {
		return "", err
	}
	if len(secgroups) == 0 {
		return "", nil
	}
	return secgroups[len(secgroups)-1].ID, nil
}

// IsEmpty method checks whether the current SecurityGroupRulePage is empty.
func (p SecurityGroupRulePage) IsEmpty() (bool, error) {
	secgroups, err := ExtractSecurityGroupRules(p)
	return len(secgroups) == 0, err
}

// ExtractSecurityGroupRules is a method to extract the list of security group rule details.
func ExtractSecurityGroupRules(r pagination.Page) ([]SecurityGroupRule, error) {
	var s []SecurityGroupRule
	err := r.(SecurityGroupRulePage).Result.ExtractIntoSlicePtr(&s, "security_group_rules")
	return s, err
}
