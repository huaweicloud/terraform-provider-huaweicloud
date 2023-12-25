package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSyncTasksRequest Request Object
type ListSyncTasksRequest struct {

	// 查询返回同步任务列表当前页面的数量，默认查询10条。 最多返回100条迁移任务信息。
	Limit *int32 `json:"limit,omitempty"`

	// 起始的任务序号，默认为0。 取值大于等于0，取值为0时从第一条开始查询。
	Offset *int32 `json:"offset,omitempty"`

	// 同步任务状态（无该参数时代表查询所有状态的任务）： SYNCHRONIZING：同步中 STOPPED：已停止
	Status *string `json:"status,omitempty"`
}

func (o ListSyncTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSyncTasksRequest struct{}"
	}

	return strings.Join([]string{"ListSyncTasksRequest", string(data)}, " ")
}
