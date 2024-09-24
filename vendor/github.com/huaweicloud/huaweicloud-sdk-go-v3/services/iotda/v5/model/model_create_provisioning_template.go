package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateProvisioningTemplate 创建预调配模板请求体。
type CreateProvisioningTemplate struct {

	// **参数说明**：预调配模板名称。 **取值范围**：长度不超过128，只允许中文、字母、数字、下划线（_）、连接符（-）的组合。
	TemplateName string `json:"template_name"`

	// **参数说明**：预调配模板的描述信息。 **取值范围**：长度不超过2048，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合
	Description *string `json:"description,omitempty"`

	TemplateBody *ProvisioningTemplateBody `json:"template_body"`
}

func (o CreateProvisioningTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateProvisioningTemplate struct{}"
	}

	return strings.Join([]string{"CreateProvisioningTemplate", string(data)}, " ")
}
