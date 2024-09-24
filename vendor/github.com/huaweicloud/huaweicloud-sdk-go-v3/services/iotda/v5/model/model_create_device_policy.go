package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDevicePolicy 添加策略请求体。
type CreateDevicePolicy struct {

	// **参数说明**：策略名称。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	PolicyName string `json:"policy_name"`

	// **参数说明**：资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的设备归属到哪个资源空间下，否则创建的设备将会归属到[[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)](tag:hws)[[默认资源空间](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_0006.html#section0)](tag:hws_hk)下。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：策略文档。
	Statement []Statement `json:"statement"`
}

func (o CreateDevicePolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDevicePolicy struct{}"
	}

	return strings.Join([]string{"CreateDevicePolicy", string(data)}, " ")
}
