package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设备影子数据属性结构体。
type DeviceShadowProperties struct {

	// 设备影子的属性数据，Json格式，里面是一个个键值对，每个键都是产品模型中属性的参数名(property_name)，目前如样例所示只支持一层结构。 **注意**：JSON结构的key当前不支持特殊字符：点(.)、dollar符号($)、空char(十六进制的ASCII码为00),key为以上特殊字符无法正常刷新设备影子
	Properties *interface{} `json:"properties,omitempty"`

	// 事件操作时间，格式：yyyyMMdd'T'HHmmss'Z'，如20151212T121212Z。
	EventTime *string `json:"event_time,omitempty"`
}

func (o DeviceShadowProperties) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceShadowProperties struct{}"
	}

	return strings.Join([]string{"DeviceShadowProperties", string(data)}, " ")
}
