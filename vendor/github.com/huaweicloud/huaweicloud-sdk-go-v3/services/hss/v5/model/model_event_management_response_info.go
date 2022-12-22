package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 事件详情统计
type EventManagementResponseInfo struct {

	// 事件编号
	EventId *string `json:"event_id,omitempty"`

	// 事件分类，包含如下:   - container_1001 : 容器命名空间   - container_1002 : 容器开放端口   - container_1003 : 容器安全选项   - container_1004 : 容器挂载目录   - containerescape_0001 : 容器高危系统调用   - containerescape_0002 : Shocker攻击   - containerescape_0003 : DirtCow攻击   - containerescape_0004 : 容器文件逃逸攻击   - dockerfile_001 : 用户自定义容器保护文件被修改   - dockerfile_002 : 容器文件系统可执行文件被修改   - dockerproc_001 : 容器进程异常事件上报   - fileprotect_0001 : 文件提权   - fileprotect_0002 : 关键文件变更   - fileprotect_0003 : 关键文件路径变更   - fileprotect_0004 : 文件目录变更   - login_0001 : 尝试暴力破解   - login_0002 : 爆破成功   - login_1001 : 登录成功   - login_1002 : 异地登录   - login_1003 : 弱口令   - malware_0001 : shell变更事件上报   - malware_0002 : 反弹shell事件上报   - malware_1001 : 恶意程序   - procdet_0001 : 进程异常行为检测   - procdet_0002 : 进程提权   - procreport_0001 : 危险命令   - user_1001 : 账号变更   - user_1002 : 风险账号   - vmescape_0001 : 虚拟机敏感命令执行   - vmescape_0002 : 虚拟化进程访问敏感文件   - vmescape_0003 : 虚拟机异常端口访问   - webshell_0001 : 网站后门   - network_1001 : 恶意挖矿   - network_1002 : 对外DDoS攻击   - network_1003 : 恶意扫描   - network_1004 : 敏感区域攻击   - crontab_1001 : Crontab可疑任务   - vul_exploit_0001 : Redis漏洞利用攻击   - vul_exploit_0002 : Hadoop漏洞利用攻击   - vul_exploit_0003 : MySQL漏洞利用攻击
	EventClassId *string `json:"event_class_id,omitempty"`

	// 事件类型，包含如下:   - 1001 : 恶意软件   - 1010 : Rootkit   - 1011 : 勒索软件   - 1015 : Webshell   - 1017 : 反弹Shell   - 2001 : 一般漏洞利用   - 2047 : Redis漏洞利用   - 2048 : Hadoop漏洞利用   - 2049 : MySQL漏洞利用   - 3002 : 文件提权   - 3003 : 进程提权   - 3004 : 关键文件变更   - 3005 : 文件/目录变更   - 3007 : 进程异常行为   - 3015 : 高危命令执行   - 3018 : 异常Shell   - 3027 : Crontab可疑任务   - 4002 : 暴力破解   - 4004 : 异常登录   - 4006 : 非法系统账号
	EventType *int32 `json:"event_type,omitempty"`

	// 事件名称
	EventName *string `json:"event_name,omitempty"`

	// 威胁等级，包含如下:   - Security : 安全   - Low : 低危   - Medium : 中危   - High : 高危   - Critical : 危急
	Severity *string `json:"severity,omitempty"`

	// 容器实例名称
	ContainerName *string `json:"container_name,omitempty"`

	// 镜像名称
	ImageName *string `json:"image_name,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ID
	HostId *string `json:"host_id,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 弹性公网IP地址
	PublicIp *string `json:"public_ip,omitempty"`

	// 攻击阶段，包含如下：   - reconnaissance : 侦查跟踪   - weaponization : 武器构建   - delivery : 载荷投递   - exploit : 漏洞利用   - installation : 安装植入   - command_and_control : 命令与控制   - actions : 目标达成
	AttackPhase *string `json:"attack_phase,omitempty"`

	// 攻击标识，包含如下：   - attack_success : 攻击成功   - attack_attempt : 攻击尝试   - attack_blocked : 攻击被阻断   - abnormal_behavior : 异常行为   - collapsible_host : 主机失陷   - system_vulnerability : 系统脆弱性
	AttackTag *string `json:"attack_tag,omitempty"`

	// 发生时间，毫秒
	OccurTime *int64 `json:"occur_time,omitempty"`

	// 处理时间，毫秒
	HandleTime *int64 `json:"handle_time,omitempty"`

	// 处理状态，包含如下:   - unhandled ：未处理   - handled : 已处理
	HandleStatus *string `json:"handle_status,omitempty"`

	// 处理方式，包含如下:   - mark_as_handled : 手动处理   - ignore : 忽略   - add_to_alarm_whitelist : 加入告警白名单   - add_to_login_whitelist : 加入登录白名单   - isolate_and_kill : 隔离查杀
	HandleMethod *string `json:"handle_method,omitempty"`

	// 手动处理的备注
	Handler *string `json:"handler,omitempty"`

	// 支持的处理操作
	OperateAcceptList *[]string `json:"operate_accept_list,omitempty"`

	// 操作详情信息列表（页面不展示）
	OperateDetailList *[]EventDetailResponseInfo `json:"operate_detail_list,omitempty"`

	// 取证信息，json格式
	ForensicInfo *interface{} `json:"forensic_info,omitempty"`

	ResourceInfo *EventResourceResponseInfo `json:"resource_info,omitempty"`

	// 地理位置信息，json格式
	GeoInfo *interface{} `json:"geo_info,omitempty"`

	// 恶意软件信息，json格式
	MalwareInfo *interface{} `json:"malware_info,omitempty"`

	// 网络信息，json格式
	NetworkInfo *interface{} `json:"network_info,omitempty"`

	// 应用信息，json格式
	AppInfo *interface{} `json:"app_info,omitempty"`

	// 系统信息，json格式
	SystemInfo *interface{} `json:"system_info,omitempty"`

	// 事件扩展信息，json格式
	ExtendInfo *interface{} `json:"extend_info,omitempty"`

	// 处置建议
	Recommendation *string `json:"recommendation,omitempty"`

	// 进程信息列表
	ProcessInfoList *[]EventProcessResponseInfo `json:"process_info_list,omitempty"`

	// 用户信息列表
	UserInfoList *[]EventUserResponseInfo `json:"user_info_list,omitempty"`

	// 文件信息列表
	FileInfoList *[]EventFileResponseInfo `json:"file_info_list,omitempty"`
}

func (o EventManagementResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventManagementResponseInfo struct{}"
	}

	return strings.Join([]string{"EventManagementResponseInfo", string(data)}, " ")
}
