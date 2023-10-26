package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ImportPrivateKeyProtection SSH密钥对私钥托管与保护。
type ImportPrivateKeyProtection struct {

	// 导入SSH密钥对的私钥。
	PrivateKey string `json:"private_key"`

	Encryption *Encryption `json:"encryption"`
}

func (o ImportPrivateKeyProtection) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImportPrivateKeyProtection struct{}"
	}

	return strings.Join([]string{"ImportPrivateKeyProtection", string(data)}, " ")
}
