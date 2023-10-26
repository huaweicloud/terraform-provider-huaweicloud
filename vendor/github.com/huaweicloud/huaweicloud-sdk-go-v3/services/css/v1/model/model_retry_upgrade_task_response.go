package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RetryUpgradeTaskResponse Response Object
type RetryUpgradeTaskResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o RetryUpgradeTaskResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RetryUpgradeTaskResponse struct{}"
	}

	return strings.Join([]string{"RetryUpgradeTaskResponse", string(data)}, " ")
}
