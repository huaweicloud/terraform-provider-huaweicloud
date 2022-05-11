package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type ConfigByOption struct {

	// 密码强度策略的正则表达式。(当option为password_regex时返回)
	PasswordRegex *string `json:"password_regex,omitempty"`

	// 密码强度策略的描述。(当option为password_regex_description时返回)
	PasswordRegexDescription *string `json:"password_regex_description,omitempty"`
}

func (o ConfigByOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigByOption struct{}"
	}

	return strings.Join([]string{"ConfigByOption", string(data)}, " ")
}
