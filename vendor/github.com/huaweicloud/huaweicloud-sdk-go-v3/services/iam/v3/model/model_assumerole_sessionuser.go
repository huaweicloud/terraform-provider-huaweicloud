package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AssumeroleSessionuser struct {

	// 委托方对应的企业用户名。用户名需满足如下规则：长度5~64，只能包含大写字母、小写字母、数字（0-9）、特殊字符（\"-\"与\"_\"）且只能以字母开头。
	Name *string `json:"name,omitempty"`
}

func (o AssumeroleSessionuser) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssumeroleSessionuser struct{}"
	}

	return strings.Join([]string{"AssumeroleSessionuser", string(data)}, " ")
}
