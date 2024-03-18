package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTenantAccessInfoRequest Request Object
type ShowTenantAccessInfoRequest struct {

	// 客户端语言
	XLanguage *string `json:"x-language,omitempty"`
}

func (o ShowTenantAccessInfoRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTenantAccessInfoRequest struct{}"
	}

	return strings.Join([]string{"ShowTenantAccessInfoRequest", string(data)}, " ")
}
