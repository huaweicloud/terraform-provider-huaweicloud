package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExtensionScaleGroupMetadata 扩容伸缩组的基本信息
type ExtensionScaleGroupMetadata struct {

	// 扩展伸缩组的uuid，由系统自动生成
	Uid *string `json:"uid,omitempty"`

	// 扩展伸缩组的名称，不能为 **default**，长度不能超过55个字符，只能包含数字和小写字母以及**-**
	Name *string `json:"name,omitempty"`
}

func (o ExtensionScaleGroupMetadata) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtensionScaleGroupMetadata struct{}"
	}

	return strings.Join([]string{"ExtensionScaleGroupMetadata", string(data)}, " ")
}
