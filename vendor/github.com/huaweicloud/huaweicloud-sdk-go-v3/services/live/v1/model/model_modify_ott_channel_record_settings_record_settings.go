package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelRecordSettingsRecordSettings 最大回看录制时长。在此时间段内会连续不断的录制，为必选项  单位：秒。取值为“0”时，表示不支持录制；最大支持录制14天
type ModifyOttChannelRecordSettingsRecordSettings struct {

	// 最大回看录制时长。在此时间段内会连续不断的录制，为必选项  单位：秒。取值为“0”时，表示不支持录制；最大支持录制14天。
	RollingbufferDuration *int32 `json:"rollingbuffer_duration,omitempty"`
}

func (o ModifyOttChannelRecordSettingsRecordSettings) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelRecordSettingsRecordSettings struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelRecordSettingsRecordSettings", string(data)}, " ")
}
