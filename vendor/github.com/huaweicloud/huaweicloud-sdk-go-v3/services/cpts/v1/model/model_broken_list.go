package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BrokenList struct {
	BrandBrokens *BrandBrokens `json:"brand_brokens,omitempty"`

	// 时间戳
	CommonTimestamps *[]string `json:"commonTimestamps,omitempty"`

	// 摸高数据
	PerformanceLoad *interface{} `json:"performance_load,omitempty"`

	RespcodeBrokens *RespcodeBrokens `json:"respcode_brokens,omitempty"`

	RtmpBrokens *RtmpBrokens `json:"rtmp_brokens,omitempty"`

	StreamingErrorBrokens *StreamingErrorBrokens `json:"streaming_error_brokens,omitempty"`

	TpsBrokens *TpsBrokens `json:"tps_brokens,omitempty"`

	VusersBrokens *VusersBrokens `json:"vusers_brokens,omitempty"`
}

func (o BrokenList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BrokenList struct{}"
	}

	return strings.Join([]string{"BrokenList", string(data)}, " ")
}
