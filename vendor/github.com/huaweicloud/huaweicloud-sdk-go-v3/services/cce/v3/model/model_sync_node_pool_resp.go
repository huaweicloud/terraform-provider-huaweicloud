package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SyncNodePoolResp struct {
}

func (o SyncNodePoolResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncNodePoolResp struct{}"
	}

	return strings.Join([]string{"SyncNodePoolResp", string(data)}, " ")
}
