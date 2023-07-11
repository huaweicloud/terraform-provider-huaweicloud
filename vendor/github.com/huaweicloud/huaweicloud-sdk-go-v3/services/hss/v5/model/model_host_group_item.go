package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HostGroupItem 服务器组信息
type HostGroupItem struct {

	// 服务器组ID
	GroupId *string `json:"group_id,omitempty"`

	// 服务器组名称
	GroupName *string `json:"group_name,omitempty"`

	// 关联服务器数
	HostNum *int32 `json:"host_num,omitempty"`

	// 有风险服务器数
	RiskHostNum *int32 `json:"risk_host_num,omitempty"`

	// 未防护服务器数
	UnprotectHostNum *int32 `json:"unprotect_host_num,omitempty"`

	// 服务器ID列表
	HostIdList *[]string `json:"host_id_list,omitempty"`

	// 是否是线下数据中心服务器组
	IsOutside *bool `json:"is_outside,omitempty"`
}

func (o HostGroupItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostGroupItem struct{}"
	}

	return strings.Join([]string{"HostGroupItem", string(data)}, " ")
}
