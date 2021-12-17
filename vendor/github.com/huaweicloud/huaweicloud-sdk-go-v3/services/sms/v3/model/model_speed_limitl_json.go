package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 速率参数
type SpeedLimitlJson struct {
	// 时间段开始时间，格式：XX:XX。

	Start string `json:"start"`
	// 时间段结束时间，格式：XX:XX。

	End string `json:"end"`
	// 时间段的速率，0-1000的整数，单位：Mbit/s。

	Speed int32 `json:"speed"`
}

func (o SpeedLimitlJson) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SpeedLimitlJson struct{}"
	}

	return strings.Join([]string{"SpeedLimitlJson", string(data)}, " ")
}
