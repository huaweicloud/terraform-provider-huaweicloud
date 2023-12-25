package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ImportPrivateKeyAction struct {

	// SSH密钥对的名称。 - 新创建的密钥对名称不能和已有密钥对的名称相同。 - SSH密钥对名称由英文字母、数字、下划线、中划线组成,长度不能超过64个字节。
	Name string `json:"name"`

	// SSH密钥对所属的用户信息
	UserId *string `json:"user_id,omitempty"`

	KeyProtection *KeyProtection `json:"key_protection"`
}

func (o ImportPrivateKeyAction) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ImportPrivateKeyAction struct{}"
	}

	return strings.Join([]string{"ImportPrivateKeyAction", string(data)}, " ")
}
