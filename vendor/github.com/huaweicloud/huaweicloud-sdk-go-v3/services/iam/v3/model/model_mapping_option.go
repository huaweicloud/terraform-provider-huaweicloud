package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type MappingOption struct {

	// 将联邦用户映射为本地用户的规则列表。
	Rules []MappingRules `json:"rules"`
}

func (o MappingOption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MappingOption struct{}"
	}

	return strings.Join([]string{"MappingOption", string(data)}, " ")
}
