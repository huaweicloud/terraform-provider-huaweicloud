package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SecurityCheckInfoResponseInfo 配置检测结果信息
type SecurityCheckInfoResponseInfo struct {

	// 风险等级，包含如下:   - Low : 低危   - Medium : 中危   - High : 高危
	Severity *string `json:"severity,omitempty"`

	// 配置检查（基线）的名称，例如SSH、CentOS 7、Windows
	CheckName *string `json:"check_name,omitempty"`

	// 配置检查（基线）的类型,Linux系统支持的基线一般check_type和check_name相同,例如SSH、CentOS 7。 Windows系统支持的基线一般check_type和check_name不相同，例如check_name为Windows的配置检查（基线），它的check_type包含Windows Server 2019 R2、Windows Server 2016 R2等。
	CheckType *string `json:"check_type,omitempty"`

	// 标准类型，包含如下:   - cn_standard : 等保合规标准   - hw_standard : 云安全实践标准
	Standard *string `json:"standard,omitempty"`

	// 当前配置检查（基线）类型下，用户共检测了多少个检查项。例如标准类型为hw_standard的SSH基线，主机安全提供了17个检查项，但用户所有主机都只检测了SSH基线的其中5个检查项，check_rule_num就是5。用户有一台主机进行了全量检查项检测，check_rule_num就是17。
	CheckRuleNum *int32 `json:"check_rule_num,omitempty"`

	// 未通过的检查项数量，check_rule_num中只要有一台主机没通过某个检查项，这个检查项就会被计算在failed_rule_num中
	FailedRuleNum *int32 `json:"failed_rule_num,omitempty"`

	// 受影响的服务器的数量，进行了当前基线检测的服务器数量
	HostNum *int32 `json:"host_num,omitempty"`

	// 最新检测时间(ms)
	ScanTime *int64 `json:"scan_time,omitempty"`

	// 对配置检查（基线）类型的描述信息，概括当前基线包含的检查项是根据什么标准制定的，能够审计哪些方面的问题。
	CheckTypeDesc *string `json:"check_type_desc,omitempty"`
}

func (o SecurityCheckInfoResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SecurityCheckInfoResponseInfo struct{}"
	}

	return strings.Join([]string{"SecurityCheckInfoResponseInfo", string(data)}, " ")
}
