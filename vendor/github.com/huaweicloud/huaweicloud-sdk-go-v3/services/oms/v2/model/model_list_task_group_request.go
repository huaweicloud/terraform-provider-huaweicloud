package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListTaskGroupRequest struct {

	// 查询返回迁移组任务列表当前页面的数量，默认查询10条。 最多返回100条迁移任务信息。
	Limit *int32 `json:"limit,omitempty"`

	// 起始的任务序号，默认为0。 取值大于等于0，取值为0时从第一条开始查询。
	Offset *int32 `json:"offset,omitempty"`

	// 迁移任务组状态（无该参数时代表查询所有状态的任务） 0 – 等待中 1 – 执行中/创建中 2 – 监控任务执行 3 – 暂停 4 – 创建任务失败 5 – 迁移失败 6 – 迁移完成 7 – 暂停中 8 – 等待删除中 9 – 删除
	Status *int32 `json:"status,omitempty"`
}

func (o ListTaskGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTaskGroupRequest struct{}"
	}

	return strings.Join([]string{"ListTaskGroupRequest", string(data)}, " ")
}
