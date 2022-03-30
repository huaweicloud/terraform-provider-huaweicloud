package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type CreateGaussMySqlProxyRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID

	InstanceId string `json:"instance_id"`

	Body *OpenMysqlProxyRequestBody `json:"body,omitempty"`
}

func (o CreateGaussMySqlProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateGaussMySqlProxyRequest struct{}"
	}

	return strings.Join([]string{"CreateGaussMySqlProxyRequest", string(data)}, " ")
}
