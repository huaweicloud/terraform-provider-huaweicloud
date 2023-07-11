package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PropertiesDto 属性数据
type PropertiesDto struct {

	// **参数说明**：MQTT 5.0版本请求和响应模式中的相关数据，可选。用户可以通过该参数配置MQTT协议请求和响应模式中的相关数据。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	CorrelationData *string `json:"correlation_data,omitempty"`

	// **参数说明**：MQTT 5.0版本请求和响应模式中的响应主题，可选。用户可以通过该参数配置MQTT协议请求和响应模式中的响应主题。 **取值范围**：长度不超过128, 只允许字母、数字、以及_-?=$#+/等字符的组合。
	ResponseTopic *string `json:"response_topic,omitempty"`

	// **参数说明**：用户自定义属性，可选。用户可以通过该参数配置用户自定义属性。可以配置的最大自定义属性数量为20。
	UserProperties *[]UserPropDto `json:"user_properties,omitempty"`
}

func (o PropertiesDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PropertiesDto struct{}"
	}

	return strings.Join([]string{"PropertiesDto", string(data)}, " ")
}
