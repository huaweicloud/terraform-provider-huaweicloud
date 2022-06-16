package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTasksRequest struct {

	// 迁移任务组group_id
	GroupId *string `json:"group_id,omitempty"`

	// 查询返回迁移任务列表当前页面的数量，默认查询10条。 最多返回100条迁移任务信息。
	Limit *int32 `json:"limit,omitempty"`

	// 起始的任务序号，默认为0。 取值大于等于0，取值为0时从第一条开始查询。
	Offset *int32 `json:"offset,omitempty"`

	// 迁移任务状态（无该参数时代表查询所有状态的任务）： 1：等待调度 2：正在执行 3：停止 4：失败 5：成功
	Status *int32 `json:"status,omitempty"`
}

func (o ListTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTasksRequest struct{}"
	}

	return strings.Join([]string{"ListTasksRequest", string(data)}, " ")
}
