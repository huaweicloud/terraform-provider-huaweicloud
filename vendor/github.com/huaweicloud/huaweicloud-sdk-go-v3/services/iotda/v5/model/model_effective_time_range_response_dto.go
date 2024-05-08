package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EffectiveTimeRangeResponseDto struct {

	// 设备代理开始生效的时间，使用UTC时区，格式：yyyyMMdd'T'HHmmss'Z'
	StartTime string `json:"start_time"`

	// 设备代理失效的时间，必须大于start_time，使用UTC时区，格式：yyyyMMdd'T'HHmmss'Z'
	EndTime string `json:"end_time"`
}

func (o EffectiveTimeRangeResponseDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EffectiveTimeRangeResponseDto struct{}"
	}

	return strings.Join([]string{"EffectiveTimeRangeResponseDto", string(data)}, " ")
}
