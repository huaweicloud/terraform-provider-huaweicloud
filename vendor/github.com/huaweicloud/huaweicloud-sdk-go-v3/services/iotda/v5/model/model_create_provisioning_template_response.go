package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateProvisioningTemplateResponse Response Object
type CreateProvisioningTemplateResponse struct {

	// **参数说明**：预调配模板ID。
	TemplateId *string `json:"template_id,omitempty"`

	// **参数说明**：预调配模板名称。 **取值范围**：长度不超过128，只允许中文、字母、数字、下划线（_）、连接符（-）的组合。
	TemplateName *string `json:"template_name,omitempty"`

	// **参数说明**：预调配模板的描述信息。 **取值范围**：长度不超过2048，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合
	Description *string `json:"description,omitempty"`

	TemplateBody *ProvisioningTemplateBody `json:"template_body,omitempty"`

	// 在物联网平台创建预调配模板的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 在物联网平台更新预调配模板的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	UpdateTime     *string `json:"update_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateProvisioningTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateProvisioningTemplateResponse struct{}"
	}

	return strings.Join([]string{"CreateProvisioningTemplateResponse", string(data)}, " ")
}
