package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RecoveryCondition struct {

	// 告警恢复周期的个数。
	RecoveryTimeframe *int32 `json:"recovery_timeframe,omitempty"`
}

func (o RecoveryCondition) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecoveryCondition struct{}"
	}

	return strings.Join([]string{"RecoveryCondition", string(data)}, " ")
}
