package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SkippedCheckItemList struct {

	// 跳过的检查项名称
	Name *string `json:"name,omitempty"`

	ResourceSelector *ResourceSelector `json:"resourceSelector,omitempty"`
}

func (o SkippedCheckItemList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SkippedCheckItemList struct{}"
	}

	return strings.Join([]string{"SkippedCheckItemList", string(data)}, " ")
}
