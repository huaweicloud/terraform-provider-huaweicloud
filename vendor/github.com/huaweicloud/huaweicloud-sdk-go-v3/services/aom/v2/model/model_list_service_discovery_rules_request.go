package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListServiceDiscoveryRulesRequest struct {

	// 具体的服务发现规则ID,可以精确匹配到一条服务发现规则。不传时返回project下所有服务发现规则的列表。
	Id *string `json:"id,omitempty"`
}

func (o ListServiceDiscoveryRulesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListServiceDiscoveryRulesRequest struct{}"
	}

	return strings.Join([]string{"ListServiceDiscoveryRulesRequest", string(data)}, " ")
}
