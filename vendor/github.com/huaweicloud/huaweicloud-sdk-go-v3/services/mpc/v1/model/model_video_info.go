package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VideoInfo struct {

	// 视频宽度
	Width *int32 `json:"width,omitempty"`

	// 视频高度
	Height *int32 `json:"height,omitempty"`

	// 视频码率，单位: kbit/s
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 视频码率，单位: bit/s
	BitrateBps *int64 `json:"bitrate_bps,omitempty"`

	// 帧率。    取值范围：0或[5,60]，0表示自适应。    单位：帧每秒。    > 若设置的帧率不在取值范围内，则自动调整为0，若设置的帧率高于片源帧率，则自动调整为片源帧率。
	FrameRate *int32 `json:"frame_rate,omitempty"`

	// 视频编码格式
	Codec *string `json:"codec,omitempty"`
}

func (o VideoInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VideoInfo struct{}"
	}

	return strings.Join([]string{"VideoInfo", string(data)}, " ")
}
