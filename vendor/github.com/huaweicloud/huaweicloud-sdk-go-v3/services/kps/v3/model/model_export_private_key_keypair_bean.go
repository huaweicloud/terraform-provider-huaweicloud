package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ExportPrivateKeyKeypairBean struct {

	// SSH密钥对的名称。
	Name string `json:"name"`

	// SSH密钥对的私钥
	PrivateKey string `json:"private_key"`
}

func (o ExportPrivateKeyKeypairBean) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExportPrivateKeyKeypairBean struct{}"
	}

	return strings.Join([]string{"ExportPrivateKeyKeypairBean", string(data)}, " ")
}
