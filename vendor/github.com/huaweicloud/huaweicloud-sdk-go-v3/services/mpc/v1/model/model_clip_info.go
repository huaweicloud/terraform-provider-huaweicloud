package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ClipInfo struct {
	Input *ObsObjInfo `json:"input,omitempty"`

	// 剪切开始时间，单位：秒。可以有正数或者负数，正数表示从开始往后的时间点，负数表示从结尾往前的时间点。
	TimelineStart *string `json:"timeline_start,omitempty"`

	// 剪切结束时间，单位：秒。可以有正数或者负数，正数表示从开始往后的时间点，负数表示从结尾往前的时间点。
	TimelineEnd *string `json:"timeline_end,omitempty"`
}

func (o ClipInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClipInfo struct{}"
	}

	return strings.Join([]string{"ClipInfo", string(data)}, " ")
}
