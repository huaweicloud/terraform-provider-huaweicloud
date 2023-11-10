package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddDeviceGroupDto 创建设备组请求结构体
type AddDeviceGroupDto struct {

	// **参数说明**：设备组名称，单个资源空间下不可重复。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_? '#().,&%@!-等字符的组合。
	Name *string `json:"name,omitempty"`

	// **参数说明**：设备组描述。 **取值范围**：长度不超过64，只允许中文、字母、数字、以及_? '#().,&%@!-等字符的组合。
	Description *string `json:"description,omitempty"`

	// **参数说明**：父设备组ID，携带该参数时表示在该设备组下创建一个子设备组,动态群组不支持该参数。 **取值范围**：长度不超过36，十六进制字符串和连接符（-）的组合。
	SuperGroupId *string `json:"super_group_id,omitempty"`

	// **参数说明**：资源空间ID。此参数为非必选参数，存在多资源空间的用户需要使用该接口时，建议携带该参数指定创建的设备组归属到哪个资源空间下，否则创建的设备组将会归属到[[默认资源空间](https://support.huaweicloud.com/usermanual-iothub/iot_01_0006.html#section0)](tag:hws)[[默认资源空间](https://support.huaweicloud.com/intl/zh-cn/usermanual-iothub/iot_01_0006.html#section0)](tag:hws_hk)下。 **取值范围**：长度不超过36，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AppId *string `json:"app_id,omitempty"`

	// **参数说明**：设备组类型，默认为静态设备组；当设备组类型为动态设备组时，需要填写动态设备组组规则
	GroupType *string `json:"group_type,omitempty"`

	// **参数说明**：动态设备组规则语法和高级搜索保持一致，只需要填写where 子句内容，其余子句无需填写，todo补充说明
	DynamicGroupRule *string `json:"dynamic_group_rule,omitempty"`
}

func (o AddDeviceGroupDto) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddDeviceGroupDto struct{}"
	}

	return strings.Join([]string{"AddDeviceGroupDto", string(data)}, " ")
}
