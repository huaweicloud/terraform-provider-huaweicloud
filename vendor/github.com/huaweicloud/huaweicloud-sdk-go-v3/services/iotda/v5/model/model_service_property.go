package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 属性服务对象。
type ServiceProperty struct {

	// **参数说明**：设备属性名称。注：设备服务内不允许重复。属性名称作为设备影子JSON文档中的key不支持特殊字符：点(.)、dollar符号($)、空char(十六进制的ASCII码为00)，如果包含了以上特殊字符则无法正常刷新影子文档。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	PropertyName string `json:"property_name"`

	// **参数说明**：设备属性的数据类型。 **取值范围**：int，long，decimal，string，DateTime，jsonObject，enum，boolean，string list。
	DataType string `json:"data_type"`

	// **参数说明**：设备属性是否必选。默认为false。
	Required *bool `json:"required,omitempty"`

	// **参数说明**：设备属性的枚举值列表。
	EnumList *[]string `json:"enum_list,omitempty"`

	// **参数说明**：设备属性的最小值。 **取值范围**：长度1-16。
	Min *string `json:"min,omitempty"`

	// **参数说明**：设备属性的最大值。 **取值范围**：长度1-16。
	Max *string `json:"max,omitempty"`

	// **参数说明**：设备属性的最大长度。
	MaxLength *int32 `json:"max_length,omitempty"`

	// **参数说明**：设备属性的步长。
	Step *float64 `json:"step,omitempty"`

	// **参数说明**：设备属性的单位。 **取值范围**：长度不超过16。
	Unit *string `json:"unit,omitempty"`

	// **参数说明**：设备属性的访问模式。 **取值范围**：RWE，RW，RE，WE，E，W，R。 - R：属性值可读 - W：属性值可写 - E：属性值可订阅，即属性值变化时上报事件
	Method string `json:"method"`

	// **参数说明**：设备属性的描述。 **取值范围**：长度不超过128，只允许中文、字母、数字、空白字符、以及_?'#().,;&%@!- ，、：；。/等字符的组合。
	Description *string `json:"description,omitempty"`

	// **参数说明**：设备属性的默认值。如果设置了默认值，使用该产品创建设备时，会将该属性的默认值写入到该设备的设备影子预期数据中，待设备上线时将该属性默认值下发给设备。
	DefaultValue *interface{} `json:"default_value,omitempty"`
}

func (o ServiceProperty) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServiceProperty struct{}"
	}

	return strings.Join([]string{"ServiceProperty", string(data)}, " ")
}
