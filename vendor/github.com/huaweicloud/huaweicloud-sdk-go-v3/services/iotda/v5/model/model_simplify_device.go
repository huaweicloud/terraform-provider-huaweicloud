package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设备信息结构体，批量查询返回。
type SimplifyDevice struct {

	// 设备ID，用于唯一标识一个设备。在注册设备时直接指定，或者由物联网平台分配获得。由物联网平台分配时，生成规则为\"product_id\" + \"_\" + \"node_id\"拼接而成。
	DeviceId *string `json:"device_id,omitempty"`

	// 设备标识码，通常使用IMEI、MAC地址或Serial No作为nodeId。
	NodeId *string `json:"node_id,omitempty"`

	// 设备名称。
	DeviceName *string `json:"device_name,omitempty"`

	// 设备关联的产品ID，用于唯一标识一个产品模型。
	ProductId *string `json:"product_id,omitempty"`
}

func (o SimplifyDevice) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SimplifyDevice struct{}"
	}

	return strings.Join([]string{"SimplifyDevice", string(data)}, " ")
}
