package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryUpgradeClusterTaskResponse Response Object
type RetryUpgradeClusterTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RetryUpgradeClusterTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryUpgradeClusterTaskResponse struct{}"
	}

	return strings.Join([]string{"RetryUpgradeClusterTaskResponse", string(data)}, " ")
}
