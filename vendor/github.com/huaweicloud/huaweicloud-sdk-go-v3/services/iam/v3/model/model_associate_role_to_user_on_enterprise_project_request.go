package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type AssociateRoleToUserOnEnterpriseProjectRequest struct {

	// 企业项目ID。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 用户ID。
	UserId string `json:"user_id"`

	// 权限ID。
	RoleId string `json:"role_id"`
}

func (o AssociateRoleToUserOnEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociateRoleToUserOnEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"AssociateRoleToUserOnEnterpriseProjectRequest", string(data)}, " ")
}
