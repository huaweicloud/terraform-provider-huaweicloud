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

	// 功能说明：地址组最大条目数，限制地址组可以包含的地址数量 取值范围：0-20 默认值：20
	MaxCapacity int32 `json:"max_capacity"`

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

	// 功能说明：企业项目ID。 取值范围：最大长度36字节，带“-”连字符的UUID格式，或者是字符串“0”。“0”表示默认企业项目。
	EnterpriseProjectId string `json:"enterprise_project_id"`

	// IP地址组资源标签
	Tags []ResourceTag `json:"tags"`

	// 功能说明：地址组状态 取值范围：       NORMAL：正常       UPDATING：更新中       UPDATE_FAILED：更新失败 默认值：NORMAL 约束：当地址组处于UPDATING（更新中）状态时，不允许再次更新
	Status string `json:"status"`

	// 功能说明：地址组状态详情信息
	StatusMessage string `json:"status_message"`

	// 功能说明：地址组包含的地址集及其备注信息
	IpExtraSet []IpExtraSetRespOption `json:"ip_extra_set"`
}

func (o AddressGroup) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddressGroup struct{}"
	}

	return strings.Join([]string{"AddressGroup", string(data)}, " ")
}
