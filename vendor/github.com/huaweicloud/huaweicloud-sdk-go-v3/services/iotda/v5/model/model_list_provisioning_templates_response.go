package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProvisioningTemplatesResponse Response Object
type ListProvisioningTemplatesResponse struct {

	// **参数说明**：预调配模板详情。
	Templates *[]ProvisioningTemplateSimple `json:"templates,omitempty"`

	Page           *Page `json:"page,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ListProvisioningTemplatesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProvisioningTemplatesResponse struct{}"
	}

	return strings.Join([]string{"ListProvisioningTemplatesResponse", string(data)}, " ")
}
