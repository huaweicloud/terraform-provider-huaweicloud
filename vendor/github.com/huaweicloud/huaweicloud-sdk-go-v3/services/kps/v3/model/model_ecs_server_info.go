package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EcsServerInfo 需要绑定密钥对的虚拟机信息。
type EcsServerInfo struct {

	// 需要绑定(替换或重置)SSH密钥对的虚拟机id
	Id string `json:"id"`

	Auth *Auth `json:"auth,omitempty"`

	// - true：禁用虚拟机的ssh登录。 - false：不禁用虚拟机的ssh登录。
	DisablePassword *bool `json:"disable_password,omitempty"`

	// SSH监听端口。
	Port *int64 `json:"port,omitempty"`
}

func (o EcsServerInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EcsServerInfo struct{}"
	}

	return strings.Join([]string{"EcsServerInfo", string(data)}, " ")
}
