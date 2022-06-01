package sources

import "github.com/chnsz/golangsdk/pagination"

// SourceServerPage is a single page result representing a query by offset page.
type SourceServerPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a SourceServerPage struct is empty.
func (r SourceServerPage) IsEmpty() (bool, error) {
	servers, err := ExtractSourceServers(r)
	return len(servers) == 0, err
}

// ExtractSourceServers is a method to extract the list of source servers details.
func ExtractSourceServers(r pagination.Page) ([]SourceServer, error) {
	var s struct {
		Servers []SourceServer `json:"source_servers"`
	}

	rst := r.(SourceServerPage).Result
	err := rst.ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s.Servers, nil
}

// SourceServer 源端服务器列表信息
type SourceServer struct {
	// 源端服务器id
	Id string `json:"id"`
	// 源端服务器的ip地址
	Ip string `json:"ip"`
	// 源端服务器名称
	Name string `json:"name"`
	// 源端服务器状态
	State string `json:"state"`
	// 源端服务器与主机迁移服务端是否连接
	Connected bool `json:"connected"`
	// 源端服务器的注册时间
	AddDate int64 `json:"add_date"`
	// Agent 版本
	AgentVersion string `json:"agent_version"`
	// 操作系统类型: WINDOWS/LINUX
	OsType string `json:"os_type"`
	// 系统详细版本号,如CENTOS7.6等
	OsVersion string `json:"os_version"`
	// 是否是OEM操作系统(Windows)
	OemSystem bool `json:"oem_system"`
	// 企业项目id
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// 源端CPU核心数
	CPU int `json:"cpu_quantity"`
	// 源端物理内存大小(单位:字节)
	Memory int64 `json:"memory"`
	// 源端列表中关联的任务
	CurrentTask Task `json:"current_task"`
	// 源端校验检查项列表
	Checks []Check `json:"checks"`
	// 推荐的目的端服务器配置
	InitTargetServer InitTargetServer `json:"init_target_server"`
	// 迁移周期
	MigrationCycle string `json:"migration_cycle"`
	// 已复制的大小(单位:字节)
	Replicatesize int64 `json:"replicatesize"`
	// 需要迁移的数据量总大小(单位:字节)
	Totalsize int64 `json:"totalsize"`
	// 迁移周期(migration_cycle)上一次变化的时间
	StageActionTime int64 `json:"stage_action_time"`
	// Agent上一次连接状态发生变化的时间
	LastVisitTime int64 `json:"last_visit_time"`
	// 源端状态(state)上次发生变化的时间
	StateActionTime int64 `json:"state_action_time"`
}

// Task 源端列表中关联的任务
type Task struct {
	// 任务id
	Id string `json:"id"`
	// 任务名称
	Name string `json:"name"`
	// 任务类型
	Type string `json:"type"`
	// 任务状态
	State string `json:"state"`
	// 预估结束时间
	EstimateCompleteTime int64 `json:"estimate_complete_time"`
	// 开始时间
	StartDate int64 `json:"start_date"`
	// 限速
	SpeedLimit float32 `json:"speed_limit"`
	// 迁移速率
	MigrateSpeed float32 `json:"migrate_speed"`
	// 压缩率
	CompressRate float32 `json:"compress_rate"`
	// 是否启动虚拟机
	StartTargetServer bool `json:"start_target_server"`
	// 虚拟机模板id
	VmTemplateId string `json:"vm_template_id"`
	// region_id
	RegionId string `json:"region_id"`
	// 项目名称
	ProjectName string `json:"project_name"`
	// 项目id
	ProjectId string `json:"project_id"`
	// 目的端
	TargetServer TargetServer `json:"target_server"`
	// 克隆服务器端
	CloneServer CloneServer `json:"clone_server"`
	// 日志收集状态
	LogCollectStatus string `json:"log_collect_status"`
	// 是否使用已有虚拟机
	ExistServer bool `json:"exist_server"`
	// 是否使用公网ip
	UsePublicIp bool `json:"use_public_ip"`
	// 已迁移时长
	RemainSeconds int64 `json:"remain_seconds"`
}

// TargetServer 目的端
type TargetServer struct {
	// 目的端服务器ID
	VmId string `json:"vm_id"`
	// 目的端服务器名称
	Name string `json:"name"`
}

// CloneServer 克隆服务器类
type CloneServer struct {
	// 克隆服务器ID
	VmId string `json:"vm_id"`
	// 克隆虚拟机的名称
	Name string `json:"name"`
	// 克隆错误信息
	CloneError string `json:"clone_error"`
	// 克隆状态
	CloneState string `json:"clone_state"`
	// 克隆错误信息描述
	ErrorMsg string `json:"error_msg"`
}

// Check 源端校验项
type Check struct {
	// 该检查项的ID
	Id int `json:"id"`
	// 参数
	Params []string `json:"params"`
	// 检查项名称
	Name string `json:"name"`
	// 检查结果
	Result string `json:"result"`
	// 检查不通过的错误码
	ErrorCode string `json:"error_code"`
	// 检查不通过的错误参数
	ErrorParams string `json:"error_params"`
}

// InitTargetServer 推荐的目的端服务器配置
type InitTargetServer struct {
	// 推荐的目的端服务器的磁盘信息
	Disks []ServerDisk `json:"disks"`
}

// ServerDisk 目的端服务器关联磁盘
type ServerDisk struct {
	// 磁盘名称
	Name string `json:"name"`
	// 磁盘大小,单位:字节
	Size int64 `json:"size"`
	// 磁盘的作用
	DeviceUse string `json:"device_use"`
	// 逻辑卷信息
	PhysicalVolumes []PhysicalVolumes `json:"physical_volumes"`
}

// PhysicalVolumes 物理分区
type PhysicalVolumes struct {
	// 分区类型
	DeviceType string `json:"device_use"`
	// 文件系统
	FileSystem string `json:"file_system"`
	// 编号
	Index int `json:"index"`
	// 挂载点
	MountPoint string `json:"mount_point"`
	// 名称
	Name string `json:"name"`
	// 大小
	Size int64 `json:"size"`
	// 使用大小
	UsedSize int64 `json:"used_size"`
	// uuid
	UUID string `json:"uuid"`
}
