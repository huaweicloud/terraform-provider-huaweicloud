package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeInfoStatus 升级状态信息
type UpgradeInfoStatus struct {

	// 升级任务状态. > Init：初始化 > Running：运行中 > Pause：暂停 > Success：成功 > Failed：失败
	Phase *string `json:"phase,omitempty"`

	// 升级任务进度
	Progress *string `json:"progress,omitempty"`

	// 升级任务结束时间
	CompletionTime *string `json:"completionTime,omitempty"`
}

func (o UpgradeInfoStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeInfoStatus struct{}"
	}

	return strings.Join([]string{"UpgradeInfoStatus", string(data)}, " ")
}
