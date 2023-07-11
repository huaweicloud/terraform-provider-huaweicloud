package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PwdPolicyInfoResponseInfo 服务器的口令复杂度策略
type PwdPolicyInfoResponseInfo struct {

	// 服务器id(鼠标在“服务器名称”放置后上浮显示)
	HostId *string `json:"host_id,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器IP
	HostIp *string `json:"host_ip,omitempty"`

	// 口令最小长度
	MinLength *bool `json:"min_length,omitempty"`

	// 大写字母
	UppercaseLetter *bool `json:"uppercase_letter,omitempty"`

	// 小写字母
	LowercaseLetter *bool `json:"lowercase_letter,omitempty"`

	// 数字
	Number *bool `json:"number,omitempty"`

	// 特殊字符
	SpecialCharacter *bool `json:"special_character,omitempty"`

	// 修改建议
	Suggestion *string `json:"suggestion,omitempty"`
}

func (o PwdPolicyInfoResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PwdPolicyInfoResponseInfo struct{}"
	}

	return strings.Join([]string{"PwdPolicyInfoResponseInfo", string(data)}, " ")
}
