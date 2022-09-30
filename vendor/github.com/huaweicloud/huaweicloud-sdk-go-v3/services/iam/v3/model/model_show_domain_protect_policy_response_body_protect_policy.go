package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 操作保护策略。
type ShowDomainProtectPolicyResponseBodyProtectPolicy struct {
	AllowUser *AllowUserBody `json:"allow_user"`

	// 是否开启操作保护，取值范围true或false。
	OperationProtection bool `json:"operation_protection"`

	// 操作保护验证指定手机号码。示例：0086-123456789。
	Mobile string `json:"mobile"`

	// 是否指定人员验证。on为指定人员验证，必须填写scene参数。off为操作员验证。
	AdminCheck string `json:"admin_check"`

	// 操作保护验证指定邮件地址。示例：example@email.com。
	Email string `json:"email"`

	// 操作保护指定人员验证方式，admin_check为on时，必须填写。包括mobile、email。
	Scene string `json:"scene"`
}

func (o ShowDomainProtectPolicyResponseBodyProtectPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainProtectPolicyResponseBodyProtectPolicy struct{}"
	}

	return strings.Join([]string{"ShowDomainProtectPolicyResponseBodyProtectPolicy", string(data)}, " ")
}
