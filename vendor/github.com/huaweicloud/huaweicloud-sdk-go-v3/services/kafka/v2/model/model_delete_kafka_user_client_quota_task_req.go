package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteKafkaUserClientQuotaTaskReq 要删除的客户端流控配置
type DeleteKafkaUserClientQuotaTaskReq struct {

	// 用户名
	User *string `json:"user,omitempty"`

	// 客户端ID
	Client *string `json:"client,omitempty"`

	// 是否使用用户默认设置（是则表示对全部用户限流）。
	UserDefault *bool `json:"user-default,omitempty"`

	// 是否使用客户端默认设置（是则表示对全部客户端限流）。
	ClientDefault *bool `json:"client-default,omitempty"`
}

func (o DeleteKafkaUserClientQuotaTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteKafkaUserClientQuotaTaskReq struct{}"
	}

	return strings.Join([]string{"DeleteKafkaUserClientQuotaTaskReq", string(data)}, " ")
}
