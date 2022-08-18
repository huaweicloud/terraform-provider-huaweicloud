package groups

import (
	"github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
	"github.com/chnsz/golangsdk/pagination"
)

// SecurityGroup is a struct that represents the detail of the security group.
type SecurityGroup struct {
	// Specifies the security group name.
	Name string `json:"name"`
	// Provides supplementary information about the security group.
	Description string `json:"description"`
	// Specifies the security group ID, which uniquely identifies the security group.
	ID string `json:"id"`
	// Specifies the resource ID of the VPC to which the security group belongs.
	// Note: This parameter has been discarded, it is not recommended to use it.
	VpcId string `json:"vpc_id"`
	// Enterprise project ID, the default value is 0.
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// Specifies the default security group rules, which ensure that resources in the security group can communicate
	// with one another.
	SecurityGroupRules []rules.SecurityGroupRule `json:"security_group_rules"`
	// ID of the project to which the security group rule belongs.
	ProjectId string `json:"project_id"`
	// Security group creation time, in UTC format: yyyy-MM-ddTHH:mm:ss.
	CreatedAt string `json:"created_at"`
	// Security group rule udpate time, in UTC format: yyyy-MM-ddTHH:mm:ss.
	UpdatedAt string `json:"updated_at"`
}

type SecurityGroupPage struct {
	pagination.MarkerPageBase
}

// LastMarker method returns the last security group ID in a security group page.
func (p SecurityGroupPage) LastMarker() (string, error) {
	secgroups, err := ExtractSecurityGroups(p)
	if err != nil {
		return "", err
	}
	if len(secgroups) == 0 {
		return "", nil
	}
	return secgroups[len(secgroups)-1].ID, nil
}

// IsEmpty method checks whether the current security group page is empty.
func (p SecurityGroupPage) IsEmpty() (bool, error) {
	secgroups, err := ExtractSecurityGroups(p)
	return len(secgroups) == 0, err
}

// ExtractSecurityGroups is a method to extract the list of security group details.
func ExtractSecurityGroups(r pagination.Page) ([]SecurityGroup, error) {
	var s []SecurityGroup
	err := r.(SecurityGroupPage).Result.ExtractIntoSlicePtr(&s, "security_groups")
	return s, err
}
