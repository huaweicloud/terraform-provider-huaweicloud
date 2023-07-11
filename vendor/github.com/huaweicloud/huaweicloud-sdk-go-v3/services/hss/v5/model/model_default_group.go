package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DefaultGroup 是否是默认策略组
type DefaultGroup struct {
}

func (o DefaultGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DefaultGroup struct{}"
	}

	return strings.Join([]string{"DefaultGroup", string(data)}, " ")
}
