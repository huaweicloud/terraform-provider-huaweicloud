package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type WorkFlowStatus struct {
	Phase *WorkFlowPhase `json:"phase,omitempty"`

	// 升级流程中的各个任务项的执行状态
	PointStatuses *[]PointStatus `json:"pointStatuses,omitempty"`

	// 表示该升级流程的任务执行线路
	LineStatuses *[]LineStatus `json:"lineStatuses,omitempty"`
}

func (o WorkFlowStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WorkFlowStatus struct{}"
	}

	return strings.Join([]string{"WorkFlowStatus", string(data)}, " ")
}
