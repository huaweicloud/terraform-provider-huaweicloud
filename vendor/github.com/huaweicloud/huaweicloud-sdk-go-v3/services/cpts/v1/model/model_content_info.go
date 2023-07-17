package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ContentInfo struct {

	// body_type
	BodyType *int32 `json:"body_type,omitempty"`

	// bodys
	Bodys *[]interface{} `json:"bodys,omitempty"`

	// TCP/UDP协议返回数据长度
	CheckEndLength *interface{} `json:"check_end_length,omitempty"`

	// TCP/UDP协议返回结束符
	CheckEndStr *interface{} `json:"check_end_str,omitempty"`

	// TCP/UDP协议返回结束类型，1：返回数据长度；2：结束符
	CheckEndType *interface{} `json:"check_end_type,omitempty"`

	// connect_timeout
	ConnectTimeout *int32 `json:"connect_timeout,omitempty"`

	// connect_type
	ConnectType *int32 `json:"connect_type,omitempty"`

	// headers
	Headers *[]ContentHeader `json:"headers,omitempty"`

	// http_version
	HttpVersion *string `json:"http_version,omitempty"`

	// method
	Method *string `json:"method,omitempty"`

	// name
	Name *string `json:"name,omitempty"`

	// protocol_type
	ProtocolType *int32 `json:"protocol_type,omitempty"`

	// return_timeout
	ReturnTimeout *int32 `json:"return_timeout,omitempty"`

	// return_timeout_param
	ReturnTimeoutParam *string `json:"return_timeout_param,omitempty"`

	// url
	Url *string `json:"url,omitempty"`

	// rtmp地址
	RtmpUrl *string `json:"rtmp_url,omitempty"`

	// flv地址
	FlvUrl *string `json:"flv_url,omitempty"`

	// 分辨率策略
	BitrateType *int32 `json:"bitrate_type,omitempty"`

	// duration
	Duration *int32 `json:"duration,omitempty"`

	// HLS重试延迟时间
	RetryDelay *int32 `json:"retry_delay,omitempty"`

	// HLS重试次数
	RetryTime *int32 `json:"retry_time,omitempty"`
}

func (o ContentInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContentInfo struct{}"
	}

	return strings.Join([]string{"ContentInfo", string(data)}, " ")
}
