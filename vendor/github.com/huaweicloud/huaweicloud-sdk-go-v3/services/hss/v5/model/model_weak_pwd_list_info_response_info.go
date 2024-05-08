package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// WeakPwdListInfoResponseInfo 服务器弱口令列表
type WeakPwdListInfoResponseInfo struct {

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP（私有IP），为兼容用户使用，不删除此字段
	HostIp *string `json:"host_ip,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 服务器公网IP
	PublicIp *string `json:"public_ip,omitempty"`

	// 弱口令账号列表
	WeakPwdAccounts *[]WeakPwdAccountInfoResponseInfo `json:"weak_pwd_accounts,omitempty"`
}

func (o WeakPwdListInfoResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "WeakPwdListInfoResponseInfo struct{}"
	}

	return strings.Join([]string{"WeakPwdListInfoResponseInfo", string(data)}, " ")
}
