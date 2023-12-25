package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeStrategy 升级配置
type UpgradeStrategy struct {

	// 升级策略类型，当前仅支持原地升级类型\"inPlaceRollingUpdate\"
	Type string `json:"type"`

	InPlaceRollingUpdate *InPlaceRollingUpdate `json:"inPlaceRollingUpdate,omitempty"`
}

func (o UpgradeStrategy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeStrategy struct{}"
	}

	return strings.Join([]string{"UpgradeStrategy", string(data)}, " ")
}
