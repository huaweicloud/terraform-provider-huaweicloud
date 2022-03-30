package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OpenMysqlProxyRequestBody struct {
	// 代理规格码。

	FlavorRef *string `json:"flavor_ref,omitempty"`
	// 代理实例节点数，取值整数2-32。

	NodeNum *int32 `json:"node_num,omitempty"`
}

func (o OpenMysqlProxyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OpenMysqlProxyRequestBody struct{}"
	}

	return strings.Join([]string{"OpenMysqlProxyRequestBody", string(data)}, " ")
}
