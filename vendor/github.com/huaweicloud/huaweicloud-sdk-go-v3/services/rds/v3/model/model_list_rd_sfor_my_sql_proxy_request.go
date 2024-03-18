package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRdSforMySqlProxyRequest Request Object
type ListRdSforMySqlProxyRequest struct {
	ContentType *string `json:"Content-Type,omitempty"`

	// 实例ID，严格匹配UUID规则。
	InstanceId string `json:"instance_id"`

	// 语言。
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ListRdSforMySqlProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRdSforMySqlProxyRequest struct{}"
	}

	return strings.Join([]string{"ListRdSforMySqlProxyRequest", string(data)}, " ")
}
