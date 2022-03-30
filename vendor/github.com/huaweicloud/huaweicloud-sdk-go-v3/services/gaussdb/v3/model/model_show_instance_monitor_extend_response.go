package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowInstanceMonitorExtendResponse struct {
	// 实例秒级监控开关。为true时表示开启，为false时表示关闭。

	MonitorSwitch *bool `json:"monitor_switch,omitempty"`
	// 采集周期，仅在monitor_switch为true时返回。1：采集周期为1s； 5：采集周期为5s。

	Period         *int32 `json:"period,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowInstanceMonitorExtendResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceMonitorExtendResponse struct{}"
	}

	return strings.Join([]string{"ShowInstanceMonitorExtendResponse", string(data)}, " ")
}
