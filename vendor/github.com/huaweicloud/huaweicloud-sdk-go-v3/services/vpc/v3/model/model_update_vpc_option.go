package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateVpcOption struct {

	// 功能说明：虚拟私有云名称 取值范围：1-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点） 约束：与description至少选填一个
	Name *string `json:"name,omitempty"`

	// 功能说明：虚拟私有云描述 取值范围：0-255个字符，不能包含“<”和“>”。 约束：与name至少选填一个
	Description *string `json:"description,omitempty"`
}

func (o UpdateVpcOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateVpcOption struct{}"
	}

	return strings.Join([]string{"UpdateVpcOption", string(data)}, " ")
}
