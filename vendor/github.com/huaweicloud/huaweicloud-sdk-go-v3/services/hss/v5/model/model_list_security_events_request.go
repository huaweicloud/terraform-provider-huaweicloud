package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ListSecurityEventsRequest struct {

	// region id
	Region string `json:"region"`

	// 租户企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 查询时间范围天数，与自定义查询时间begin_time，end_time互斥
	LastDays *int32 `json:"last_days,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 容器名称
	ContainerName *string `json:"container_name,omitempty"`

	// 偏移量：指定返回记录的开始位置，必须为数字，取值范围为大于或等于0，默认0
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`

	// 事件类型，包含如下:   - 1001 : 恶意软件   - 1010 : Rootkit   - 1011 : 勒索软件   - 1015 : Webshell   - 1017 : 反弹Shell   - 2001 : 一般漏洞利用   - 2047 : Redis漏洞利用   - 2048 : Hadoop漏洞利用   - 2049 : MySQL漏洞利用   - 3002 : 文件提权   - 3003 : 进程提权   - 3004 : 关键文件变更   - 3005 : 文件/目录变更   - 3007 : 进程异常行为   - 3015 : 高危命令执行   - 3018 : 异常Shell   - 3027 : Crontab可疑任务   - 4002 : 暴力破解   - 4004 : 异常登录   - 4006 : 非法系统账号
	EventTypes *[]int32 `json:"event_types,omitempty"`

	// 处置状态，包含如下:   - unhandled ：未处理   - handled : 已处理
	HandleStatus *string `json:"handle_status,omitempty"`

	// 威胁等级，包含如下:   - Security ：安全   - Low : 低危   - Medium : 中危   - High : 高危   - Critical : 危急
	Severity *string `json:"severity,omitempty"`

	// 事件类别，包含如下:   - host : 主机安全事件   - container : 容器安全事件
	Category string `json:"category"`

	// 自定义查询时间，与查询时间范围天数互斥，查询时间段的起始时间，毫秒级时间戳，end_time减去begin_time小于等于2天，与查询时间范围天数互斥
	BeginTime *string `json:"begin_time,omitempty"`

	// 自定义时间，查询时间段的终止时间，毫秒级时间戳，end_time减去begin_time小于等于2天，与查询时间范围天数互斥
	EndTime *string `json:"end_time,omitempty"`
}

func (o ListSecurityEventsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSecurityEventsRequest struct{}"
	}

	return strings.Join([]string{"ListSecurityEventsRequest", string(data)}, " ")
}
