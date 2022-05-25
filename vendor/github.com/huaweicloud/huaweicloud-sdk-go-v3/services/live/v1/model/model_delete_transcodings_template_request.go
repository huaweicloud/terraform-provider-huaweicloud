package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteTranscodingsTemplateRequest struct {

	// 推流域名
	Domain string `json:"domain"`

	// 应用名称
	AppName string `json:"app_name"`
}

func (o DeleteTranscodingsTemplateRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodingsTemplateRequest struct{}"
	}

	return strings.Join([]string{"DeleteTranscodingsTemplateRequest", string(data)}, " ")
}
