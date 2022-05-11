package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type SecurityCompliance struct {

	// 密码强度策略的正则表达式。
	PasswordRegex string `json:"password_regex"`

	// 密码强度策略的描述。
	PasswordRegexDescription string `json:"password_regex_description"`
}

func (o SecurityCompliance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SecurityCompliance struct{}"
	}

	return strings.Join([]string{"SecurityCompliance", string(data)}, " ")
}
