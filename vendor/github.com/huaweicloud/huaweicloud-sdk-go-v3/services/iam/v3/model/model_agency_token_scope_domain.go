package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AgencyTokenScopeDomain struct {

	// 委托方A的账号ID，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Id *string `json:"id,omitempty"`

	// 委托方A的账号名，获取方式请参见：[获取账号、IAM用户、项目、用户组、委托的名称和ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	Name *string `json:"name,omitempty"`
}

func (o AgencyTokenScopeDomain) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AgencyTokenScopeDomain struct{}"
	}

	return strings.Join([]string{"AgencyTokenScopeDomain", string(data)}, " ")
}
