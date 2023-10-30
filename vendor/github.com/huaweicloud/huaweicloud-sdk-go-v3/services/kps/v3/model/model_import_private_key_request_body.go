package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImportPrivateKeyRequestBody 导入私钥请求体
type ImportPrivateKeyRequestBody struct {
	Keypair *ImportPrivateKeyKeypairBean `json:"keypair"`
}

func (o ImportPrivateKeyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImportPrivateKeyRequestBody struct{}"
	}

	return strings.Join([]string{"ImportPrivateKeyRequestBody", string(data)}, " ")
}
