package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 生成接入凭证的结构体。
type CreateAccessCodeRequestBody struct {

	// **参数说明**：接入凭证类型，默认为AMQP的接入凭证类型。 **取值范围**： - AMQP
	Type *string `json:"type,omitempty"`
}

func (o CreateAccessCodeRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAccessCodeRequestBody struct{}"
	}

	return strings.Join([]string{"CreateAccessCodeRequestBody", string(data)}, " ")
}
