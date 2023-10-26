package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchAssociateKeypairRequestBody struct {

	// 最多可同时选择10个弹性云服务器绑定密钥对。  约束：只支持选择相同的密钥对，弹性云服务器处于“运行中”状态，并未绑定密钥对。
	BatchKeypairs []AssociateKeypairRequestBody `json:"batch_keypairs"`
}

func (o BatchAssociateKeypairRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAssociateKeypairRequestBody struct{}"
	}

	return strings.Join([]string{"BatchAssociateKeypairRequestBody", string(data)}, " ")
}
