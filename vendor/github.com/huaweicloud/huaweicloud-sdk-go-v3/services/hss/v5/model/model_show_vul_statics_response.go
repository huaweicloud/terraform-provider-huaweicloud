package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowVulStaticsResponse Response Object
type ShowVulStaticsResponse struct {

	// 需紧急修复的漏洞数
	NeedUrgentRepair *int32 `json:"need_urgent_repair,omitempty"`

	// 未完成修复的漏洞数
	Unrepair *int32 `json:"unrepair,omitempty"`

	// 存在漏洞的服务器数
	ExistedVulHosts *int32 `json:"existed_vul_hosts,omitempty"`

	// 今日处理漏洞数
	TodayHandle *int32 `json:"today_handle,omitempty"`

	// 累计处理漏洞数
	AllHandle *int32 `json:"all_handle,omitempty"`

	// 已支持漏洞数
	Supported *int32 `json:"supported,omitempty"`

	// 漏洞库更新时间
	VulLibraryUpdateTime *int64 `json:"vul_library_update_time,omitempty"`
	HttpStatusCode       int    `json:"-"`
}

func (o ShowVulStaticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowVulStaticsResponse struct{}"
	}

	return strings.Join([]string{"ShowVulStaticsResponse", string(data)}, " ")
}
