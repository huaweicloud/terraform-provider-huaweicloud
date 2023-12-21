package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowUpgradeClusterTaskResponse Response Object
type ShowUpgradeClusterTaskResponse struct {

	// api版本，默认为v3
	ApiVersion *string `json:"apiVersion,omitempty"`

	// 资源类型，默认为UpgradeTask
	Kind *string `json:"kind,omitempty"`

	Metadata *UpgradeTaskMetadata `json:"metadata,omitempty"`

	Spec *UpgradeTaskSpec `json:"spec,omitempty"`

	Status         *UpgradeTaskStatus `json:"status,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ShowUpgradeClusterTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUpgradeClusterTaskResponse struct{}"
	}

	return strings.Join([]string{"ShowUpgradeClusterTaskResponse", string(data)}, " ")
}
