package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 模板音频信息
type AudioTemplateInfo struct {

	// 音频采样率(有效值范围) - 1：AUDIO_SAMPLE_AUTO - 2：AUDIO_SAMPLE_22050 - 3：AUDIO_SAMPLE_32000 - 4：AUDIO_SAMPLE_44100 - 5：AUDIO_SAMPLE_48000 - 6：AUDIO_SAMPLE_96000  默认值为1。
	SampleRate int32 `json:"sample_rate"`

	// 音频码率（单位：Kbps）。
	Bitrate *int32 `json:"bitrate,omitempty"`

	// 声道数(有效值范围) - 1：AUDIO_CHANNELS_1 - 2：AUDIO_CHANNELS_2
	Channels int32 `json:"channels"`
}

func (o AudioTemplateInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AudioTemplateInfo struct{}"
	}

	return strings.Join([]string{"AudioTemplateInfo", string(data)}, " ")
}
