package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowDeviceResponse struct {

	// 资源空间ID。
	AppId *string `json:"app_id,omitempty"`

	// 资源空间名称。
	AppName *string `json:"app_name,omitempty"`

	// 设备ID，用于唯一标识一个设备。在注册设备时直接指定，或者由物联网平台分配获得。由物联网平台分配时，生成规则为\"product_id\" + \"_\" + \"node_id\"拼接而成。
	DeviceId *string `json:"device_id,omitempty"`

	// 设备标识码，通常使用IMEI、MAC地址或Serial No作为node_id。
	NodeId *string `json:"node_id,omitempty"`

	// 网关ID，用于标识设备所属的父设备，即父设备的设备ID。当设备是直连设备时，gateway_id与设备的device_id一致。当设备是非直连设备时，gateway_id为设备所关联的父设备的device_id。
	GatewayId *string `json:"gateway_id,omitempty"`

	// 设备名称。
	DeviceName *string `json:"device_name,omitempty"`

	// 设备节点类型。 - ENDPOINT：非直连设备。 - GATEWAY：直连设备或网关。 - UNKNOWN：未知。
	NodeType *string `json:"node_type,omitempty"`

	// 设备的描述信息。
	Description *string `json:"description,omitempty"`

	// 设备的固件版本。
	FwVersion *string `json:"fw_version,omitempty"`

	// 设备的软件版本。
	SwVersion *string `json:"sw_version,omitempty"`

	// 设备的sdk信息。
	DeviceSdkVersion *string `json:"device_sdk_version,omitempty"`

	AuthInfo *AuthInfo `json:"auth_info,omitempty"`

	// 设备关联的产品ID，用于唯一标识一个产品模型。
	ProductId *string `json:"product_id,omitempty"`

	// 设备关联的产品名称。
	ProductName *string `json:"product_name,omitempty"`

	// 设备的状态。 - ONLINE：设备在线。 - OFFLINE：设备离线。 - ABNORMAL：设备异常。 - INACTIVE：设备未激活。 - FROZEN：设备冻结。
	Status *string `json:"status,omitempty"`

	// 在物联网平台注册设备的时间。格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 设备的标签列表。
	Tags *[]TagV5Dto `json:"tags,omitempty"`

	// 设备扩展信息。用户可以自定义任何想要的扩展信息，如果在创建设备时为子设备指定该字段，将会通过MQTT接口“平台通知网关子设备新增“将该信息通知给网关。
	ExtensionInfo  *interface{} `json:"extension_info,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ShowDeviceResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDeviceResponse struct{}"
	}

	return strings.Join([]string{"ShowDeviceResponse", string(data)}, " ")
}
