package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeTaskStatus 升级任务状态信息
type UpgradeTaskStatus struct {

	// 升级任务状态. > Init：初始化 > Queuing：等待 > Running：运行中 > Pause：暂停 > Success：成功 > Failed：失败
	Phase *string `json:"phase,omitempty"`

	// 升级任务进度
	Progress *string `json:"progress,omitempty"`

	// 升级任务结束时间
	CompletionTime *string `json:"completionTime,omitempty"`
}

func (o UpgradeTaskStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeTaskStatus struct{}"
	}

	return strings.Join([]string{"UpgradeTaskStatus", string(data)}, " ")
}
