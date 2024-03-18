package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PrecheckStatus 升级前检查状态
type PrecheckStatus struct {

	// 状态，取值如下 - Init: 初始化 - Running 运行中 - Success 成功 - Failed 失败 - Error 错误
	Phase *string `json:"phase,omitempty"`

	// 检查结果过期时间
	ExpireTimeStamp *string `json:"expireTimeStamp,omitempty"`

	// 信息，一般是执行错误的日志信息
	Message *string `json:"message,omitempty"`

	ClusterCheckStatus *ClusterCheckStatus `json:"clusterCheckStatus,omitempty"`

	AddonCheckStatus *AddonCheckStatus `json:"addonCheckStatus,omitempty"`

	NodeCheckStatus *NodeCheckStatus `json:"nodeCheckStatus,omitempty"`
}

func (o PrecheckStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrecheckStatus struct{}"
	}

	return strings.Join([]string{"PrecheckStatus", string(data)}, " ")
}
