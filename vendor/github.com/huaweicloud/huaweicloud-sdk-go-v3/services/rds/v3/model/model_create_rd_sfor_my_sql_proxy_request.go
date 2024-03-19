package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateRdSforMySqlProxyRequest Request Object
type CreateRdSforMySqlProxyRequest struct {
	ContentType *string `json:"Content-Type,omitempty"`

	// 实例ID，严格匹配UUID规则。
	InstanceId string `json:"instance_id"`

	// 语言。
	XLanguage *string `json:"X-Language,omitempty"`

	Body *CreateMysqlMultiProxyRequest `json:"body,omitempty"`
}

func (o CreateRdSforMySqlProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRdSforMySqlProxyRequest struct{}"
	}

	return strings.Join([]string{"CreateRdSforMySqlProxyRequest", string(data)}, " ")
}
