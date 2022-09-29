package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowCeshierarchyRespPartitions struct {

	// 分区名称。
	Name *string `json:"name,omitempty"`
}

func (o ShowCeshierarchyRespPartitions) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCeshierarchyRespPartitions struct{}"
	}

	return strings.Join([]string{"ShowCeshierarchyRespPartitions", string(data)}, " ")
}
