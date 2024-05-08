package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateOttChannelInfoReqRecordSettings 最大回看录制时长。在此时间段内会连续不断的录制，为必选项  单位：秒。取值为“0”时，表示不支持录制；最大支持录制14天
type CreateOttChannelInfoReqRecordSettings struct {

	// 最大回看录制时长。在此时间段内会连续不断的录制，为必选项  单位：秒。取值为“0”时，表示不支持录制；最大支持录制14天
	RollingbufferDuration int32 `json:"rollingbuffer_duration"`
}

func (o CreateOttChannelInfoReqRecordSettings) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateOttChannelInfoReqRecordSettings struct{}"
	}

	return strings.Join([]string{"CreateOttChannelInfoReqRecordSettings", string(data)}, " ")
}
