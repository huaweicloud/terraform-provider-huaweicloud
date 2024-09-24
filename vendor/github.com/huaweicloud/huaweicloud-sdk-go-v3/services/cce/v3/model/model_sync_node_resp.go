package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SyncNodeResp struct {
}

func (o SyncNodeResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncNodeResp struct{}"
	}

	return strings.Join([]string{"SyncNodeResp", string(data)}, " ")
}
