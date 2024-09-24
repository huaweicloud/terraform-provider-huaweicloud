package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListProtectionServerRequest Request Object
type ListProtectionServerRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 操作系统类型，包含如下2种。   - Linux ：Linux。   - Windows ：Windows。
	OsType *string `json:"os_type,omitempty"`

	// 服务器IP地址
	HostIp *string `json:"host_ip,omitempty"`

	// 主机状态，包含如下3种。   - 不传参默认为全部。   - ACTIVE ：正在运行。   - SHUTOFF ：关机。
	HostStatus *string `json:"host_status,omitempty"`

	// 查询时间范围天数，最近7天为last_days=7，若不填，则默认查询一天内的防护事件和已有备份数
	LastDays *int32 `json:"last_days,omitempty"`
}

func (o ListProtectionServerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProtectionServerRequest struct{}"
	}

	return strings.Join([]string{"ListProtectionServerRequest", string(data)}, " ")
}
