package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleAssignmentScope struct {
	Project *RoleProjectAssignmentId `json:"project,omitempty"`

	Domain *RoleDomainAssignmentId `json:"domain,omitempty"`

	EnterpriseProject *RoleEnterpriseProjectAssignmentId `json:"enterprise_project,omitempty"`
}

func (o RoleAssignmentScope) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleAssignmentScope struct{}"
	}

	return strings.Join([]string{"RoleAssignmentScope", string(data)}, " ")
}
