package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryAutopilotUpgradeClusterTaskResponse Response Object
type RetryAutopilotUpgradeClusterTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RetryAutopilotUpgradeClusterTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryAutopilotUpgradeClusterTaskResponse struct{}"
	}

	return strings.Join([]string{"RetryAutopilotUpgradeClusterTaskResponse", string(data)}, " ")
}
