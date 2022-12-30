package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteHostsGroupRequest struct {

	// region id
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 服务器组ID
	GroupId *string `json:"group_id,omitempty"`
}

func (o DeleteHostsGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteHostsGroupRequest struct{}"
	}

	return strings.Join([]string{"DeleteHostsGroupRequest", string(data)}, " ")
}
