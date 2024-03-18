package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRdSforMySqlProxyResponse Response Object
type ListRdSforMySqlProxyResponse struct {

	// 数据库实例下的数据库代理信息列表。
	ProxyQueryInfoList *[]QueryProxyResponseV3 `json:"proxy_query_info_list,omitempty"`

	// 支持同时开启的数据库代理的最大数量。
	MaxProxyNum *int32 `json:"max_proxy_num,omitempty"`

	// 单个数据库代理支持选择的代理节点的最大数量。
	MaxProxyNodeNum *int32 `json:"max_proxy_node_num,omitempty"`

	// 是否支持创建数据库代理时设置负载均衡路由模式。
	SupportBalanceRouteModeForFavoredVersion *bool `json:"support_balance_route_mode_for_favored_version,omitempty"`
	HttpStatusCode                           int   `json:"-"`
}

func (o ListRdSforMySqlProxyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRdSforMySqlProxyResponse struct{}"
	}

	return strings.Join([]string{"ListRdSforMySqlProxyResponse", string(data)}, " ")
}
