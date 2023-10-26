package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExportPrivateKeyRequestBody 导出私钥请求体
type ExportPrivateKeyRequestBody struct {
	Keypair *KeypairBean `json:"keypair"`
}

func (o ExportPrivateKeyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExportPrivateKeyRequestBody struct{}"
	}

	return strings.Join([]string{"ExportPrivateKeyRequestBody", string(data)}, " ")
}
