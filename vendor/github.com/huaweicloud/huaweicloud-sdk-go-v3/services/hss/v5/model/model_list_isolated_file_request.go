package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListIsolatedFileRequest Request Object
type ListIsolatedFileRequest struct {

	// region id
	Region string `json:"region"`

	// 租户企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 查询时间范围天数，与自定义查询时间begin_time，end_time互斥
	LastDays *int32 `json:"last_days,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 隔离状态，包含如下:   - isolated : 已隔离   - restored : 已恢复   - isolating : 已下发隔离任务   - restoring : 已下发恢复任务
	IsolationStatus *string `json:"isolation_status,omitempty"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListIsolatedFileRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListIsolatedFileRequest struct{}"
	}

	return strings.Join([]string{"ListIsolatedFileRequest", string(data)}, " ")
}
