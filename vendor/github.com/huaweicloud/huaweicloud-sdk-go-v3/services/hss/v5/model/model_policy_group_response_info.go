package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PolicyGroupResponseInfo 策略组详情
type PolicyGroupResponseInfo struct {

	// 策略组名
	GroupName *string `json:"group_name,omitempty"`

	// 策略组ID
	GroupId *string `json:"group_id,omitempty"`

	// 策略组的描述信息
	Description *string `json:"description,omitempty"`

	// 是否允许删除该策略组
	Deletable *bool `json:"deletable,omitempty"`

	// 关联服务器数
	HostNum *int32 `json:"host_num,omitempty"`

	// 是否是默认策略组
	DefaultGroup *bool `json:"default_group,omitempty"`

	// 支持的操作系统，包含如下:   - Linux ：支持Linux系统   - Windows : 支持Windows系统
	SupportOs *string `json:"support_os,omitempty"`

	// 支持的版本，包含如下:   - hss.version.basic ：基础版策略组   - hss.version.advanced : 专业版策略组   - hss.version.enterprise : 企业版策略组   - hss.version.premium : 旗舰版策略组   - hss.version.wtp : 网页防篡改版策略组   - hss.version.container.enterprise : 容器版策略组
	SupportVersion *string `json:"support_version,omitempty"`
}

func (o PolicyGroupResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PolicyGroupResponseInfo struct{}"
	}

	return strings.Join([]string{"PolicyGroupResponseInfo", string(data)}, " ")
}
