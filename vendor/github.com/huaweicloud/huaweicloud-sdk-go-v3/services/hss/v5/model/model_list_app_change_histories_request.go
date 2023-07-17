package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAppChangeHistoriesRequest Request Object
type ListAppChangeHistoriesRequest struct {

	// 主机id
	HostId *string `json:"host_id,omitempty"`

	// 主机ip
	HostIp *string `json:"host_ip,omitempty"`

	// 主机名称
	HostName *string `json:"host_name,omitempty"`

	// 软件名称
	AppName *string `json:"app_name,omitempty"`

	// 变更类型:   - add ：新建   - delete ：删除   - modify ：修改
	VariationType *string `json:"variation_type,omitempty"`

	// 企业项目
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 排序的key值
	SortKey *string `json:"sort_key,omitempty"`

	// 升序还是降序，默认升序，asc
	SortDir *string `json:"sort_dir,omitempty"`

	// 默认10
	Limit *int32 `json:"limit,omitempty"`

	// 默认是0
	Offset *int32 `json:"offset,omitempty"`

	// 变更开始时间，13位时间戳
	StartTime *int64 `json:"start_time,omitempty"`

	// 变更结束时间，13位时间戳
	EndTime *int64 `json:"end_time,omitempty"`
}

func (o ListAppChangeHistoriesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAppChangeHistoriesRequest struct{}"
	}

	return strings.Join([]string{"ListAppChangeHistoriesRequest", string(data)}, " ")
}
