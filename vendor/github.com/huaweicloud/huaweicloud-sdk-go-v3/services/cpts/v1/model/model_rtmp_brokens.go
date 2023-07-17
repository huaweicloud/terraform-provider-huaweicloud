package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RtmpBrokens struct {

	// 音频接收字节数
	AudioRecBytes *[]float64 `json:"audioRecBytes,omitempty"`

	// 音频发送字节数
	AudioSentBytes *[]float64 `json:"audioSentBytes,omitempty"`

	// RTMP接收数据包数
	RtmpReceivedPackets *[]float64 `json:"rtmpReceivedPackets,omitempty"`

	// RTMP发送数据包数
	RtmpSentPackets *[]float64 `json:"rtmpSentPackets,omitempty"`

	// 视频接收字节数
	VideoRecBytes *[]float64 `json:"videoRecBytes,omitempty"`

	// 视频发送字节数
	VideoSentBytes *[]float64 `json:"videoSentBytes,omitempty"`
}

func (o RtmpBrokens) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RtmpBrokens struct{}"
	}

	return strings.Join([]string{"RtmpBrokens", string(data)}, " ")
}
