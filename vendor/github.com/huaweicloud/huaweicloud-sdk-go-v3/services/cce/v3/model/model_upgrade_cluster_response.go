package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeClusterResponse Response Object
type UpgradeClusterResponse struct {
	Metadata *UpgradeCluserResponseMetadata `json:"metadata,omitempty"`

	Spec           *UpgradeResponseSpec `json:"spec,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o UpgradeClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeClusterResponse struct{}"
	}

	return strings.Join([]string{"UpgradeClusterResponse", string(data)}, " ")
}
