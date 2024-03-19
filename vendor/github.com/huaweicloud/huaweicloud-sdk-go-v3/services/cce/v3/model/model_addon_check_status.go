package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddonCheckStatus 插件限制检查状态
type AddonCheckStatus struct {

	// 状态，取值如下 - Init: 初始化 - Running 运行中 - Success 成功 - Failed 失败
	Phase *string `json:"phase,omitempty"`

	// 检查项状态集合
	ItemsStatus *[]PreCheckItemStatus `json:"itemsStatus,omitempty"`
}

func (o AddonCheckStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddonCheckStatus struct{}"
	}

	return strings.Join([]string{"AddonCheckStatus", string(data)}, " ")
}
