package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PageInfo struct {

	// 当前页事件、告警总数
	CurrentCount int32 `json:"current_count"`

	// 前一个marker
	PreviousMarker string `json:"previous_marker"`

	// 下一个marker
	NextMarker string `json:"next_marker"`
}

func (o PageInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PageInfo struct{}"
	}

	return strings.Join([]string{"PageInfo", string(data)}, " ")
}
