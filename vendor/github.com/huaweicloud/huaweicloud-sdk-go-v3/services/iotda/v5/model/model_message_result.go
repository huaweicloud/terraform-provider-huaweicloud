package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 下发的消息响应结果
type MessageResult struct {

	// 消息状态, PENDING，DELIVERED，FAILED和TIMEOUT。如果设备不在线，则平台缓存消息，并且返回PENDING，等设备数据上报之后再下发；如果设备在线，则消息直接进行下发，下发成功后接口返回DELIVERED，失败返回FAILED；如果消息在平台默认时间内（1天）还没有下发给设备，则平台会将消息设置为超时，状态为TIMEOUT。另外应用可以订阅消息的执行结果，平台会将消息结果推送给订阅的应用。
	Status *string `json:"status,omitempty"`

	// 消息的创建时间，\"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串。
	CreatedTime *string `json:"created_time,omitempty"`

	// 消息结束时间, \"yyyyMMdd'T'HHmmss'Z'\"格式的UTC字符串，包含消息转换到DELIVERED，FAILED和TIMEOUT状态的时间。
	FinishedTime *string `json:"finished_time,omitempty"`
}

func (o MessageResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MessageResult struct{}"
	}

	return strings.Join([]string{"MessageResult", string(data)}, " ")
}
