package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StreamingErrorBrokens struct {

	// 创建流媒体失败数
	CreateStreamFailed *[]float64 `json:"createStreamFailed,omitempty"`

	// 建立握手失败数
	HandShakeFailed *[]float64 `json:"handShakeFailed,omitempty"`

	// 文件解析失败数
	ParseFileFailed *[]float64 `json:"parseFileFailed,omitempty"`

	// FLV文件解析失败数
	ParseFlvFileFailed *[]float64 `json:"parseFlvFileFailed,omitempty"`

	// 播放失败数
	PlayFailed *[]float64 `json:"playFailed,omitempty"`

	// 发布失败数
	PublishFailed *[]float64 `json:"publishFailed,omitempty"`

	// 重试失败数
	RetryFailed *[]float64 `json:"retryFailed,omitempty"`

	// RTMP连接失败数
	RtmpConnectFailed *[]float64 `json:"rtmpConnectFailed,omitempty"`

	// TCP连接失败数
	TcpConnectFailed *[]float64 `json:"tcpConnectFailed,omitempty"`
}

func (o StreamingErrorBrokens) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StreamingErrorBrokens struct{}"
	}

	return strings.Join([]string{"StreamingErrorBrokens", string(data)}, " ")
}
