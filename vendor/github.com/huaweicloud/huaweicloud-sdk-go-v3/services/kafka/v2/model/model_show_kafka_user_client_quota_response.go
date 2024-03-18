package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowKafkaUserClientQuotaResponse Response Object
type ShowKafkaUserClientQuotaResponse struct {

	// 客户端流控配置列表。
	Quotas *[]Quota `json:"quotas,omitempty"`

	// 用户/客户端流控配置数目。
	Count          *int32 `json:"count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowKafkaUserClientQuotaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowKafkaUserClientQuotaResponse struct{}"
	}

	return strings.Join([]string{"ShowKafkaUserClientQuotaResponse", string(data)}, " ")
}
