package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CurrentNodeDetail 当前正在升级节点任务详情。
type CurrentNodeDetail struct {

	// 升级任务序号。
	Order *int32 `json:"order,omitempty"`

	// 升级任务名称。
	Name *string `json:"name,omitempty"`

	// 当前任务状态。
	Status *string `json:"status,omitempty"`

	// 当前任务描述。
	Desc *string `json:"desc,omitempty"`

	// 当前任务开始时间。
	BeginTime *string `json:"beginTime,omitempty"`

	// 当前任务结束时间。
	EndTime *string `json:"endTime,omitempty"`
}

func (o CurrentNodeDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CurrentNodeDetail struct{}"
	}

	return strings.Join([]string{"CurrentNodeDetail", string(data)}, " ")
}
