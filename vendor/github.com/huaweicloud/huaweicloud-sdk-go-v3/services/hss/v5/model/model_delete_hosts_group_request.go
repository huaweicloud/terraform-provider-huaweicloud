package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteHostsGroupRequest Request Object
type DeleteHostsGroupRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 服务器组ID
	GroupId string `json:"group_id"`
}

func (o DeleteHostsGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteHostsGroupRequest struct{}"
	}

	return strings.Join([]string{"DeleteHostsGroupRequest", string(data)}, " ")
}
