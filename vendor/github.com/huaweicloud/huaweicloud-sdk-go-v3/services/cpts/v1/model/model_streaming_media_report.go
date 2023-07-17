package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type StreamingMediaReport struct {

	// 流媒体播放次数(流媒体aw执行次数)
	StreamingPlayTimes *int64 `json:"streamingPlayTimes,omitempty"`

	// 流媒体播放出现失败的次数(失败的流媒体aw次数)
	StreamingErrorTimes *int64 `json:"streamingErrorTimes,omitempty"`

	// 流媒体播放成功率
	StreamingSuccessRate *float64 `json:"streamingSuccessRate,omitempty"`

	// 每秒发送数据包大小
	SentPacketsPerSecond *float64 `json:"sentPacketsPerSecond,omitempty"`

	// 每秒接收数据包大小
	ReceivedPacketsPerSecond *float64 `json:"receivedPacketsPerSecond,omitempty"`

	// 接收数据包大小
	RecvPackets *float64 `json:"recvPackets,omitempty"`

	// 发送数据包大小
	SendPackets *float64 `json:"sendPackets,omitempty"`

	// 音频发送字节大小
	AudioSentBytes *float64 `json:"audioSentBytes,omitempty"`

	// 音频接收字节大小
	AudioRecBytes *float64 `json:"audioRecBytes,omitempty"`

	// 视频发送字节大小
	VideoSentBytes *float64 `json:"videoSentBytes,omitempty"`

	// 视频接收字节大小
	VideoRecBytes *float64 `json:"videoRecBytes,omitempty"`

	// 接收关键帧延迟之和
	SumRecvKeyFrameDelay *float64 `json:"sumRecvKeyFrameDelay,omitempty"`

	// 平均接收关键帧延迟
	AvgRecvKeyFrameDelay *float64 `json:"avgRecvKeyFrameDelay,omitempty"`

	// 最小接收关键帧延迟
	MinRecvKeyFrameDelay *float64 `json:"minRecvKeyFrameDelay,omitempty"`

	// 最大接收关键帧延迟
	MaxRecvKeyFrameDelay *float64 `json:"maxRecvKeyFrameDelay,omitempty"`

	// 发送关键帧延迟之和
	SumSendKeyFrameDelay *float64 `json:"sumSendKeyFrameDelay,omitempty"`

	// 平均发送关键帧延迟
	AvgSendKeyFrameDelay *float64 `json:"avgSendKeyFrameDelay,omitempty"`

	// 最小发送关键帧延迟
	MinSendKeyFrameDelay *float64 `json:"minSendKeyFrameDelay,omitempty"`

	// 最大发送关键帧延迟
	MaxSendKeyFrameDelay *float64 `json:"maxSendKeyFrameDelay,omitempty"`

	// 关键帧发送次数
	KeyFrameSendCnt *float64 `json:"keyFrameSendCnt,omitempty"`

	// 关键帧接收次数
	KeyFrameReceiveCnt *float64 `json:"keyFrameReceiveCnt,omitempty"`

	// TCP连接失败数
	TcpConnectFailed *float64 `json:"tcpConnectFailed,omitempty"`

	// 握手失败数
	HandShakeFailed *float64 `json:"handShakeFailed,omitempty"`

	// RTMP连接失败数
	RtmpConnectFailed *float64 `json:"rtmpConnectFailed,omitempty"`

	// 创建流失败数
	CreateStreamFailed *float64 `json:"createStreamFailed,omitempty"`

	// 播放失败数
	PlayFailed *float64 `json:"playFailed,omitempty"`

	// 发布失败数
	PublishFailed *float64 `json:"publishFailed,omitempty"`

	// 重试失败数
	RetryFailed *float64 `json:"retryFailed,omitempty"`

	// 解析文件失败数
	ParseFileFailed *float64 `json:"parseFileFailed,omitempty"`

	// 非法URL失败数
	IllegalUrlFailed *float64 `json:"illegalUrlFailed,omitempty"`

	// 非法FLV Header失败数
	IllegalFlvHeaderFailed *float64 `json:"illegalFlvHeaderFailed,omitempty"`

	// HTTP连接超时数
	HttpTimeoutFailed *float64 `json:"httpTimeoutFailed,omitempty"`

	// 解析FLV文件失败数
	ParseFlvFileFailed *float64 `json:"parseFlvFileFailed,omitempty"`

	// 未知错误数
	UnknownFailed *float64 `json:"unknownFailed,omitempty"`
}

func (o StreamingMediaReport) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StreamingMediaReport struct{}"
	}

	return strings.Join([]string{"StreamingMediaReport", string(data)}, " ")
}
