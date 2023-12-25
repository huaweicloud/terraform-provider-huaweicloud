package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InPlaceRollingUpdate 原地升级配置
type InPlaceRollingUpdate struct {

	// 节点升级步长，取值范围为[1, 40]，建议取值20
	UserDefinedStep *int32 `json:"userDefinedStep,omitempty"`
}

func (o InPlaceRollingUpdate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InPlaceRollingUpdate struct{}"
	}

	return strings.Join([]string{"InPlaceRollingUpdate", string(data)}, " ")
}
