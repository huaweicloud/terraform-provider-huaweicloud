package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ContinueUpgradeClusterTaskResponse Response Object
type ContinueUpgradeClusterTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ContinueUpgradeClusterTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ContinueUpgradeClusterTaskResponse struct{}"
	}

	return strings.Join([]string{"ContinueUpgradeClusterTaskResponse", string(data)}, " ")
}
