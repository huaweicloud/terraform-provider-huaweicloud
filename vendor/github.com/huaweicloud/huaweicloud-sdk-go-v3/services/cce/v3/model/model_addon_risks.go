package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddonRisks 节点风险来源
type AddonRisks struct {

	// 插件模板名称
	AddonTemplateName *string `json:"addonTemplateName,omitempty"`

	// 插件别名
	Alias *string `json:"alias,omitempty"`
}

func (o AddonRisks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddonRisks struct{}"
	}

	return strings.Join([]string{"AddonRisks", string(data)}, " ")
}
