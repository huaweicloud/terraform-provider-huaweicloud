package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateKafkaUserClientQuotaTaskReq 用户/客户端流控配置
type CreateKafkaUserClientQuotaTaskReq struct {

	// 用户名。  不对全部用户/客户端限流时，用户名和客户端ID不能同时为空。
	User *string `json:"user,omitempty"`

	// 客户端ID。  不对全部用户/客户端限流时，用户名和客户端ID不能同时为空。
	Client *string `json:"client,omitempty"`

	// 是否使用用户默认设置。   - 是，表示对全部用户限流。此时不能同时设置用户名。   - 否，表示对特定用户限流。此时需要设置用户名。
	UserDefault *bool `json:"user-default,omitempty"`

	// 是否使用客户端默认设置。   - 是，表示对全部客户端限流。此时不能设置客户端ID。   - 否，表示对特定客户端限流。此时需要设置客户端ID。
	ClientDefault *bool `json:"client-default,omitempty"`

	// 生产上限速率（单位为B/s）。
	ProducerByteRate *int64 `json:"producer-byte-rate,omitempty"`

	// 消费上限速率（单位为B/s）。  > “生产上限速率”和“消费上限速率”不可同时为空。
	ConsumerByteRate *int64 `json:"consumer-byte-rate,omitempty"`
}

func (o CreateKafkaUserClientQuotaTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKafkaUserClientQuotaTaskReq struct{}"
	}

	return strings.Join([]string{"CreateKafkaUserClientQuotaTaskReq", string(data)}, " ")
}
