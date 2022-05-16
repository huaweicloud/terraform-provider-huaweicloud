package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListRolesForUserOnEnterpriseProjectRequest struct {

	// 企业项目ID。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 用户ID。
	UserId string `json:"user_id"`
}

func (o ListRolesForUserOnEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRolesForUserOnEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"ListRolesForUserOnEnterpriseProjectRequest", string(data)}, " ")
}
