package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ShowCeshierarchyRespInstanceIds struct {

	// 实例ID。
	Name *string `json:"name,omitempty"`
}

func (o ShowCeshierarchyRespInstanceIds) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowCeshierarchyRespInstanceIds struct{}"
	}

	return strings.Join([]string{"ShowCeshierarchyRespInstanceIds", string(data)}, " ")
}
