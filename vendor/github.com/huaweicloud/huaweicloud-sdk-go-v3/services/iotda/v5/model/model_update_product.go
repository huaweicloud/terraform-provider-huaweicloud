package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProduct 修改产品信息结构体。
type UpdateProduct struct {

	// **参数说明**：资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数，指定要修改的产品属于哪个资源空间；若不携带，则优先修改默认资源空间下产品，如默认资源空间下无对应产品，则按照产品创建时间修改最早创建产品。如果用户存在多资源空间，同时又不想携带该参数，可以联系华为技术支持对用户数据做资源空间合并。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：产品名称。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Name *string `json:"name,omitempty"`

	// **参数说明**：设备类型。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	DeviceType *string `json:"device_type,omitempty"`

	// **参数说明**：设备使用的协议类型。注：禁止其他协议类型修改为CoAP。 **取值范围**：MQTT，CoAP，HTTP，HTTPS，Modbus，ONVIF，OPC-UA，OPC-DA，Other，TCP，UDP。
	ProtocolType *string `json:"protocol_type,omitempty"`

	// **参数说明**：设备上报数据的格式。 **取值范围**： - json：JSON格式 - binary：二进制码流格式
	DataFormat *string `json:"data_format,omitempty"`

	// **参数说明**：设备的服务能力列表。
	ServiceCapabilities *[]ServiceCapability `json:"service_capabilities,omitempty"`

	// **参数说明**：厂商名称。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	ManufacturerName *string `json:"manufacturer_name,omitempty"`

	// **参数说明**：设备所属行业。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Industry *string `json:"industry,omitempty"`

	// **参数说明**：产品的描述信息。 **取值范围**：长度不超过128，只允许中文、字母、数字、空白字符、以及_?'#().,;&%@!- ，、：；。/等字符的组合。
	Description *string `json:"description,omitempty"`
}

func (o UpdateProduct) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProduct struct{}"
	}

	return strings.Join([]string{"UpdateProduct", string(data)}, " ")
}
