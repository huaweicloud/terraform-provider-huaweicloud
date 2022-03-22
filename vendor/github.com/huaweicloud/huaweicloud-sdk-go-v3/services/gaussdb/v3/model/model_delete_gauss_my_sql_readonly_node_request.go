package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteGaussMySqlReadonlyNodeRequest struct {
	// 语言

	XLanguage *string `json:"X-Language,omitempty"`
	// 实例ID，严格匹配UUID规则。

	InstanceId string `json:"instance_id"`
	// 节点ID，严格匹配UUID规则。

	NodeId string `json:"node_id"`
}

func (o DeleteGaussMySqlReadonlyNodeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteGaussMySqlReadonlyNodeRequest struct{}"
	}

	return strings.Join([]string{"DeleteGaussMySqlReadonlyNodeRequest", string(data)}, " ")
}
