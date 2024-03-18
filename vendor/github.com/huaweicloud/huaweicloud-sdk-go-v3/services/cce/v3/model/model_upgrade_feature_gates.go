package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeFeatureGates 集群升级特性开关
type UpgradeFeatureGates struct {

	// 集群升级Console界面是否支持V4版本，该字段一般由CCE Console使用。
	SupportUpgradePageV4 *bool `json:"supportUpgradePageV4,omitempty"`
}

func (o UpgradeFeatureGates) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeFeatureGates struct{}"
	}

	return strings.Join([]string{"UpgradeFeatureGates", string(data)}, " ")
}
