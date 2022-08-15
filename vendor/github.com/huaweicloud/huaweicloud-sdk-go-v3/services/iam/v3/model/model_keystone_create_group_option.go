package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneCreateGroupOption struct {

	// 用户组描述信息，长度小于等于255字节。
	Description *string `json:"description,omitempty"`

	// 用户组所属账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId *string `json:"domain_id,omitempty"`

	// 用户组名，长度小于等于128字符。
	Name string `json:"name"`
}

func (o KeystoneCreateGroupOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneCreateGroupOption struct{}"
	}

	return strings.Join([]string{"KeystoneCreateGroupOption", string(data)}, " ")
}
