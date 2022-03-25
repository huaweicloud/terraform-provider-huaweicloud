package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 解绑密钥对描述消息体
type DisassociateKeypairRequestBody struct {
	Server *DisassociateEcsServerInfo `json:"server"`
}

func (o DisassociateKeypairRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DisassociateKeypairRequestBody struct{}"
	}

	return strings.Join([]string{"DisassociateKeypairRequestBody", string(data)}, " ")
}
