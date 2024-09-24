package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeAutopilotClusterResponse Response Object
type UpgradeAutopilotClusterResponse struct {
	Metadata *UpgradeCluserResponseMetadata `json:"metadata,omitempty"`

	Spec           *UpgradeResponseSpec `json:"spec,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o UpgradeAutopilotClusterResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeAutopilotClusterResponse struct{}"
	}

	return strings.Join([]string{"UpgradeAutopilotClusterResponse", string(data)}, " ")
}
