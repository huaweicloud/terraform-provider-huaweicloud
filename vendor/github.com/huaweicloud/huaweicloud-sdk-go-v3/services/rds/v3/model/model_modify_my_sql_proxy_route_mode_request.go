package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyMySqlProxyRouteModeRequest 修改数据库代理路由模式请求体。
type ModifyMySqlProxyRouteModeRequest struct {

	// 数据库主实例读权重。     - 当route_mode选择0（权重负载）时，该字段取值范围为0~1000。     - 当route_mode选择1或2（负载均衡）时，该字段不生效。
	MasterWeight int32 `json:"master_weight"`

	// 数据库节点的读权重设置。      - 只能为只读实例选择权重。     - 该列表可以为空列表。
	ReadonlyInstances []InstancesWeight `json:"readonly_instances"`

	// 数据库代理路由模式。 取值范围:     0：表示权重负载模式。     1：表示负载均衡模式（数据库主实例不接受读请求）。     2：表示负载均衡模式（数据库主实例接受读请求）。      - 如需使用负载均衡模式，请联系客服申请
	RouteMode int32 `json:"route_mode"`
}

func (o ModifyMySqlProxyRouteModeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyMySqlProxyRouteModeRequest struct{}"
	}

	return strings.Join([]string{"ModifyMySqlProxyRouteModeRequest", string(data)}, " ")
}
