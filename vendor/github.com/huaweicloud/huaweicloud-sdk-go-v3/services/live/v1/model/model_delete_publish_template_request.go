package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePublishTemplateRequest Request Object
type DeletePublishTemplateRequest struct {

	// 推流域名
	Domain string `json:"domain"`
}

func (o DeletePublishTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePublishTemplateRequest struct{}"
	}

	return strings.Join([]string{"DeletePublishTemplateRequest", string(data)}, " ")
}
