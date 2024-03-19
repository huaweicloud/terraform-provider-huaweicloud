package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteKafkaUserClientQuotaTaskReq 待删除的用户/客户端流控配置
type DeleteKafkaUserClientQuotaTaskReq struct {

	// 用户名。  不对全部用户/客户端限流时，用户名和客户端ID不能同时为空。
	User *string `json:"user,omitempty"`

	// 客户端ID。  不对全部用户/客户端限流时，用户名和客户端ID不能同时为空。
	Client *string `json:"client,omitempty"`

	// 是否使用用户默认设置。   - 是，表示对全部用户限流。此时不能同时设置用户名。   - 否，表示对特定用户限流。此时需要设置用户名。
	UserDefault *bool `json:"user-default,omitempty"`

	// 是否使用客户端默认设置。   - 是，表示对全部客户端限流。此时不能设置客户端ID。   - 否，表示对特定客户端限流。此时需要设置客户端ID。
	ClientDefault *bool `json:"client-default,omitempty"`
}

func (o DeleteKafkaUserClientQuotaTaskReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteKafkaUserClientQuotaTaskReq struct{}"
	}

	return strings.Join([]string{"DeleteKafkaUserClientQuotaTaskReq", string(data)}, " ")
}
