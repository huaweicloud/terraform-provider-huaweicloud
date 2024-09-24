package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutoLaunchStatisticsRequest Request Object
type ListAutoLaunchStatisticsRequest struct {

	// 自启动项名称
	Name *string `json:"name,omitempty"`

	// 自启动项类型   - 0 ：自启动服务   - 1 ：定时任务   - 2 ：预加载动态库   - 3 ：Run注册表键   - 4 ：开机启动文件夹
	Type *string `json:"type,omitempty"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`
}

func (o ListAutoLaunchStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutoLaunchStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListAutoLaunchStatisticsRequest", string(data)}, " ")
}
