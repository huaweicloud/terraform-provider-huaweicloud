package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// WeakPwdListInfoResponseInfo 服务器弱口令列表
type WeakPwdListInfoResponseInfo struct {

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP
	HostIp *string `json:"host_ip,omitempty"`

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
