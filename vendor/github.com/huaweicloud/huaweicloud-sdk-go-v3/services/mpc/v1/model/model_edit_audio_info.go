package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type EditAudioInfo struct {

	// 音频编码格式,取值有：[AAC, HEAAC, MP3]。
	Codec *string `json:"codec,omitempty"`

	// 视频码率，单位: bit/s
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 采样率, 单位: HZ
	Sample *int32 `json:"sample,omitempty"`

	// 声道数。
	Channels *string `json:"channels,omitempty"`
}

func (o EditAudioInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EditAudioInfo struct{}"
	}

	return strings.Join([]string{"EditAudioInfo", string(data)}, " ")
}
