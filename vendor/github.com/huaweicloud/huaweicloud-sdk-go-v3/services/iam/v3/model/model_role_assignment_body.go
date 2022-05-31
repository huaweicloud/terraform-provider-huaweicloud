package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RoleAssignmentBody struct {
	User *RoleUserAssignmentId `json:"user,omitempty"`

	Role *RoleAssignmentId `json:"role,omitempty"`

	Group *RoleGroupAssignmentId `json:"group,omitempty"`

	Agency *RoleAgencyAssignmentId `json:"agency,omitempty"`

	Scope *RoleAssignmentScope `json:"scope,omitempty"`

	// 是否基于所有项目授权。
	IsInherited *bool `json:"is_inherited,omitempty"`
}

func (o RoleAssignmentBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RoleAssignmentBody struct{}"
	}

	return strings.Join([]string{"RoleAssignmentBody", string(data)}, " ")
}
