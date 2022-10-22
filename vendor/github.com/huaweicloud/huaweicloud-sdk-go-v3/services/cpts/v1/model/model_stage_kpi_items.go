package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StageKpiItems struct {
	AverageResponseTime *StageKpiItem `json:"average_response_time,omitempty"`

	SuccessRate *StageKpiItem `json:"success_rate,omitempty"`
}

func (o StageKpiItems) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StageKpiItems struct{}"
	}

	return strings.Join([]string{"StageKpiItems", string(data)}, " ")
}
