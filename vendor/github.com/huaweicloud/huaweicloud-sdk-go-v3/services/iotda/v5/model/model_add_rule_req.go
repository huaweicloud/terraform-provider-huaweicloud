package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 规则触发条件请求结构体
type AddRuleReq struct {

	// **参数说明**：规则名称。 **取值范围**：长度不超过256，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合
	RuleName *string `json:"rule_name,omitempty"`

	// **参数说明**：用户自定义的规则描述。
	Description *string `json:"description,omitempty"`

	Subject *RoutingRuleSubject `json:"subject"`

	// **参数说明**：租户规则的生效范围，默认GLOBAL，。 **取值范围**： - GLOBAL：生效范围为租户级。 - APP：生效范围为资源空间级。如果类型为APP，创建的规则生效范围为携带的app_id指定的资源空间，不携带app_id则创建规则生效范围为[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)。
	AppType *string `json:"app_type,omitempty"`

	// **参数说明**：资源空间ID。。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：用户自定义sql select语句，最大长度500，该参数仅供标准版和企业版用户使用。
	Select *string `json:"select,omitempty"`

	// **参数说明**：用户自定义sql where语句，最大长度500，该参数仅供标准版和企业版用户使用。
	Where *string `json:"where,omitempty"`
}

func (o AddRuleReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddRuleReq struct{}"
	}

	return strings.Join([]string{"AddRuleReq", string(data)}, " ")
}
