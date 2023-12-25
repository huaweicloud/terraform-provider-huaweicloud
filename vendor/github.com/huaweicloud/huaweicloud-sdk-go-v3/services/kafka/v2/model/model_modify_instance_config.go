package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyInstanceConfig struct {

	// 修改的配置名称。
	Name *string `json:"name,omitempty"`

	// 配置的修改值。
	Value *string `json:"value,omitempty"`
}

func (o ModifyInstanceConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyInstanceConfig struct{}"
	}

	return strings.Join([]string{"ModifyInstanceConfig", string(data)}, " ")
}
