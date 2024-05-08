package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EffectiveTimeRange struct {

	// 设备代理开始生效的时间，使用UTC时区，格式：yyyyMMdd'T'HHmmss'Z'
	StartTime *string `json:"start_time,omitempty"`

	// 设备代理失效的时间，必须大于start_time，使用UTC时区，格式：yyyyMMdd'T'HHmmss'Z'
	EndTime *string `json:"end_time,omitempty"`
}

func (o EffectiveTimeRange) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EffectiveTimeRange struct{}"
	}

	return strings.Join([]string{"EffectiveTimeRange", string(data)}, " ")
}
