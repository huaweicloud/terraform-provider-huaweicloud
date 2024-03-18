package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTenantAccessInfoRequest Request Object
type UpdateTenantAccessInfoRequest struct {

	// 客户端语言
	XLanguage *string `json:"x-language,omitempty"`

	Body *UpdateTenantAccessInfoReq `json:"body,omitempty"`
}

func (o UpdateTenantAccessInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTenantAccessInfoRequest struct{}"
	}

	return strings.Join([]string{"UpdateTenantAccessInfoRequest", string(data)}, " ")
}
