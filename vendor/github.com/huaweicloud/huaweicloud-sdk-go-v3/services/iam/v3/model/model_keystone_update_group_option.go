package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type KeystoneUpdateGroupOption struct {

	// 用户组描述信息，长度小于等于255字节。name与description至少填写一个。
	Description *string `json:"description,omitempty"`

	// 用户组所属账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId *string `json:"domain_id,omitempty"`

	// 用户组名，长度小于等于128字符。name与description至少填写一个。
	Name *string `json:"name,omitempty"`
}

func (o KeystoneUpdateGroupOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneUpdateGroupOption struct{}"
	}

	return strings.Join([]string{"KeystoneUpdateGroupOption", string(data)}, " ")
}
