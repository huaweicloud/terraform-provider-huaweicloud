package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClusterUpgradeFeatureGatesRequest Request Object
type ListClusterUpgradeFeatureGatesRequest struct {
}

func (o ListClusterUpgradeFeatureGatesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClusterUpgradeFeatureGatesRequest struct{}"
	}

	return strings.Join([]string{"ListClusterUpgradeFeatureGatesRequest", string(data)}, " ")
}
