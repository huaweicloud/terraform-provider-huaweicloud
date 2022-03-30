package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowGaussMySqlProxyRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`
}

func (o ShowGaussMySqlProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGaussMySqlProxyRequest struct{}"
	}

	return strings.Join([]string{"ShowGaussMySqlProxyRequest", string(data)}, " ")
}
