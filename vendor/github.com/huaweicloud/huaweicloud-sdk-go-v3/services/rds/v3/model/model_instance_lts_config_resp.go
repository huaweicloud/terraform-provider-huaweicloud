package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InstanceLtsConfigResp struct {

	// LTS配置信息
	LtsConfigs *[]InstanceLtsConfigDetailResp `json:"lts_configs,omitempty"`

	Instance *InstanceLtsBasicInfoResp `json:"instance,omitempty"`
}

func (o InstanceLtsConfigResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceLtsConfigResp struct{}"
	}

	return strings.Join([]string{"InstanceLtsConfigResp", string(data)}, " ")
}
