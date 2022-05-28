package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建规则动作请求结构体
type AddActionReq struct {

	// **参数说明**：规则触发条件ID，用于唯一标识一条规则触发条件，在创建规则时由物联网平台分配获得。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	RuleId string `json:"rule_id"`

	// **参数说明**：规则动作的类型。 **取值范围**： - HTTP_FORWARDING：HTTP服务消息类型。 - DIS_FORWARDING：转发DIS服务消息类型。 - OBS_FORWARDING：转发OBS服务消息类型。 - AMQP_FORWARDING：转发AMQP服务消息类型。 - DMS_KAFKA_FORWARDING：转发kafka消息类型。
	Channel string `json:"channel"`

	ChannelDetail *ChannelDetail `json:"channel_detail"`
}

func (o AddActionReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddActionReq struct{}"
	}

	return strings.Join([]string{"AddActionReq", string(data)}, " ")
}
