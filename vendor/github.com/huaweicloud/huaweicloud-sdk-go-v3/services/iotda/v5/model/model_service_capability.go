package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ServiceCapability结构体。
type ServiceCapability struct {

	// **参数说明**：设备的服务ID。注：产品内不允许重复。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-$等字符的组合。
	ServiceId string `json:"service_id"`

	// **参数说明**：设备的服务类型。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_?'#().,&%@!-$等字符的组合。
	ServiceType string `json:"service_type"`

	// **参数说明**：设备服务支持的属性列表。 **取值范围**：数组长度大小不超过500。
	Properties *[]ServiceProperty `json:"properties,omitempty"`

	// **参数说明**：设备服务支持的命令列表。 **取值范围**：数组长度大小不超过500。
	Commands *[]ServiceCommand `json:"commands,omitempty"`

	// **参数说明**：设备服务支持的事件列表。 **取值范围**：数组长度大小不超过500。
	Events *[]ServiceEvent `json:"events,omitempty"`

	// **参数说明**：设备服务的描述信息。 **取值范围**：长度不超过128，只允许中文、字母、数字、空白字符、以及_?'#().,;&%@!- ，、：；。/等字符的组合。
	Description *string `json:"description,omitempty"`

	// **参数说明**：指定设备服务是否必选。目前本字段为非功能性字段，仅起到标识作用。 **取值范围**： - Master：主服务 - Mandatory：必选服务 - Optional：可选服务 默认值为Optional。
	Option *string `json:"option,omitempty"`
}

func (o ServiceCapability) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServiceCapability struct{}"
	}

	return strings.Join([]string{"ServiceCapability", string(data)}, " ")
}
