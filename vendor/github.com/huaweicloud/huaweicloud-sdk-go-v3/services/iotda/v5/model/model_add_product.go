package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddProduct 创建产品信息结构体。
type AddProduct struct {

	// **参数说明**：产品ID，资源空间下唯一。常app_id一起使用，用于唯一标识一个产品。如果携带此参数，平台将产品ID设置为该参数值；如果不携带此参数，产品ID在物联网平台创建产品后由平台分配获得。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	ProductId *string `json:"product_id,omitempty"`

	// **参数说明**：产品名称。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Name string `json:"name"`

	// **参数说明**：设备类型。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	DeviceType string `json:"device_type"`

	// **参数说明**：设备使用的协议类型。 **取值范围**：MQTT，CoAP，HTTP，HTTPS，Modbus，ONVIF，OPC-UA，OPC-DA，Other，TCP，UDP。
	ProtocolType string `json:"protocol_type"`

	// **参数说明**：设备上报数据的格式。 **取值范围**： - json：JSON格式 - binary：二进制码流格式 默认值json。
	DataFormat string `json:"data_format"`

	// **参数说明**：设备的服务能力列表。 **取值范围**：数组长度大小不超过500，内容大小不超过500k。
	ServiceCapabilities []ServiceCapability `json:"service_capabilities"`

	// **参数说明**：厂商名称。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	ManufacturerName *string `json:"manufacturer_name,omitempty"`

	// **参数说明**：设备所属行业。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	Industry *string `json:"industry,omitempty"`

	// **参数说明**：产品的描述信息。 **取值范围**：长度不超过128，只允许中文、字母、数字、空白字符、以及_?'#().,;&%@!- ，、：；。/等字符的组合。
	Description *string `json:"description,omitempty"`

	// **参数说明**：资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的产品归属到哪个资源空间下，否则创建的产品将会归属到[[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)](tag:hws)[[默认资源空间](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_0006.html#section0)](tag:hws_hk)下。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`
}

func (o AddProduct) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddProduct struct{}"
	}

	return strings.Join([]string{"AddProduct", string(data)}, " ")
}
