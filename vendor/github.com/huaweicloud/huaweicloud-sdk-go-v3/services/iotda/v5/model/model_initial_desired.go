package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 设备初始配置数据结构体。
type InitialDesired struct {

	// **参数说明**：设备的服务ID，在设备关联的产品模型中定义。 **取值范围**：长度不超过32，只允许中文、字母、数字、以及_?'#().,&%@!-等字符的组合。
	ServiceId string `json:"service_id"`

	// **参数说明**：设备初始配置属性数据，Json格式，里面是一个个键值对，每个键都是产品模型中属性的参数名(property_name)，目前如样例所示只支持一层结构；这里设置的属性值与产品中对应属性的默认值比对，如果不同，则将以该字段中设置的属性值为准写入到设备影子中；如果想要删除整个desired可以填写空object(例如\"desired\":{})，如果想要删除某一个属性期望值可以将属性置位null(例如{\"temperature\":null})
	Desired *interface{} `json:"desired"`
}

func (o InitialDesired) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InitialDesired struct{}"
	}

	return strings.Join([]string{"InitialDesired", string(data)}, " ")
}
