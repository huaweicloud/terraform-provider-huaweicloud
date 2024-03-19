package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClusterCheckStatus 集群限制检查状态
type ClusterCheckStatus struct {

	// 状态，取值如下 - Init: 初始化 - Running 运行中 - Success 成功 - Failed 失败
	Phase *string `json:"phase,omitempty"`

	// 检查项状态集合
	ItemsStatus *[]PreCheckItemStatus `json:"itemsStatus,omitempty"`
}

func (o ClusterCheckStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterCheckStatus struct{}"
	}

	return strings.Join([]string{"ClusterCheckStatus", string(data)}, " ")
}
