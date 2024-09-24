package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppStatisticsRequest Request Object
type ListAppStatisticsRequest struct {

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
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
