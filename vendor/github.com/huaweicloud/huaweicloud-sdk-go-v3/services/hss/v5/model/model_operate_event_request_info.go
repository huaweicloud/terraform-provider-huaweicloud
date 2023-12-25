package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OperateEventRequestInfo 操作的事件
type OperateEventRequestInfo struct {

	// 事件分类，包含如下:   - container_1001 : 容器命名空间   - container_1002 : 容器开放端口   - container_1003 : 容器安全选项   - container_1004 : 容器挂载目录   - containerescape_0001 : 容器高危系统调用   - containerescape_0002 : Shocker攻击   - containerescape_0003 : DirtCow攻击   - containerescape_0004 : 容器文件逃逸攻击   - dockerfile_001 : 用户自定义容器保护文件被修改   - dockerfile_002 : 容器文件系统可执行文件被修改   - dockerproc_001 : 容器进程异常事件上报   - fileprotect_0001 : 文件提权   - fileprotect_0002 : 关键文件变更   - fileprotect_0003 : 关键文件路径变更   - fileprotect_0004 : 文件/目录变更   - av_1002 : 病毒   - av_1003 : 蠕虫   - av_1004 : 木马   - av_1005 : 僵尸网络   - av_1006 : 后门   - av_1007 : 间谍软件   - av_1008 : 恶意广告软件   - av_1009 : 钓鱼   - av_1010 : Rootkit   - av_1011 : 勒索软件   - av_1012 : 黑客工具   - av_1013 : 灰色软件   - av_1015 : Webshell   - av_1016 : 挖矿软件   - login_0001 : 尝试暴力破解   - login_0002 : 爆破成功   - login_1001 : 登录成功   - login_1002 : 异地登录   - login_1003 : 弱口令   - malware_0001 : shell变更事件上报   - malware_0002 : 反弹shell事件上报   - malware_1001 : 恶意程序   - procdet_0001 : 进程异常行为检测   - procdet_0002 : 进程提权   - procreport_0001 : 危险命令   - user_1001 : 账号变更   - user_1002 : 风险账号   - vmescape_0001 : 虚拟机敏感命令执行   - vmescape_0002 : 虚拟化进程访问敏感文件   - vmescape_0003 : 虚拟机异常端口访问   - webshell_0001 : 网站后门   - network_1001 : 恶意挖矿   - network_1002 : 对外DDoS攻击   - network_1003 : 恶意扫描   - network_1004 : 敏感区域攻击   - ransomware_0001 : 勒索攻击   - ransomware_0002 : 勒索攻击   - ransomware_0003 : 勒索攻击   - fileless_0001 : 进程注入   - fileless_0002 : 动态库注入进程   - fileless_0003 : 关键配置变更   - fileless_0004 : 环境变量变更   - fileless_0005 : 内存文件进程   - fileless_0006 : vdso劫持   - crontab_1001 : Crontab可疑任务   - vul_exploit_0001 : Redis漏洞利用攻击   - vul_exploit_0002 : Hadoop漏洞利用攻击   - vul_exploit_0003 : MySQL漏洞利用攻击   - rootkit_0001 : 可疑rootkit文件   - rootkit_0002 : 可疑内核模块   - RASP_0004 : 上传Webshell   - RASP_0018 : 无文件Webshell   - blockexec_001 : 已知勒索攻击   - hips_0001 : Windows Defender防护被禁用   - hips_0002 : 可疑的黑客工具   - hips_0003 : 可疑的勒索加密行为   - hips_0004 : 隐藏账号创建   - hips_0005 : 读取用户密码凭据   - hips_0006 : 可疑的SAM文件导出   - hips_0007 : 可疑shadow copy删除操作   - hips_0008 : 备份文件删除   - hips_0009 : 可疑勒索病毒操作注册表   - hips_0010 : 可疑的异常进程行为   - hips_0011 : 可疑的扫描探测   - hips_0012 : 可疑的勒索病毒脚本运行   - hips_0013 : 可疑的挖矿命令执行   - hips_0014 : 可疑的禁用windows安全中心   - hips_0015 : 可疑的停止防火墙服务行为   - hips_0016 : 可疑的系统自动恢复禁用   - hips_0017 : Offies 创建可执行文件   - hips_0018 : 带宏Offies文件异常创建   - hips_0019 : 可疑的注册表操作   - hips_0020 : Confluence远程代码执行   - hips_0021 : MSDT远程代码执行   - portscan_0001 : 通用端口扫描   - portscan_0002 : 秘密端口扫描   - k8s_1001 : Kubernetes事件删除   - k8s_1002 : 创建特权Pod   - k8s_1003 : Pod中使用交互式shell   - k8s_1004 : 创建敏感目录Pod   - k8s_1005 : 创建主机网络的Pod   - k8s_1006 : 创建主机Pid空间的Pod   - k8s_1007 : 普通pod访问APIserver认证失败   - k8s_1008 : 普通Pod通过Curl访问APIServer   - k8s_1009 : 系统管理空间执行exec   - k8s_1010 : 系统管理空间创建Pod   - k8s_1011 : 创建静态Pod   - k8s_1012 : 创建DaemonSet   - k8s_1013 : 创建集群计划任务   - k8s_1014 : Secrets操作   - k8s_1015 : 枚举用户可执行的操作   - k8s_1016 : 高权限RoleBinding或ClusterRoleBinding   - k8s_1017 : ServiceAccount创建   - k8s_1018 : 创建Cronjob   - k8s_1019 : Pod中exec使用交互式shell   - k8s_1020 : 无权限访问Apiserver   - k8s_1021 : 使用curl访问APIServer   - k8s_1022 : Ingress漏洞   - k8s_1023 : 中间人攻击   - k8s_1024 : 蠕虫挖矿木马   - k8s_1025 : K8s事件删除   - k8s_1026 : SelfSubjectRulesReview场景   - imgblock_0001 : 镜像白名单阻断   - imgblock_0002 : 镜像黑名单阻断   - imgblock_0003 : 镜像标签白名单阻断   - imgblock_0004 : 镜像标签黑名单阻断   - imgblock_0005 : 创建容器白名单阻断   - imgblock_0006 : 创建容器黑名单阻断   - imgblock_0007 : 容器mount proc阻断   - imgblock_0008 : 容器seccomp unconfined阻断   - imgblock_0009 : 容器特权阻断   - imgblock_0010 : 容器capabilities阻断
	EventClassId string `json:"event_class_id"`

	// 事件编号
	EventId string `json:"event_id"`

	// 事件类型，包含如下:   - 1001 : 通用恶意软件   - 1002 : 病毒   - 1003 : 蠕虫   - 1004 : 木马   - 1005 : 僵尸网络   - 1006 : 后门   - 1010 : Rootkit   - 1011 : 勒索软件   - 1012 ：黑客工具   - 1015 : Webshell   - 1016 : 挖矿   - 1017 : 反弹Shell   - 2001 : 一般漏洞利用   - 2012 : 远程代码执行   - 2047 : Redis漏洞利用   - 2048 : Hadoop漏洞利用   - 2049 : MySQL漏洞利用   - 3002 : 文件提权   - 3003 : 进程提权   - 3004 : 关键文件变更   - 3005 : 文件/目录变更   - 3007 : 进程异常行为   - 3015 : 高危命令执行   - 3018 : 异常Shell   - 3027 : Crontab可疑任务   - 3029 ：系统安全防护被禁用   - 3030 ：备份删除   - 3031 ：异常注册表操作   - 3036 : 容器镜像阻断   - 4002 : 暴力破解   - 4004 : 异常登录   - 4006 : 非法系统账号   - 4014 : 用户账号添加   - 4020 : 用户密码窃取   - 6002 : 端口扫描   - 6003 : 主机扫描   - 13001 : Kubernetes事件删除   - 13002 : Pod异常行为   - 13003 : 枚举用户信息   - 13004 : 绑定集群用户角色
	EventType int32 `json:"event_type"`

	// 发生时间，毫秒
	OccurTime int64 `json:"occur_time"`

	// 操作详情信息列表，当 operate_type 为 add_to_alarm_whitelist 或 remove_from_alarm_whitelist 时，必传 keyword 和 hash；当 operate_type 为 add_to_login_whitelist 或 remove_from_login_whitelist 时，必传 login_ip， private_ip 和 login_user_name；当 operate_type 为 isolate_and_kill 或 do_not_isolate_or_kill 时，必传 agent_id，file_hash，file_path，process_pid；其余情况可不填写内容。
	OperateDetailList []EventDetailRequestInfo `json:"operate_detail_list"`
}

func (o OperateEventRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OperateEventRequestInfo struct{}"
	}

	return strings.Join([]string{"OperateEventRequestInfo", string(data)}, " ")
}
