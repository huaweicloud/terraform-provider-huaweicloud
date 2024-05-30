package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateClouddcnSubnetOption
type UpdateClouddcnSubnetOption struct {

	// 功能说明：子网名称 取值范围：1-64，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name *string `json:"name,omitempty"`

	// 功能说明：子网描述 取值范围：0-255个字符，不能包含“<”和“>”。
	Description *string `json:"description,omitempty"`
}

func (o UpdateClouddcnSubnetOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateClouddcnSubnetOption struct{}"
	}

	return strings.Join([]string{"UpdateClouddcnSubnetOption", string(data)}, " ")
}
