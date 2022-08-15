package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowHistoryTaskDetailsRequest struct {

	// 当用户开启企业项目功能时，该参数生效，表示查询资源所属项目，\"all\"表示所有项目。注意：当使用子账号调用接口时，该参数必传。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 刷新任务ID。
	HistoryTasksId string `json:"history_tasks_id"`

	// 刷新预热的urls所显示单页最大数量，取值范围为1-10000。page_size和page_number必须同时传值。默认值30。
	PageSize *int32 `json:"page_size,omitempty"`

	// 刷新预热的urls当前查询为第几页，取值范围为1-65535。默认值1。
	PageNumber *int32 `json:"page_number,omitempty"`

	// url的状态 processing 处理中，succeed 完成，failed 失败，waiting 等待，refreshing 刷新中，preheating 预热中。
	Status *string `json:"status,omitempty"`

	// url的地址。
	Url *string `json:"url,omitempty"`

	// 刷新预热任务的创建时间。不传参默认为查询7天内的任务。最长可查询15天内数据。
	CreateTime *int64 `json:"create_time,omitempty"`
}

func (o ShowHistoryTaskDetailsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowHistoryTaskDetailsRequest struct{}"
	}

	return strings.Join([]string{"ShowHistoryTaskDetailsRequest", string(data)}, " ")
}
