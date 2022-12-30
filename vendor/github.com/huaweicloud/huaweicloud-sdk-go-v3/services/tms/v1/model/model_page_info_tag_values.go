package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PageInfoTagValues struct {

	// 分页位置标识（索引）
	NextMarker string `json:"next_marker"`

	// 当前页标签值的数量
	CurrentCount int32 `json:"current_count"`
}

func (o PageInfoTagValues) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PageInfoTagValues struct{}"
	}

	return strings.Join([]string{"PageInfoTagValues", string(data)}, " ")
}
