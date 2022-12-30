package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PageInfoTagKeys struct {

	// 分页位置标识（索引）
	NextMarker string `json:"next_marker"`

	// 当前页标签键的数量
	CurrentCount int32 `json:"current_count"`
}

func (o PageInfoTagKeys) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PageInfoTagKeys struct{}"
	}

	return strings.Join([]string{"PageInfoTagKeys", string(data)}, " ")
}
