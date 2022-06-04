package tasks

import "github.com/chnsz/golangsdk"

// CreateOpts 创建任务的参数
type CreateOpts struct {
	// 任务名称
	Name string `json:"name" required:"true"`
	// 任务类型: MIGRATE_FILE, MIGRATE_BLOCK
	Type string `json:"type" required:"true"`
	// 操作系统类型: WINDOWS, LINUX
	OsType string `json:"os_type" required:"true"`
	// region的名称
	Region string `json:"region_name" required:"true"`
	// region id
	RegionID string `json:"region_id" required:"true"`
	// 项目名称
	Project string `json:"project_name" required:"true"`
	// 项目id
	ProjectID string `json:"project_id" required:"true"`
	// 源端服务器信息
	SourceServer SourceServerRequest `json:"source_server" required:"true"`
	// 目的端虚拟机信息
	TargetServer TargetServerRequest `json:"target_server" required:"true"`
	// 进程优先级, 0:低, 1:标准, 2:高
	Priority int `json:"priority,omitempty"`
	// 迁移ip,如果是自动创建虚拟机,不需要此参数
	MigrationIp string `json:"migration_ip,omitempty"`
	// 自动创建虚拟机使用模板
	VmTemplateId string `json:"vm_template_id,omitempty"`
	// 是否使用已有虚拟机
	ExistServer *bool `json:"exist_server,omitempty"`
	// 迁移后是否启动目的端虚拟机
	StartServer *bool `json:"start_target_server,omitempty"`
	// 是否使用公网ip
	UsePublicIp *bool `json:"use_public_ip,omitempty"`
	// 复制或者同步后是否会继续持续同步,不添加则默认是false
	Syncing *bool `json:"syncing,omitempty"`
}

// SourceServerRequest 源端服务器信息
type SourceServerRequest struct {
	// 源端服务器id
	Id string `json:"id" required:"true"`
}

// TargetServerRequest 目的端虚拟机信息
type TargetServerRequest struct {
	// 虚拟机名称
	Name string `json:"name,omitempty"`
	// 虚拟机id, 如果是自动创建虚拟机,不需要此参数
	VMID string `json:"vm_id,omitempty"`
	// 磁盘信息
	Disks []DiskRequest `json:"disks,omitempty"`
	// 卷组,数据从源端获取
	VolumeGroups []VGRequest `json:"volume_groups,omitempty"`
	// btrfs信息,数据从源端获取
	Btrfs []BtrfsFileSystem `json:"btrfs_list,omitempty"`
}

// DiskRequest 目的端磁盘信息
type DiskRequest struct {
	// 名称,根据磁盘顺序设置为disk X
	Name string `json:"name" required:"true"`
	// 大小
	Size int64 `json:"size" required:"true"`
	// 磁盘类型,普通磁盘,OS所在磁盘,BOOT所在磁盘
	DeviceType string `json:"device_use,omitempty"`
	// 磁盘id,自动创建虚拟机不用设置
	DiskId string `json:"disk_id,omitempty"`
	// 物理卷信息
	PhysicalVolumes []PVRequest `json:"physical_volumes,omitempty"`
	// 使用大小
	UsedSize int64 `json:"used_size,omitempty"`
}

// PVRequest 使用大小
type PVRequest struct {
	// 分区类型,普通分区,启动分区,系统分区
	DeviceType string `json:"device_use,omitempty"`
	// 文件系统类型
	FileSystem string `json:"file_system,omitempty"`
	// 顺序
	Index *int `json:"index,omitempty"`
	// 挂载点
	MountPoint string `json:"mount_point,omitempty"`
	// 名称,windows表示盘符,Linux表示设备号
	Name string `json:"name,omitempty"`
	// 大小
	Size int64 `json:"size,omitempty"`
	// 使用大小
	UsedSize int64 `json:"used_size,omitempty"`
	// GUID,可从源端查询
	UUID string `json:"uuid,omitempty"`
}

// VGRequest 逻辑卷组信息
type VGRequest struct {
	// 名称
	Name string `json:"name" required:"true"`
	// 大小
	Size int64 `json:"size,omitempty"`
	// Pv信息
	Components string `json:"components,omitempty"`
	// 剩余空间
	FreeSize int64 `json:"free_size,omitempty"`
	// lv信息
	LogicalVolumes []LVRequest `json:"logical_volumes,omitempty"`
}

// LVRequest 逻辑卷信息
type LVRequest struct {
	// 名称
	Name string `json:"name" required:"true"`
	// 大小
	Size int64 `json:"size" required:"true"`
	// 块数量
	BlockCount int `json:"block_count,omitempty"`
	// 块大小
	BlockSize int `json:"block_size,omitempty"`
	// 文件系统
	FileSystem string `json:"file_system"`
	// 挂载点
	MountPoint string `json:"mount_point"`
	// inode数量
	InodeSize int `json:"inode_size,omitempty"`
	// 使用大小
	UsedSize int64 `json:"used_size,omitempty"`
	// 剩余空间
	FreeSize int `json:"free_size,omitempty"`
}

// Create 创建迁移任务
func Create(c *golangsdk.ServiceClient, opts *CreateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var rst golangsdk.Result
	_, rst.Err = c.Post(rootURL(c), b, &rst.Body, nil)

	var r CreateResp
	if err := rst.ExtractInto(&r); err != nil {
		return "", err
	}

	return r.ID, nil
}

// ActionOpts is an object to manage migration tasks
type ActionOpts struct {
	// Operation specifies the operation to be performed on the task. The value can be:
	// start, stop, collect_log, test, clone_test, restart, sync_failed_rollback
	Operation string `json:"operation" required:"true"`
	// Param specifies the operation parameters
	Param map[string]string `json:"param,omitempty"`
}

// Action is the method to manage migration tasks
func Action(c *golangsdk.ServiceClient, id string, opts ActionOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(actionURL(c, id), b, nil, nil)
	return err
}

// Get 查询指定ID模板信息
func Get(c *golangsdk.ServiceClient, id string) (*MigrateTask, error) {
	var rst golangsdk.Result
	_, rst.Err = c.Get(taskURL(c, id), &rst.Body, nil)

	var r MigrateTask
	err := rst.ExtractInto(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Delete 删除指定ID的迁移任务
func Delete(c *golangsdk.ServiceClient, id string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(taskURL(c, id), nil)
	return &r
}
