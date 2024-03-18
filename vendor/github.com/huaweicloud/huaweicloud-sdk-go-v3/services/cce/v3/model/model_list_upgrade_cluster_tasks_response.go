package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUpgradeClusterTasksResponse Response Object
type ListUpgradeClusterTasksResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型
	Kind *string `json:"kind,omitempty"`

	Metadata *UpgradeTaskMetadata `json:"metadata,omitempty"`

	// 集群升级任务列表
	Items          *[]UpgradeTaskResponseBody `json:"items,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListUpgradeClusterTasksResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUpgradeClusterTasksResponse struct{}"
	}

	return strings.Join([]string{"ListUpgradeClusterTasksResponse", string(data)}, " ")
}
