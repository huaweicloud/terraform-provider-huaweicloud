package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PauseUpgradeClusterTaskResponse Response Object
type PauseUpgradeClusterTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o PauseUpgradeClusterTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PauseUpgradeClusterTaskResponse struct{}"
	}

	return strings.Join([]string{"PauseUpgradeClusterTaskResponse", string(data)}, " ")
}
