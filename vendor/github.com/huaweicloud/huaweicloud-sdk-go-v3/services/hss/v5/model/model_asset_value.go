package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssetValue 资产重要性
type AssetValue struct {
}

func (o AssetValue) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssetValue struct{}"
	}

	return strings.Join([]string{"AssetValue", string(data)}, " ")
}
