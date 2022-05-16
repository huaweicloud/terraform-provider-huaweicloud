package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SSH密钥对私钥托管与保护。
type KeyProtection struct {

	// 导入SSH密钥对的私钥。
	PrivateKey *string `json:"private_key,omitempty"`

	Encryption *Encryption `json:"encryption,omitempty"`
}

func (o KeyProtection) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeyProtection struct{}"
	}

	return strings.Join([]string{"KeyProtection", string(data)}, " ")
}
