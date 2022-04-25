package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AddressGroup struct {

	// 功能说明：地址组唯一标识 取值范围：合法UUID的字符串
	Id string `json:"id"`

	// 功能说明：地址组名称 取值范围：0-64个字符，支持数字、字母、中文、_(下划线)、-（中划线）、.（点）
	Name string `json:"name"`

	// 功能说明：地址组描述信息 取值范围：0-255个字符 约束：不能包含“<”和“>”。
	Description string `json:"description"`

	// 功能说明：地址组可包含地址集 取值范围：可以是单个ip地址，ip地址范围，ip地址cidr 约束：当前一个地址组ip_set数量限制默认值为20，即配置的ip地址、ip地址范围或ip地址cidr的总数默认限制20
	IpSet []string `json:"ip_set"`

	// 功能说明：IP地址组ip版本 取值范围：4, 表示ipv4地址组；6, 表示ipv6地址组
	IpVersion int32 `json:"ip_version"`

	// 功能说明：地址组创建时间 取值范围：UTC时间格式：yyyy-MM-ddTHH:mm:ss；系统自动生成
	CreatedAt *sdktime.SdkTime `json:"created_at"`

	// 功能描述：地址组最近一次更新资源的时间 取值范围：UTC时间格式：yyyy-MM-ddTHH:mm:ss；系统自动生成
	UpdatedAt *sdktime.SdkTime `json:"updated_at"`

	// 功能说明：资源所属项目ID
	TenantId string `json:"tenant_id"`
}

func (o AddressGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddressGroup struct{}"
	}

	return strings.Join([]string{"AddressGroup", string(data)}, " ")
}
