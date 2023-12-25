package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePublishTemplateRequest Request Object
type UpdatePublishTemplateRequest struct {

	// 推流域名
	Domain string `json:"domain"`

	Body *CallbackUrl `json:"body,omitempty"`
}

func (o UpdatePublishTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePublishTemplateRequest struct{}"
	}

	return strings.Join([]string{"UpdatePublishTemplateRequest", string(data)}, " ")
}
