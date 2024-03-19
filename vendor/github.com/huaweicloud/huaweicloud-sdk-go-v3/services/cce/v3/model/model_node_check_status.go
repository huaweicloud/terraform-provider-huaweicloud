package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeCheckStatus 节点限制检查状态
type NodeCheckStatus struct {

	// 状态，取值如下 - Init: 初始化 - Running 运行中 - Success 成功 - Failed 失败
	Phase *string `json:"phase,omitempty"`

	// 节点检查状态
	NodeStageStatus *[]NodeStageStatus `json:"nodeStageStatus,omitempty"`
}

func (o NodeCheckStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeCheckStatus struct{}"
	}

	return strings.Join([]string{"NodeCheckStatus", string(data)}, " ")
}
