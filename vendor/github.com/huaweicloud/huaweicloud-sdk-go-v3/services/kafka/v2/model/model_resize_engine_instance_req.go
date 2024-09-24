package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResizeEngineInstanceReq 实例规格变更请求体。
type ResizeEngineInstanceReq struct {

	// 变更类型。  取值范围：   [- storage：存储空间扩容，代理数量不变。](tag:hws,hws_eu,hws_hk,ocb,hws_ocb,ctc,g42,hk_g42,tm,hk_tm,dt)    - horizontal：代理数量扩容，每个broker的存储空间不变。    [- vertical：垂直扩容，broker的底层虚机规格变更，代理数量和存储空间不变。](tag:hws,hws_eu,hws_hk,ocb,hws_ocb,ctc,g42,hk_g42,tm,hk_tm,dt)
	OperType string `json:"oper_type"`

	// 扩容后的存储空间。  [当oper_type类型是storage或horizontal时，该参数有效且必填。  实例存储空间 = 代理数量 * 每个broker的存储空间。  当oper_type类型是storage时，代理数量不变，每个broker存储空间最少扩容100GB。  当oper_type类型是horizontal时，每个broker的存储空间不变。](tag:hws,hws_eu,hws_hk,ocb,hws_ocb,ctc,g42,hk_g42,tm,hk_tm,dt)  [实例存储空间 = 代理数量 * 每个broker的存储空间。 每个broker的存储空间不变。](tag:hcs,fcs)
	NewStorageSpace *int32 `json:"new_storage_space,omitempty"`

	// 当oper_type参数为horizontal时，该参数有效。  [取值范围：最多支持30个broker。](tag:hws,hws_eu,hws_hk,ocb,hws_ocb,ctc,sbc,hk_sbc,g42,hk_g42,tm,hk_tm)
	NewBrokerNum *int32 `json:"new_broker_num,omitempty"`

	// 垂直扩容时的新产品ID。  当oper_type类型是vertical时，该参数才有效且必填。  产品ID可以从[查询产品规格列表](ListEngineProducts.xml)获取。
	NewProductId *string `json:"new_product_id,omitempty"`

	// 实例绑定的弹性IP地址的ID。  以英文逗号隔开多个弹性IP地址的ID。  当oper_type类型是horizontal时，该参数必填。
	PublicipId *string `json:"publicip_id,omitempty"`

	// 指定的内网IP地址，仅支持指定IPv4。  指定的IP数量只能小于等于新增节点数量。  当指定IP小于节点数量时，未指定的节点随机分配内网IP地址。
	TenantIps *[]string `json:"tenant_ips,omitempty"`

	// 实例扩容时新节点使用备用子网的id  当实例扩容使用备用子网，则传入此值  需要联系客服添加白名单才能传入此值
	SecondTenantSubnetId *string `json:"second_tenant_subnet_id,omitempty"`
}

func (o ResizeEngineInstanceReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeEngineInstanceReq struct{}"
	}

	return strings.Join([]string{"ResizeEngineInstanceReq", string(data)}, " ")
}
