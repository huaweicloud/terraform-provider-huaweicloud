package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 修改规则动作请求结构体
type UpdateActionReq struct {

	// **参数说明**：规则动作的类型。 **取值范围**： - HTTP_FORWARDING：HTTP服务消息类型。 - DIS_FORWARDING：转发DIS服务消息类型。 - OBS_FORWARDING：转发OBS服务消息类型。 - AMQP_FORWARDING：转发AMQP服务消息类型。 - DMS_KAFKA_FORWARDING：转发kafka消息类型。
	Channel *string `json:"channel,omitempty"`

	ChannelDetail *ChannelDetail `json:"channel_detail,omitempty"`
}

func (o UpdateActionReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateActionReq struct{}"
	}

	return strings.Join([]string{"UpdateActionReq", string(data)}, " ")
}
