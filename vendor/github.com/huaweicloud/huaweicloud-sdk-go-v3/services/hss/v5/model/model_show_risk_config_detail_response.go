package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRiskConfigDetailResponse Response Object
type ShowRiskConfigDetailResponse struct {

	// 风险等级，包含如下:   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 配置检查（基线）的类型,例如SSH、CentOS 7、Windows Server 2019 R2、Windows Server 2016 R2、MySQL5-Windows
	CheckType *string `json:"check_type,omitempty"`

	// 对配置检查（基线）类型的描述信息，概括当前基线包含的检查项是根据什么标准制定的，能够审计哪些方面的问题。
	CheckTypeDesc *string `json:"check_type_desc,omitempty"`

	// 当前配置检查（基线）类型下，用户共检测了多少个检查项。例如标准类型为hw_standard的SSH基线，主机安全提供了17个检查项，但用户所有主机都只检测了SSH基线的其中5个检查项，check_rule_num就是5。用户有一台主机进行了全量检查项检测，check_rule_num就是17。
	CheckRuleNum *int32 `json:"check_rule_num,omitempty"`

	// 未通过的检查项数量，check_rule_num中只要有一台主机没通过某个检查项，这个检查项就会被计算在failed_rule_num中
	FailedRuleNum *int32 `json:"failed_rule_num,omitempty"`

	// 已通过的检查项数量，check_rule_num中只要有一台主机通过了某个检查项，这个检查项就会被计算在passed_rule_num中
	PassedRuleNum *int32 `json:"passed_rule_num,omitempty"`

	// 已忽略的检查项数量，check_rule_num中只要有一台主机忽略了某个检查项，这个检查项就会被计算在ignored_rule_num中
	IgnoredRuleNum *int32 `json:"ignored_rule_num,omitempty"`

	// 受影响的服务器的数量，进行了当前基线检测的服务器数量
	HostNum        *int64 `json:"host_num,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowRiskConfigDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRiskConfigDetailResponse struct{}"
	}

	return strings.Join([]string{"ShowRiskConfigDetailResponse", string(data)}, " ")
}
