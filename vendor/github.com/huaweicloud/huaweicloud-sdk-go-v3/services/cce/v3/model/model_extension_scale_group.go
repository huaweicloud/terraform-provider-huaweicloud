package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ExtensionScaleGroup 节点池扩展伸缩组配置
type ExtensionScaleGroup struct {
	Metadata *ExtensionScaleGroupMetadata `json:"metadata,omitempty"`

	Spec *ExtensionScaleGroupSpec `json:"spec,omitempty"`
}

func (o ExtensionScaleGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtensionScaleGroup struct{}"
	}

	return strings.Join([]string{"ExtensionScaleGroup", string(data)}, " ")
}
