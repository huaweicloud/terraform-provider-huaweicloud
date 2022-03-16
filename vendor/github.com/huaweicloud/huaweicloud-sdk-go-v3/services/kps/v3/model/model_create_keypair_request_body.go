package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建密钥对请求体
type CreateKeypairRequestBody struct {
	Keypair *CreateKeypairAction `json:"keypair"`
}

func (o CreateKeypairRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateKeypairRequestBody struct{}"
	}

	return strings.Join([]string{"CreateKeypairRequestBody", string(data)}, " ")
}
