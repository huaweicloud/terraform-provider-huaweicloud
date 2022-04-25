package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateAddressGroupOption struct {

	// 功能说明：地址组名称 取值范围：0-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name *string `json:"name,omitempty"`

	// 功能说明：IP地址组描述信息 取值范围：0-255个字符 约束：不能包含“<”和“>”。
	Description *string `json:"description,omitempty"`

	// 功能说明：IP地址组可包含地址集 取值范围：可以是单个ip地址，ip地址范围，ip地址cidr 约束：当前一个地址组ip_set数量限制默认值为20，即配置的ip地址、ip地址范围或ip地址cidr的总数默认限制20
	IpSet *[]string `json:"ip_set,omitempty"`
}

func (o UpdateAddressGroupOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAddressGroupOption struct{}"
	}

	return strings.Join([]string{"UpdateAddressGroupOption", string(data)}, " ")
}
