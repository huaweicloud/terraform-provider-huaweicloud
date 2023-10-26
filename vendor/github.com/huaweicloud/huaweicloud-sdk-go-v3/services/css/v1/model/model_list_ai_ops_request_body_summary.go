package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAiOpsRequestBodySummary 风险概要。
type ListAiOpsRequestBodySummary struct {

	// 检测项判定为高风险的数量。
	High *int32 `json:"high,omitempty"`

	// 检测项判定为中风险的数量。
	Medium *int32 `json:"medium,omitempty"`

	// 检测项判定为建议的数量。
	Suggestion *int32 `json:"suggestion,omitempty"`
}

func (o ListAiOpsRequestBodySummary) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAiOpsRequestBodySummary struct{}"
	}

	return strings.Join([]string{"ListAiOpsRequestBodySummary", string(data)}, " ")
}
