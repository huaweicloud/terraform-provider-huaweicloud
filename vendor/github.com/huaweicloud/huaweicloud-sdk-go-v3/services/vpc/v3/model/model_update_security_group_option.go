package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateSecurityGroupOption struct {

	// 功能说明：安全组名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name *string `json:"name,omitempty"`

	// 功能说明：安全组描述 取值范围：0-255个字符，不能包含“<”和“>”
	Description *string `json:"description,omitempty"`
}

func (o UpdateSecurityGroupOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSecurityGroupOption struct{}"
	}

	return strings.Join([]string{"UpdateSecurityGroupOption", string(data)}, " ")
}
