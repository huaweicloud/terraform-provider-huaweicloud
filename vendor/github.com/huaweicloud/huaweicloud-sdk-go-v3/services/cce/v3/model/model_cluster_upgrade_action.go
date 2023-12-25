package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ClusterUpgradeAction struct {

	// 插件配置列表
	Addons *[]UpgradeAddonConfig `json:"addons,omitempty"`

	// 节点池内节点升级顺序配置。 > key表示节点池ID，默认节点池取值为\"DefaultPool\"
	NodeOrder map[string][]NodePriority `json:"nodeOrder,omitempty"`

	// 节点池升级顺序配置，key/value对格式。 > key表示节点池ID，默认节点池取值为\"DefaultPool\" > value表示对应节点池的优先级，默认值为0，优先级最低，数值越大优先级越高
	NodePoolOrder map[string]int32 `json:"nodePoolOrder,omitempty"`

	Strategy *UpgradeStrategy `json:"strategy"`

	// 目标集群版本，例如\"v1.23\"
	TargetVersion string `json:"targetVersion"`
}

func (o ClusterUpgradeAction) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ClusterUpgradeAction struct{}"
	}

	return strings.Join([]string{"ClusterUpgradeAction", string(data)}, " ")
}
