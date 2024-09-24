package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowReplayDelayStatusResponse Response Object
type ShowReplayDelayStatusResponse struct {

	// 当前配置的延迟时间，单位ms
	CurDelayTimeMills *int32 `json:"cur_delay_time_mills,omitempty"`

	// 延迟时间参数取值范围
	DelayTimeValueRange *string `json:"delay_time_value_range,omitempty"`

	// 真实延迟时间，单位ms
	RealDelayTimeMills *int32 `json:"real_delay_time_mills,omitempty"`

	// 当前日志回放状态。true表示回放暂停，false表示回放正常
	CurLogReplayPaused *bool `json:"cur_log_replay_paused,omitempty"`

	// 最新接收的日志
	LatestReceiveLog *string `json:"latest_receive_log,omitempty"`

	// 最新回放的日志位点
	LatestReplayLog *string `json:"latest_replay_log,omitempty"`
	HttpStatusCode  int     `json:"-"`
}

func (o ShowReplayDelayStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowReplayDelayStatusResponse struct{}"
	}

	return strings.Join([]string{"ShowReplayDelayStatusResponse", string(data)}, " ")
}
