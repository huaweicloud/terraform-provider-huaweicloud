package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowCeshierarchyRespNodes struct {

	// 节点名称。
	Name *string `json:"name,omitempty"`
}

func (o ShowCeshierarchyRespNodes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCeshierarchyRespNodes struct{}"
	}

	return strings.Join([]string{"ShowCeshierarchyRespNodes", string(data)}, " ")
}
