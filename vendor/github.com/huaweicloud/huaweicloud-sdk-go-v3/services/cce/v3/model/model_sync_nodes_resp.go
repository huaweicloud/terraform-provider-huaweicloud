package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SyncNodesResp struct {
}

func (o SyncNodesResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SyncNodesResp struct{}"
	}

	return strings.Join([]string{"SyncNodesResp", string(data)}, " ")
}
