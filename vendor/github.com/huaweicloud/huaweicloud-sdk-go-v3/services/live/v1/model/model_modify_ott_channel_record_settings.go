package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyOttChannelRecordSettings OTT频道修改录制消息体
type ModifyOttChannelRecordSettings struct {

	// 频道推流域名
	Domain string `json:"domain"`

	// 组名或应用名
	AppName string `json:"app_name"`

	// 频道ID。频道唯一标识，为必填项
	Id string `json:"id"`

	RecordSettings *ModifyOttChannelRecordSettingsRecordSettings `json:"record_settings,omitempty"`
}

func (o ModifyOttChannelRecordSettings) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyOttChannelRecordSettings struct{}"
	}

	return strings.Join([]string{"ModifyOttChannelRecordSettings", string(data)}, " ")
}
