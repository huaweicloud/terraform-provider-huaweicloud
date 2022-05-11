package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 需要绑定密钥对的虚拟机信息。
type EcsServerInfo struct {

	// 需要绑定(替换或重置)SSH密钥对的虚拟机id
	Id string `json:"id"`

	Auth *Auth `json:"auth,omitempty"`

	// - true：禁用虚拟机的ssh登陆。 - false：不禁用虚拟机的ssh登陆。
	DisablePassword *bool `json:"disable_password,omitempty"`
}

func (o EcsServerInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EcsServerInfo struct{}"
	}

	return strings.Join([]string{"EcsServerInfo", string(data)}, " ")
}
