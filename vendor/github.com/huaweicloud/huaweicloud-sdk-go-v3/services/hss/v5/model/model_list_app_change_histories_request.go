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

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 排序的key值，目前只支持按照recent_scan_time排序，按照recent_scan_time排序时，根据sort_dir的值决定升序还是降序
	SortKey *string `json:"sort_key,omitempty"`

	// 排序方式，默认为降序，当sort_key为按照recent_scan_time排序时，根据当前值决定升序还是降序，当sort_key为其他值时均为降序   - asc ：升序   - desc ：降序
	SortDir *string `json:"sort_dir,omitempty"`

	// 每页显示数量
	Limit *int32 `json:"limit,omitempty"`

	// 偏移量：指定返回记录的开始位置
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
