package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StageKpiItem struct {

	// 比较符
	Comparison *string `json:"comparison,omitempty"`

	// 比较值
	Value *int32 `json:"value,omitempty"`
}

func (o StageKpiItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StageKpiItem struct{}"
	}

	return strings.Join([]string{"StageKpiItem", string(data)}, " ")
}
