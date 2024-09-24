package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProvisioningTemplate 更新预调配模板请求体。
type UpdateProvisioningTemplate struct {

	// **参数说明**：预调配模板的描述信息。 **取值范围**：长度不超过2048，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合
	Description *string `json:"description,omitempty"`

	TemplateBody *ProvisioningTemplateBody `json:"template_body,omitempty"`
}

func (o UpdateProvisioningTemplate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProvisioningTemplate struct{}"
	}

	return strings.Join([]string{"UpdateProvisioningTemplate", string(data)}, " ")
}
