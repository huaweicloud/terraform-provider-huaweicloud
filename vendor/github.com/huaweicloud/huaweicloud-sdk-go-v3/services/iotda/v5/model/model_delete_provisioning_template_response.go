package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteProvisioningTemplateResponse Response Object
type DeleteProvisioningTemplateResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteProvisioningTemplateResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteProvisioningTemplateResponse struct{}"
	}

	return strings.Join([]string{"DeleteProvisioningTemplateResponse", string(data)}, " ")
}
