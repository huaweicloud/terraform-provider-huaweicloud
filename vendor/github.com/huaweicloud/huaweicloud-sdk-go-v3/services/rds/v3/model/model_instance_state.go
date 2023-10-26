package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InstanceState 实例状态
type InstanceState struct {

	// 实例状态
	Status string `json:"status"`

	// 参数变更，是否需要重启
	WaitRestartForParams bool `json:"wait_restart_for_params"`
}

func (o InstanceState) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceState struct{}"
	}

	return strings.Join([]string{"InstanceState", string(data)}, " ")
}
