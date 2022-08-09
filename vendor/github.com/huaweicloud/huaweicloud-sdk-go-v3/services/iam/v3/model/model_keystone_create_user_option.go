package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateUserOption struct {

	// IAM用户名，长度1~64之间，只能包含如下字符：大小写字母、空格、数字或特殊字符（-_.）且不能以数字开头。
	Name string `json:"name"`

	// IAM用户所属账号ID。
	DomainId *string `json:"domain_id,omitempty"`

	// IAM用户密码。 - 系统默认密码最小长度为6位字符，在6-32位之间支持用户自定义密码长度。 - 至少包含以下四种字符中的两种： 大写字母、小写字母、数字和特殊字符。 - 不能包含手机号和邮箱。 - 必须满足账户设置中密码策略的要求。
	Password *string `json:"password,omitempty"`

	// 是否启用IAM用户。true为启用，false为停用，默认为true。
	Enabled *bool `json:"enabled,omitempty"`

	// IAM用户描述信息。
	Description *string `json:"description,omitempty"`
}

func (o KeystoneCreateUserOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateUserOption struct{}"
	}

	return strings.Join([]string{"KeystoneCreateUserOption", string(data)}, " ")
}
