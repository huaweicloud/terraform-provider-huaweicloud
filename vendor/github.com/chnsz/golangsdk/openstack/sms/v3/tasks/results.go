package tasks

// CreateResp is a auto create Response Object
type CreateResp struct {
	// 创建成功返回的任务id
	ID string `json:"id"`
}

// MigrateTask 查询指定迁移任务的返回体
type MigrateTask struct {
	// 任务名称
	Name string `json:"name"`
	// 任务类型,创建时必选,更新时可选
	Type string `json:"type"`
	// 操作系统类型,分为WINDOWS和LINUX,创建时必选,更新时可选
	OsType string `json:"os_type"`
	// 迁移任务id
	Id string `json:"id"`
	// 进程优先级
	//   0:低
	//   1:标准(默认)
	//   2:高
	Priority int `json:"priority"`
	// 迁移完成后是否启动目的端服务器
	//   true:启动
	//   false:停止
	StartTargetServer bool `json:"start_target_server"`
	// 企业项目id
	EnterpriseProjectId string `json:"enterprise_project_id"`
	// 目的端服务器的IP地址。
	//   公网迁移时请填写弹性IP地址
	//   专线迁移时请填写私有IP地址
	MigrationIp string `json:"migration_ip"`
	// 目的端服务器的区域名称
	Region string `json:"region_name"`
	// 目的端服务器的区域ID
	RegionID string `json:"region_id"`
	// 目的端服务器所在项目名称
	Project string `json:"project_name"`
	// 目的端服务器所在项目ID
	ProjectID string `json:"project_id"`
	// 模板ID
	VmTemplateId string `json:"vm_template_id"`
	// 源端服务器信息
	SourceServer SourceServer `json:"source_server"`
	// 目的端服务器信息
	TargetServer TaskTargetServer `json:"target_server"`
	// 任务状态
	State string `json:"state"`

	// 连接状态
	Connected bool `json:"connected"`
	// 迁移速率,单位:MB/S
	MigrateSpeed float32 `json:"migrate_speed"`
	// 压缩率
	CompressRate float32 `json:"compress_rate"`
	// 错误信息
	ErrorMsg string `json:"error_json"`

	// 任务创建时间
	CreateDate int `json:"create_date"`
	// 任务开始时间
	StartDate int `json:"start_date"`
	// 任务结束时间
	FinishDate int `json:"finish_date"`
	// 预估完成时间
	EstimateCompleteTime int `json:"estimate_complete_time"`
	// 迁移剩余时间(秒)
	RemainSeconds int `json:"remain_seconds"`
	// 任务总耗时
	TotalTime int `json:"total_time"`
	// 目的端的快照id
	TargetSnapshotId string `json:"target_snapshot_id"`
	// 克隆服务器信息
	CloneServer CloneServer `json:"clone_server"`
	// 任务包含的子任务列表
	SubTasks []SubTask `json:"sub_tasks"`
}

// SourceServer 源端服务器
type SourceServer struct {
	// 源端在SMS数据库中的ID
	Id string `json:"id"`
	// 源端服务器ip,注册源端时必选,更新非必选
	Ip string `json:"ip"`
	// 用来区分不同源端服务器的名称
	Name string `json:"name"`
	// 源端主机名,注册源端必选,更新非必选
	Hostname string `json:"hostname"`
	// 源端服务器的OS类型,分为Windows和Linux,注册必选,更新非必选
	OsType string `json:"os_type"`
	// 操作系统版本,注册必选,更新非必选
	OsVersion string `json:"os_version"`
	// 源端服务器启动类型,如BIOS或者UEFI
	Firmware string `json:"firmware"`
	// CPU个数,单位vCPU
	CpuQuantity int `json:"cpu_quantity"`
	// 内存大小,单位MB
	Memory int `json:"memory"`
	// 源端服务器的磁盘信息
	Disks []ServerDisk `json:"disks"`
	// Linux 必选,源端的Btrfs信息。如果源端不存在Btrfs,则为[]
	BtrfsList []BtrfsFileSystem `json:"btrfs_list"`
	// 源端服务器的网卡信息
	Networks []NetWork `json:"networks"`
	// 租户的domainId
	DomainId string `json:"domain_id"`
	// 是否安装rsync组件,Linux系统此参数为必选
	HasRsync bool `json:"has_rsync"`
	// Linux场景必选,源端是否是半虚拟化
	Paravirtualization bool `json:"paravirtualization"`
	// Linux必选,裸设备列表
	RawDevices string `json:"raw_devices"`
	// Windows 必选,是否缺少驱动文件
	DriverFiles bool `json:"driver_files"`
	// Windows必选,是否存在不正常服务
	SystemServices bool `json:"system_services"`
	// Windows必选,权限是否满足要求
	AccountRights bool `json:"account_rights"`
	// Linux必选,系统引导类型,BOOT_LOADER(GRUB/LILO)
	BootLoader string `json:"boot_loader"`
	// Windows必选,系统目录
	SystemDir string `json:"system_dir"`
	// Linux必选,如果没有卷组,输入[]
	VolumeGroups []VolumeGroups `json:"volume_groups"`
	// Agent版本
	AgentVersion string `json:"agent_version"`
}

// ServerDisk 磁盘信息
type ServerDisk struct {
	// 磁盘名称
	Name string `json:"name"`
	// 磁盘的分区类型,添加源端时源端磁盘必选
	PartitionStyle string `json:"partition_style"`
	// 磁盘类型
	DeviceType string `json:"device_use"`
	// 磁盘总大小,以字节为单位
	Size int64 `json:"size"`
	// 磁盘已使用大小,以字节为单位
	UsedSize int64 `json:"used_size"`
	// 磁盘上的物理分区信息
	PhysicalVolumes []PhysicalVolume `json:"physical_volumes"`
	// 是否为系统盘
	OsDisk bool `json:"os_disk"`
	// Linux系统 目的端ECS中与源端关联的磁盘名称
	RelationName string `json:"relation_name"`
}

// PhysicalVolume 使用大小
type PhysicalVolume struct {
	// 分区类型,普通分区,启动分区,系统分区
	DeviceType string `json:"device_use"`
	// 文件系统类型
	FileSystem string `json:"file_system"`
	// 顺序
	Index int `json:"index"`
	// 挂载点
	MountPoint string `json:"mount_point"`
	// 名称,windows表示盘符,Linux表示设备号
	Name string `json:"name"`
	// 大小
	Size int64 `json:"size"`
	// 使用大小
	UsedSize int64 `json:"used_size"`
	// GUID,可从源端查询
	UUID string `json:"uuid"`
	// 每个cluster大小
	SizePerCluster int `json:"size_per_cluster"`
}

// BtrfsFileSystem btrfs分区类型
type BtrfsFileSystem struct {
	// 文件系统名称
	Name string `json:"name"`
	// 文件系统标签,若无标签为空字符串
	Label string `json:"label"`
	// 文件系统的uuid
	UUID string `json:"uuid"`
	// btrfs包含的设备名称
	Device string `json:"device"`
	// 文件系统数据占用大小
	Size int64 `json:"size"`
	// btrfs节点大小
	Nodesize int `json:"nodesize"`
	// 扇区大小
	Sectorsize int `json:"sectorsize"`
	// 数据配置(RAD)
	DataProfile string `json:"data_profile"`
	// 文件系统配置(RAD)
	SystemProfile string `json:"system_profile"`
	// 元数据配置(RAD)
	MetadataProfile string `json:"metadata_profile"`
	// Btrfs文件系统信息
	GlobalReserve1 string `json:"global_reserve1"`
	// Btrfs卷已使用空间大小
	UsedSize int64 `json:"g_vol_used_size"`
	// 默认子卷ID
	DefaultSubvolid string `json:"default_subvolid"`
	// 默认子卷名称
	DefaultSubvolName string `json:"default_subvol_name"`
	// 默认子卷挂载路径/BTRFS文件系统的挂载路径
	DefaultSubvolMountpath string `json:"default_subvol_mountpath"`
	// 子卷信息
	Subvolumn []BtrfsSubvolumn `json:"subvolumn"`
}

// BtrfsSubvolumn btrfs子卷信息
type BtrfsSubvolumn struct {
	// 父卷的uuid
	UUID string `json:"uuid"`
	// 子卷是否为快照
	IsSnapshot string `json:"is_snapshot"`
	// 子卷的id
	SubvolId string `json:"subvol_id"`
	// 父卷id
	ParentId string `json:"parent_id"`
	// 子卷的名称
	SubvolName string `json:"subvol_name"`
	// 子卷的挂载路径
	SubvolMountPath string `json:"subvol_mount_path"`
}

// NetWork 网卡实体类
type NetWork struct {
	// 网卡的名称
	Name string `json:"name"`
	// 该网卡绑定的IP
	Ip string `json:"ip"`
	// 掩码
	Netmask string `json:"netmask"`
	// 网关
	Gateway string `json:"gateway"`
	// Linux必选,网卡的MTU
	Mtu int `json:"mtu"`
	// Mac地址
	Mac string `json:"mac"`
	// 数据库Id
	Id string `json:"id"`
}

// VolumeGroups 逻辑卷组信息
type VolumeGroups struct {
	// Pv信息
	Components string `json:"components"`
	// 剩余空间
	FreeSize int64 `json:"free_size"`
	// lv信息
	LogicalVolumes []LogicalVolumes `json:"logical_volumes"`
	// 名称
	Name string `json:"name"`
	// 大小
	Size int64 `json:"size"`
}

// LogicalVolumes 逻辑卷信息
type LogicalVolumes struct {
	// 块数量
	BlockCount int `json:"block_count"`
	// 块大小
	BlockSize int `json:"block_size"`
	// 文件系统
	FileSystem string `json:"file_system"`
	// inode数量
	InodeSize int `json:"inode_size"`
	// 挂载点
	MountPoint string `json:"mount_point"`
	// 名称
	Name string `json:"name"`
	// 大小
	Size int64 `json:"size"`
	// 使用大小
	UsedSize int64 `json:"used_size"`
	// 剩余空间
	FreeSize int64 `json:"free_size"`
}

// TaskTargetServer 目的端服务器
type TaskTargetServer struct {
	// 目的端在SMS数据库中的ID
	Id string `json:"id"`
	// 目的端服务器ID,自动创建虚拟机不需要这个参数
	VMID string `json:"vm_id"`
	// 目的端服务器的名称
	Name string `json:"name"`
	// 目的端服务器ip
	Ip string `json:"ip"`
	// 源端服务器的OS类型,分为Windows和Linux,注册必选,更新非必选
	OsType string `json:"os_type"`
	// 操作系统版本,注册必选,更新非必选
	OsVersion string `json:"os_version"`
	// Windows必选,系统目录
	SystemDir string `json:"system_dir"`
	// 目的端磁盘信息,一般和源端保持一致
	Disks []TargetDisk `json:"disks"`
	// lvm信息,一般和源端保持一致
	VolumeGroups []VolumeGroups `json:"volume_groups"`
	// Linux 必选,源端的Btrfs信息。如果源端不存在Btrfs,则为[]
	BtrfsList []string `json:"btrfs_list"`
	// 目的端代理镜像磁盘id
	ImageDiskId string `json:"image_disk_id"`
	// 目的端回滚快照id
	CutoveredSnapshotIds string `json:"cutovered_snapshot_ids"`
}

// TargetDisk 目的端磁盘
type TargetDisk struct {
	// 判断是普通分区,启动分区还是系统分区
	DeviceType string `json:"device_use"`
	// 磁盘id
	DiskId string `json:"disk_id"`
	// 磁盘名称
	Name string `json:"name"`
	// 逻辑卷信息
	PhysicalVolumes []TargetPhysicalVolumes `json:"physical_volumes"`
	// 大小
	Size int64 `json:"size"`
	// 已使用大小
	UsedSize int64 `json:"used_size"`
}

// TargetPhysicalVolumes 物理分区
type TargetPhysicalVolumes struct {
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

// CloneServer 克隆服务器类
type CloneServer struct {
	// 克隆服务器ID
	VMID string `json:"vm_id"`
	// 克隆虚拟机的名称
	Name string `json:"name"`
	// 克隆错误信息
	CloneError string `json:"clone_error"`
	// 克隆状态
	CloneState string `json:"clone_state"`
	// 克隆错误信息描述
	ErrorMsg string `json:"error_msg"`
}

// SubTask 修改任务进度
type SubTask struct {
	// 子任务名称
	Name string `json:"name"`
	// 子任务的进度,取值为0-100之间的整数
	Progress int `json:"progress"`
	// 子任务开始时间
	StartDate int `json:"start_date"`
	// 子任务结束时间(如果子任务还没有结束,则为空)
	EndDate int `json:"end_date"`
	// 迁移速率,Mbit/s
	MigrateSpeed float32 `json:"migrate_speed"`
	// 触发子任务的用户操作名称
	UserOp string `json:"user_op"`
	// 迁移或同步时,具体的迁移详情
	ProcessTrace string `json:"process_trace"`
}
