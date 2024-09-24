package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeviceResource 预调配模板设备资源详情结构体。
type DeviceResource struct {
	DeviceName *ParameterRef `json:"device_name,omitempty"`

	NodeId *ParameterRef `json:"node_id"`

	// **参数说明**：设备所属的产品id，可以是一个明确的静态字符串id，也可以是动态的模板参数引用 - 明确的静态字符串：\"642bf260f2f9030e44210d8d\"。**取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。\" - 参数引用: {\"ref\" : \"iotda::certificate::country\"}
	ProductId *interface{} `json:"product_id"`

	// **参数说明**：设备绑定的标签列表
	Tags *[]TagRef `json:"tags,omitempty"`
}

func (o DeviceResource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeviceResource struct{}"
	}

	return strings.Join([]string{"DeviceResource", string(data)}, " ")
}
