package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ContentInfo struct {

	// body类型（0：字符串；1：form-data格式；3：x-www-form-urlencoded格式）
	BodyType *int32 `json:"body_type,omitempty"`

	// bodys
	Bodys *[]interface{} `json:"bodys,omitempty"`

	// TCP/UDP协议返回数据长度
	CheckEndLength *interface{} `json:"check_end_length,omitempty"`

	// TCP/UDP协议返回结束符
	CheckEndStr *interface{} `json:"check_end_str,omitempty"`

	// TCP/UDP协议返回结束类型，1：返回数据长度；2：结束符
	CheckEndType *interface{} `json:"check_end_type,omitempty"`

	// 超时时间
	ConnectTimeout *int32 `json:"connect_timeout,omitempty"`

	// 连接设置，当前版本已未使用
	ConnectType *int32 `json:"connect_type,omitempty"`

	// 请求头
	Headers *[]ContentHeader `json:"headers,omitempty"`

	// HTTP版本
	HttpVersion *string `json:"http_version,omitempty"`

	// HTTP方法
	Method *string `json:"method,omitempty"`

	// 用例名称
	Name *string `json:"name,omitempty"`

	// 协议类型（1：HTTP；2：HTTPS；3：TCP；4：UDP；7：HLS/RTMP；9：WebSocket；10：HTTP-FLV）
	ProtocolType *int32 `json:"protocol_type,omitempty"`

	// 响应超时
	ReturnTimeout *int32 `json:"return_timeout,omitempty"`

	// 响应超时参数
	ReturnTimeoutParam *string `json:"return_timeout_param,omitempty"`

	// 请求地址
	Url *string `json:"url,omitempty"`

	// rtmp地址
	RtmpUrl *string `json:"rtmp_url,omitempty"`

	// flv地址
	FlvUrl *string `json:"flv_url,omitempty"`

	// 分辨率策略
	BitrateType *int32 `json:"bitrate_type,omitempty"`

	// 持续时间
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
