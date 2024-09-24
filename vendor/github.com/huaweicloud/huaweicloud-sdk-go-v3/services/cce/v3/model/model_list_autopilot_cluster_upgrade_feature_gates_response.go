package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotClusterUpgradeFeatureGatesResponse Response Object
type ListAutopilotClusterUpgradeFeatureGatesResponse struct {

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`

	// 特性开关信息,格式为key/value键值对。 - Key: 目前有下列值：DisplayPreCheckDetail(展示所有集群升级前检查项详情),EvsSnapshot(使用EVS快照备份集群), LabelForSkippedNode(支持为集群升级过程中跳过的节点打标签), UpgradeStrategy(集群升级策略) - Value: Support 支持,Disable 关闭,Default 使用CCE服务默认规则判断
	UpgradeFeatureGates map[string]string `json:"upgradeFeatureGates,omitempty"`
	HttpStatusCode      int               `json:"-"`
}

func (o ListAutopilotClusterUpgradeFeatureGatesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotClusterUpgradeFeatureGatesResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotClusterUpgradeFeatureGatesResponse", string(data)}, " ")
}
