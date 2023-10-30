package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AiOpsRiskInfo struct {

	// 检测项介绍。
	RiskType *string `json:"riskType,omitempty"`

	// 风险等级。 - high - medium - suggestion
	Level *string `json:"level,omitempty"`

	// 风险描述。
	Desc *string `json:"desc,omitempty"`

	// 风险建议。
	Suggestion *string `json:"suggestion,omitempty"`
}

func (o AiOpsRiskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AiOpsRiskInfo struct{}"
	}

	return strings.Join([]string{"AiOpsRiskInfo", string(data)}, " ")
}
