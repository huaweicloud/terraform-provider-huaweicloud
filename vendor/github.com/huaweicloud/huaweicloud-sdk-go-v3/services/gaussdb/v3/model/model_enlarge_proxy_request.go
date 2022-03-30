package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// proxy节点扩容信息
type EnlargeProxyRequest struct {
	// proxy节点扩容操作需要扩容的节点数。本次扩容的节点数的取值范围：1~30之间的整数。 限制条件：该实例的proxy节点的总数量小于等于32。

	NodeNum int32 `json:"node_num"`
}

func (o EnlargeProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnlargeProxyRequest struct{}"
	}

	return strings.Join([]string{"EnlargeProxyRequest", string(data)}, " ")
}
