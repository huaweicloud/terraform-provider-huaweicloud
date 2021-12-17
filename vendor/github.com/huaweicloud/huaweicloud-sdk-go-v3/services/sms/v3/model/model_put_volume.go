package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 更新分区
type PutVolume struct {
	// 数据库id

	Id *string `json:"id,omitempty"`
	// 是否迁移

	NeedMigration *bool `json:"need_migration,omitempty"`
	// 调整大小

	AdjustSize *int64 `json:"adjust_size,omitempty"`
}

func (o PutVolume) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PutVolume struct{}"
	}

	return strings.Join([]string{"PutVolume", string(data)}, " ")
}
