package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListMigprojectsRequest struct {
	// 每一页记录的迁移项目

	Limit *int32 `json:"limit,omitempty"`
	// 偏移量

	Offset *int32 `json:"offset,omitempty"`
}

func (o ListMigprojectsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMigprojectsRequest struct{}"
	}

	return strings.Join([]string{"ListMigprojectsRequest", string(data)}, " ")
}
