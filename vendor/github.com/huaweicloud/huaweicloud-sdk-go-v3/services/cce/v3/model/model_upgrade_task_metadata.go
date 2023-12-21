package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeTaskMetadata 升级任务元数据
type UpgradeTaskMetadata struct {

	// 升级任务ID
	Uid *string `json:"uid,omitempty"`

	// 任务创建时间
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	// 任务更新时间
	UpdateTimestamp *string `json:"updateTimestamp,omitempty"`
}

func (o UpgradeTaskMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeTaskMetadata struct{}"
	}

	return strings.Join([]string{"UpgradeTaskMetadata", string(data)}, " ")
}
