package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppStatisticsRequest Request Object
type ListAppStatisticsRequest struct {

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量，为页数*每页显示条数
	Offset *int32 `json:"offset,omitempty"`

	// 类别，默认为host，包含如下： - host：主机 - container：容器
	Category *string `json:"category,omitempty"`
}

func (o ListAppStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListAppStatisticsRequest", string(data)}, " ")
}
