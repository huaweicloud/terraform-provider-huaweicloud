package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SnapshotSpecItems struct {

	// 子任务ID
	Id *string `json:"id,omitempty"`

	// 子任务类型
	Type *string `json:"type,omitempty"`

	// 状态
	Status *string `json:"status,omitempty"`

	// 任务创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 任务更新时间
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`

	// 信息
	Message *string `json:"message,omitempty"`
}

func (o SnapshotSpecItems) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SnapshotSpecItems struct{}"
	}

	return strings.Join([]string{"SnapshotSpecItems", string(data)}, " ")
}
