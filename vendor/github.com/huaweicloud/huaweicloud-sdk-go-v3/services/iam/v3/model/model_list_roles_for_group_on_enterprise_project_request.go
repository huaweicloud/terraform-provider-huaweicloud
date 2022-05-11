package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListRolesForGroupOnEnterpriseProjectRequest struct {

	// 待查询企业项目ID。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// 待查询用户组。
	GroupId string `json:"group_id"`
}

func (o ListRolesForGroupOnEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRolesForGroupOnEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"ListRolesForGroupOnEnterpriseProjectRequest", string(data)}, " ")
}
