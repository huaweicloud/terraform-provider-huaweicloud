package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SmnInfo SMN信息
type SmnInfo struct {

	// smn发送消息的内容
	SmnNotifiedContent *string `json:"smn_notified_content,omitempty"`

	// smn的订阅的状态
	SmnSubscriptionStatus *int32 `json:"smn_subscription_status,omitempty"`

	// smn的订阅类型
	SmnSubscriptionType *string `json:"smn_subscription_type,omitempty"`
}

func (o SmnInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmnInfo struct{}"
	}

	return strings.Join([]string{"SmnInfo", string(data)}, " ")
}
