package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResetUserPasswrodReq struct {

	// 用户新密码。  重置用户密码时，为必选参数；  不能与名称或倒序的名称相同。 复杂度要求： - 输入长度为8到32位的字符串。 - 必须包含如下四种字符中的三种组合：   - 小写字母   - 大写字母   - 数字   - 特殊字符包括（`~!@#$%^&*()-_=+\\|[{}]:'\",<.>/?）和空格，并且不能以-开头
	NewPassword *string `json:"new_password,omitempty"`
}

func (o ResetUserPasswrodReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetUserPasswrodReq struct{}"
	}

	return strings.Join([]string{"ResetUserPasswrodReq", string(data)}, " ")
}
