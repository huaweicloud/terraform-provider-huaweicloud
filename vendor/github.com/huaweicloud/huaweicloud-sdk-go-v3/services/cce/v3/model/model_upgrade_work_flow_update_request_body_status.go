package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeWorkFlowUpdateRequestBodyStatus 更新后workflow的状态（当前仅支持Cancel）
type UpgradeWorkFlowUpdateRequestBodyStatus struct {
	Phase *WorkFlowPhase `json:"phase,omitempty"`
}

func (o UpgradeWorkFlowUpdateRequestBodyStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeWorkFlowUpdateRequestBodyStatus struct{}"
	}

	return strings.Join([]string{"UpgradeWorkFlowUpdateRequestBodyStatus", string(data)}, " ")
}
