package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutoLaunchStatisticsRequest Request Object
type ListAutoLaunchStatisticsRequest struct {

	// 自启动项名称
	Name *string `json:"name,omitempty"`

	// 自启动项类型
	Type *string `json:"type,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListAutoLaunchStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutoLaunchStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListAutoLaunchStatisticsRequest", string(data)}, " ")
}
