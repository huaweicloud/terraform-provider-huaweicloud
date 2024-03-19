package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateAccessCodeRequestBody 生成接入凭证的结构体。
type CreateAccessCodeRequestBody struct {

	// **参数说明**：接入凭证类型，默认为AMQP的接入凭证类型。 **取值范围**： - [AMQP,MQTT]
	Type *string `json:"type,omitempty"`

	// **参数说明**: 是否将AMQP/MQTT连接断开
	ForceDisconnect *bool `json:"force_disconnect,omitempty"`
}

func (o CreateAccessCodeRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAccessCodeRequestBody struct{}"
	}

	return strings.Join([]string{"CreateAccessCodeRequestBody", string(data)}, " ")
}
