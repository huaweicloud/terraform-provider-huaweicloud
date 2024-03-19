package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateMysqlMultiProxyRequest 开启数据库代理请求体。
type CreateMysqlMultiProxyRequest struct {

	// 数据库代理规格码。      - 当局点支持主备模式数据库代理时，该字段不生效。     - 当局点支持集群模式数据库代理时，该字段请参考查询数据库代理规格信息接口返回体中，[规格信息]中的code字段。
	FlavorRef string `json:"flavor_ref"`

	// 数据库代理节点数量。      - 当局点支持主备模式数据库代理时，请设置该字段为固定值2。     - 当局点支持集群模式数据库代理时，该字段最小值为2，最大值请参考查询数据库代理信息列表接口返回体中，[数据库代理信息列表]中的max_proxy_node_num字段值。
	NodeNum int32 `json:"node_num"`

	// 数据库代理名称。用于表示实例的名称，同一租户下，同类型的实例名可重名。  取值范围：最小长度为4个字符，最大不超过64个字节，以字母或中文字符开头，只能包含字母、数字、中划线、下划线、英文句号和中文。  当不选择该参数或局点仅支持主备模式数据库代理时，将随机生成名称。
	ProxyName *string `json:"proxy_name,omitempty"`

	// 数据库代理读写模式。 取值范围:     readwrite：读写模式。     readonly：只读模式。
	ProxyMode *string `json:"proxy_mode,omitempty"`

	// 数据库代理路由模式。 取值范围:     0：表示权重负载模式。     1：表示负载均衡模式（数据库主节点不接受读请求）。     2：表示负载均衡模式（数据库主节点接受读请求）。      - 如需使用负载均衡模式，请联系客服申请
	RouteMode *int32 `json:"route_mode,omitempty"`

	// 数据库节点的读权重设置。      - 在proxy_mode（数据库代理读写模式）为readonly（只读模式）或者在route_mode（路由模式）>0时，只能为只读节点选择权重。     - 在proxy_mode（数据库代理读写模式）为readonly（只读模式）时，需要至少为一个只读实例配置权重。     - 在route_mode（路由模式）>0时，为主实例配置的权重将不生效。     - 该列表可以为空列表。
	NodesReadWeight []InstancesWeight `json:"nodes_read_weight"`

	// 数据库VPC下的子网ID。 取值范围为该实例所在VPC下的所有子网ID。      - 如需使用该参数，请联系客服申请。     - 获取子网ID请参考[创建VPC和子网](https://support.huaweicloud.com/api-cce/cce_02_0100.html)
	SubnetId *string `json:"subnet_id,omitempty"`
}

func (o CreateMysqlMultiProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateMysqlMultiProxyRequest struct{}"
	}

	return strings.Join([]string{"CreateMysqlMultiProxyRequest", string(data)}, " ")
}
