package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SmnResponse SMN返回的发送结果信息
type SmnResponse struct {

	// 发送时间。
	SentTime *int64 `json:"sent_time,omitempty"`

	// 发送的通知的消息内容。
	SmnNotifiedHistory *[]SmnInfo `json:"smn_notified_history,omitempty"`

	// 请求smn服务的请求id。
	SmnRequestId *string `json:"smn_request_id,omitempty"`

	// 调用smn服务返回的信息。
	SmnResponseBody *string `json:"smn_response_body,omitempty"`

	// 调用smn服务返回的http状态码。
	SmnResponseCode *string `json:"smn_response_code,omitempty"`

	// smn的主题。
	SmnTopic *string `json:"smn_topic,omitempty"`
}

func (o SmnResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmnResponse struct{}"
	}

	return strings.Join([]string{"SmnResponse", string(data)}, " ")
}
