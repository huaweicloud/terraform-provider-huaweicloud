package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ClearPrivateKeyRequest Request Object
type ClearPrivateKeyRequest struct {

	// 密钥对名称。
	KeypairName string `json:"keypair_name"`
}

func (o ClearPrivateKeyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClearPrivateKeyRequest struct{}"
	}

	return strings.Join([]string{"ClearPrivateKeyRequest", string(data)}, " ")
}
