package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type PwdPasswordUser struct {
	Domain *PwdPasswordUserDomain `json:"domain"`

	// IAM用户名，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Name string `json:"name"`

	// IAM用户的登录密码。
	Password string `json:"password"`
}

func (o PwdPasswordUser) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PwdPasswordUser struct{}"
	}

	return strings.Join([]string{"PwdPasswordUser", string(data)}, " ")
}
