package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 参数服务对象。
type ServiceCommandPara struct {

	// **参数说明**：参数的名称。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	ParaName string `json:"para_name"`

	// **参数说明**：参数的数据类型。 **取值范围**：int，long，decimal，string，DateTime，jsonObject，enum，boolean，string list。
	DataType string `json:"data_type"`

	// **参数说明**：参数是否必选。默认为false。
	Required *bool `json:"required,omitempty"`

	// **参数说明**：参数的枚举值列表。
	EnumList *[]string `json:"enum_list,omitempty"`

	// **参数说明**：参数的最小值。 **取值范围**：长度1-16。
	Min *string `json:"min,omitempty"`

	// **参数说明**：参数的最大值。 **取值范围**：长度1-16。
	Max *string `json:"max,omitempty"`

	// **参数说明**：参数的最大长度。
	MaxLength *int32 `json:"max_length,omitempty"`

	// **参数说明**：参数的步长。
	Step *float64 `json:"step,omitempty"`

	// **参数说明**：参数的单位。 **取值范围**：长度不超过16。
	Unit *string `json:"unit,omitempty"`

	// **参数说明**：参数的描述。 **取值范围**：长度不超过128，只允许中文、字母、数字、空白字符、以及_?'#().,;&%@!- ，、：；。/等字符的组合。
	Description *string `json:"description,omitempty"`
}

func (o ServiceCommandPara) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServiceCommandPara struct{}"
	}

	return strings.Join([]string{"ServiceCommandPara", string(data)}, " ")
}
