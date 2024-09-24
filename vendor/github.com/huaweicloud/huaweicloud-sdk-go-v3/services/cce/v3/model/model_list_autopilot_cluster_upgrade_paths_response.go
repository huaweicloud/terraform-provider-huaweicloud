package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAutopilotClusterUpgradePathsResponse Response Object
type ListAutopilotClusterUpgradePathsResponse struct {

	// API版本
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`

	// 升级路径集合
	UpgradePaths   *[]UpgradePath `json:"upgradePaths,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListAutopilotClusterUpgradePathsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAutopilotClusterUpgradePathsResponse struct{}"
	}

	return strings.Join([]string{"ListAutopilotClusterUpgradePathsResponse", string(data)}, " ")
}
