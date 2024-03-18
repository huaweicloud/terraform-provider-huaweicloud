package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListClusterUpgradePathsRequest Request Object
type ListClusterUpgradePathsRequest struct {
}

func (o ListClusterUpgradePathsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListClusterUpgradePathsRequest struct{}"
	}

	return strings.Join([]string{"ListClusterUpgradePathsRequest", string(data)}, " ")
}
