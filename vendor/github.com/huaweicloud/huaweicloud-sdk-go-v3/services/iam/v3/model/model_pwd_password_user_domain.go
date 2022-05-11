package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type PwdPasswordUserDomain struct {

	// IAM用户所属账号名，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Name string `json:"name"`
}

func (o PwdPasswordUserDomain) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PwdPasswordUserDomain struct{}"
	}

	return strings.Join([]string{"PwdPasswordUserDomain", string(data)}, " ")
}
