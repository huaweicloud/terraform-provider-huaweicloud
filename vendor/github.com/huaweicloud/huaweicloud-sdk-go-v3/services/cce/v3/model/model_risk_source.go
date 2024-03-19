package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RiskSource 风险项来源
type RiskSource struct {

	// 配置风险项
	ConfigurationRisks *[]ConfigurationRisks `json:"configurationRisks,omitempty"`

	// 废弃API风险
	DeprecatedAPIRisks *[]DeprecatedApiRisks `json:"deprecatedAPIRisks,omitempty"`

	// 节点风险
	NodeRisks *[]NodeRisks `json:"nodeRisks,omitempty"`

	// 插件风险
	AddonRisks *[]AddonRisks `json:"addonRisks,omitempty"`
}

func (o RiskSource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RiskSource struct{}"
	}

	return strings.Join([]string{"RiskSource", string(data)}, " ")
}
