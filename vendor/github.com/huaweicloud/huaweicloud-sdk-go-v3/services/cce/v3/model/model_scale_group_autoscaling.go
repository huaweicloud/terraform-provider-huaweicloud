package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleGroupAutoscaling 扩展伸缩组弹性伸缩配置
type ScaleGroupAutoscaling struct {

	// 伸缩组弹性扩缩容启用开关，默认不开启
	Enable *bool `json:"enable,omitempty"`

	// 伸缩组优先级，未设置则默认为0，数值越大优先级越高
	ExtensionPriority *int32 `json:"extensionPriority,omitempty"`

	// 弹性伸缩时，伸缩组最少应保持的节点数量，必须大于0
	MinNodeCount *int32 `json:"minNodeCount,omitempty"`

	// 弹性伸缩时，伸缩组最多可保持的节点数量，应大于等于 **minNodeCount**, 不可大于集群规格所允许的节点上限，不可大于节点池节点数量上限
	MaxNodeCount *int32 `json:"maxNodeCount,omitempty"`
}

func (o ScaleGroupAutoscaling) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleGroupAutoscaling struct{}"
	}

	return strings.Join([]string{"ScaleGroupAutoscaling", string(data)}, " ")
}
