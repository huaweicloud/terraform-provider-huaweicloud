package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListGroupsForEnterpriseProjectRequest struct {

	// 待查询的企业项目Id。
	EnterpriseProjectId string `json:"enterprise_project_id"`
}

func (o ListGroupsForEnterpriseProjectRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListGroupsForEnterpriseProjectRequest struct{}"
	}

	return strings.Join([]string{"ListGroupsForEnterpriseProjectRequest", string(data)}, " ")
}
