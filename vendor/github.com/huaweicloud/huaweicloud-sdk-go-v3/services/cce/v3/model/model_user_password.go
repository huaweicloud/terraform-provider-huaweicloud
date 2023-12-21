package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UserPassword struct {

	// 登录帐号，默认为“root”
	Username *string `json:"username,omitempty"`

	// 登录密码，若创建节点通过用户名密码方式，即使用该字段，则响应体中该字段作屏蔽展示。 密码复杂度要求： - 长度为8-26位。 - 密码至少必须包含大写字母、小写字母、数字和特殊字符（!@$%^-_=+[{}]:,./?~#*）中的三种。 - 密码不能包含用户名或用户名的逆序。 创建节点时password字段需要加盐加密，具体方法请参见[创建节点时password字段加盐加密](add-salt.xml)。
	Password string `json:"password"`
}

func (o UserPassword) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserPassword struct{}"
	}

	return strings.Join([]string{"UserPassword", string(data)}, " ")
}
