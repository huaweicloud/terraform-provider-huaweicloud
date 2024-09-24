package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHostGroupsRequest Request Object
type ListHostGroupsRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`
}

func (o ListHostGroupsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHostGroupsRequest struct{}"
	}

	return strings.Join([]string{"ListHostGroupsRequest", string(data)}, " ")
}
