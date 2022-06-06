package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateDomainProtectPolicyResponseBodyProtectPolicy struct {
	AllowUser *AllowUserBody `json:"allow_user"`

	// 是否开启操作保护，取值范围true或false。
	OperationProtection bool `json:"operation_protection"`

	// 是否指定人员验证。on为指定人员验证，必须填写scene参数。off为操作员验证。
	AdminCheck string `json:"admin_check"`

	// 操作保护指定人员验证方式，admin_check为on时，必须填写。包括mobile、email。
	Scene string `json:"scene"`
}

func (o UpdateDomainProtectPolicyResponseBodyProtectPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDomainProtectPolicyResponseBodyProtectPolicy struct{}"
	}

	return strings.Join([]string{"UpdateDomainProtectPolicyResponseBodyProtectPolicy", string(data)}, " ")
}
