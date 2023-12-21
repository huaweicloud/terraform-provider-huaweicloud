package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPublishTemplateRequest Request Object
type ListPublishTemplateRequest struct {

	// 推流域名
	Domain string `json:"domain"`
}

func (o ListPublishTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPublishTemplateRequest struct{}"
	}

	return strings.Join([]string{"ListPublishTemplateRequest", string(data)}, " ")
}
