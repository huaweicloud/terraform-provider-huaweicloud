package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 逻辑卷组信息
type PutVolumeGroups struct {
	// lv信息

	LogicalVolumes *[]PutLogicalVolume `json:"logical_volumes,omitempty"`
	// 卷组id

	Id string `json:"id"`
	// 是否迁移

	NeedMigration *bool `json:"need_migration,omitempty"`
	// 调整大小

	AdjustSize *int64 `json:"adjust_size,omitempty"`
}

func (o PutVolumeGroups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PutVolumeGroups struct{}"
	}

	return strings.Join([]string{"PutVolumeGroups", string(data)}, " ")
}
