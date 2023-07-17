package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BroadcastMessageResponse Response Object
type BroadcastMessageResponse struct {

	// **参数说明**：资源空间ID。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：接收广播消息的完整Topic名称
	TopicFullName *string `json:"topic_full_name,omitempty"`

	// 消息id，由物联网平台生成，用于标识该消息。
	MessageId *string `json:"message_id,omitempty"`

	// 消息的创建时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串。
	CreatedTime    *string `json:"created_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BroadcastMessageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BroadcastMessageResponse struct{}"
	}

	return strings.Join([]string{"BroadcastMessageResponse", string(data)}, " ")
}
