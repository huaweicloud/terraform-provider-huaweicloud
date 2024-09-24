package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SwitchLogReplayRequestBody 中止/恢复wal日志回放请求体
type SwitchLogReplayRequestBody struct {

	// “true”表示中止回放，“false”表示恢复回放，其他情况表示不做操作
	PauseLogReplay string `json:"pause_log_replay"`
}

func (o SwitchLogReplayRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchLogReplayRequestBody struct{}"
	}

	return strings.Join([]string{"SwitchLogReplayRequestBody", string(data)}, " ")
}
