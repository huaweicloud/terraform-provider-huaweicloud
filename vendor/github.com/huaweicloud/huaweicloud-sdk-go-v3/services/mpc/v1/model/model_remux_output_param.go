package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RemuxOutputParam struct {

	// 输出格式。取值范围： - HLS - MP4
	Format *string `json:"format,omitempty"`

	// 分片时长，仅当“format”为“HLS”时有效。  取值范围：[2，10]。  默认值： 5。  单位：秒。
	SegmentDuration *int32 `json:"segment_duration,omitempty"`

	// 输出媒体是否去除片源的中metadata自定义信息。默认值：false
	RemoveMeta *bool `json:"remove_meta,omitempty"`
}

func (o RemuxOutputParam) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RemuxOutputParam struct{}"
	}

	return strings.Join([]string{"RemuxOutputParam", string(data)}, " ")
}
