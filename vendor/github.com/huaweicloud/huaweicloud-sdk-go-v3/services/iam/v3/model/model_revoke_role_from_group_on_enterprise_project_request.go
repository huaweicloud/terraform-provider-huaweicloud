package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type RevokeRoleFromGroupOnEnterpriseProjectRequest struct {

	// 企业项目ID。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 用户组ID。
	GroupId string `json:"group_id"`

	// 权限ID。
	RoleId string `json:"role_id"`
}

func (o RevokeRoleFromGroupOnEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeRoleFromGroupOnEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"RevokeRoleFromGroupOnEnterpriseProjectRequest", string(data)}, " ")
}
