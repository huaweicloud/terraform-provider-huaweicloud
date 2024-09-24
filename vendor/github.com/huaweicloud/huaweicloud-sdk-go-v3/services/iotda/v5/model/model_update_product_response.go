package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateProductResponse Response Object
type UpdateProductResponse struct {

	// 资源空间ID。
	AppId *string `json:"app_id,omitempty"`

	// 资源空间名称。
	AppName *string `json:"app_name,omitempty"`

	// 产品ID，用于唯一标识一个产品，在物联网平台创建产品后由平台分配获得。
	ProductId *string `json:"product_id,omitempty"`

	// 产品名称。
	Name *string `json:"name,omitempty"`

	// 设备类型。
	DeviceType *string `json:"device_type,omitempty"`

	// 设备使用的协议类型。取值范围：MQTT，CoAP，HTTP，HTTPS，Modbus，ONVIF， OPC-UA，OPC-DA，Other，TCP，UDP。
	ProtocolType *string `json:"protocol_type,omitempty"`

	// 设备上报数据的格式，取值范围：json，binary。
	DataFormat *string `json:"data_format,omitempty"`

	// 厂商名称。
	ManufacturerName *string `json:"manufacturer_name,omitempty"`

	// 设备所属行业。
	Industry *string `json:"industry,omitempty"`

	// 产品的描述信息。
	Description *string `json:"description,omitempty"`

	// 设备的服务能力列表。
	ServiceCapabilities *[]ServiceCapability `json:"service_capabilities,omitempty"`

	// 在物联网平台创建产品的时间，格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	CreateTime     *string `json:"create_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateProductResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProductResponse struct{}"
	}

	return strings.Join([]string{"UpdateProductResponse", string(data)}, " ")
}
