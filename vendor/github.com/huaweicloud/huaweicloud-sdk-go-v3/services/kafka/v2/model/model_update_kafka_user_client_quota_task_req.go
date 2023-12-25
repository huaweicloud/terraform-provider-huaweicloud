package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateKafkaUserClientQuotaTaskReq 客户端流控配置
type UpdateKafkaUserClientQuotaTaskReq struct {

	// 用户名
	User *string `json:"user,omitempty"`

	// 客户端ID
	Client *string `json:"client,omitempty"`

	// 是否使用用户默认设置（是则表示对全部用户限流）。
	UserDefault *bool `json:"user-default,omitempty"`

	// 是否使用客户端默认设置（是则表示对全部客户端限流）。
	ClientDefault *bool `json:"client-default,omitempty"`

	// 生产上限速率（单位为B/s）
	ProducerByteRate *int64 `json:"producer-byte-rate,omitempty"`

	// 消费上限速率（单位为B/s）
	ConsumerByteRate *int64 `json:"consumer-byte-rate,omitempty"`
}

func (o UpdateKafkaUserClientQuotaTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateKafkaUserClientQuotaTaskReq struct{}"
	}

	return strings.Join([]string{"UpdateKafkaUserClientQuotaTaskReq", string(data)}, " ")
}
