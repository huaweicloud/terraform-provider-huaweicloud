package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 分页信息
type PageInfo struct {

	// 返回下一页的查询地址
	NextMarker *string `json:"next_marker,omitempty"`

	// 返回前一页查询地址
	PreviousMarker *string `json:"previous_marker,omitempty"`

	// 本页返回条目数量
	CurrentCount *int32 `json:"current_count,omitempty"`
}

func (o PageInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PageInfo struct{}"
	}

	return strings.Join([]string{"PageInfo", string(data)}, " ")
}
