package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// This is a auto create Body Object
type UpdateTaskSpeedReq struct {
	// 当前上报进度的子任务名称，子任务名称包括： 创建虚拟机 CREATE_CLOUD_SERVER 配置安全通道 SSL_CONFIG 挂载代理镜像 ATTACH_AGENT_IMAGE 卸载载代理镜像 DETTACH_AGENT_IMAGE Linux分区格式化 FORMAT_DISK_LINUX Linux分区格式化(文件级级） FORMAT_DISK_LINUX_FILE Linux分区格式化(块级） FORMAT_DISK_LINUX_BLOCK Windows分区格式化 FORMAT_DISK_WINDOWS Linux文件级数据迁移 MIGRATE_LINUX_FILE, Linux块级数据迁移 MIGRATE_LINUX_BLOCK Windows块级数据迁移 MIGRATE_WINDOWS_BLOCK 克隆一个虚拟机 CLONE_VM Linux文件级数据同步 SYNC_LINUX_FILE Linux块级数据同步 SYNC_LINUX_BLOCK Windows块级数据同步 SYNC_WINDOWS_BLOCK Linux配置修改 CONFIGURE_LINUX Linux配置修改(块级）CONFIGURE_LINUX_BLOCK Linux配置修改（文件级） CONFIGURE_LINUX_FILE Windows配置修改 CONFIGURE_WINDOWS

	SubtaskName UpdateTaskSpeedReqSubtaskName `json:"subtask_name"`
	// 当前上报的子任务的最新百分比进度

	Progress int32 `json:"progress"`
	// 当前任务已经复制的数据量大小（B）

	Replicatesize int64 `json:"replicatesize"`
	// 当前任务的总迁移数据大小

	Totalsize int64 `json:"totalsize"`
	// 实施迁移速率，单位Mb/s

	MigrateSpeed float64 `json:"migrate_speed"`
}

func (o UpdateTaskSpeedReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskSpeedReq struct{}"
	}

	return strings.Join([]string{"UpdateTaskSpeedReq", string(data)}, " ")
}

type UpdateTaskSpeedReqSubtaskName struct {
	value string
}

type UpdateTaskSpeedReqSubtaskNameEnum struct {
	CREATE_CLOUD_SERVER     UpdateTaskSpeedReqSubtaskName
	SSL_CONFIG              UpdateTaskSpeedReqSubtaskName
	ATTACH_AGENT_IMAGE      UpdateTaskSpeedReqSubtaskName
	DETTACH_AGENT_IMAGE     UpdateTaskSpeedReqSubtaskName
	FORMAT_DISK_LINUX       UpdateTaskSpeedReqSubtaskName
	FORMAT_DISK_LINUX_FILE  UpdateTaskSpeedReqSubtaskName
	FORMAT_DISK_LINUX_BLOCK UpdateTaskSpeedReqSubtaskName
	FORMAT_DISK_WINDOWS     UpdateTaskSpeedReqSubtaskName
	MIGRATE_LINUX_FILE      UpdateTaskSpeedReqSubtaskName
	MIGRATE_LINUX_BLOCK     UpdateTaskSpeedReqSubtaskName
	MIGRATE_WINDOWS_BLOCK   UpdateTaskSpeedReqSubtaskName
	CLONE_VM                UpdateTaskSpeedReqSubtaskName
	SYNC_LINUX_FILE         UpdateTaskSpeedReqSubtaskName
	SYNC_LINUX_BLOCK        UpdateTaskSpeedReqSubtaskName
	SYNC_WINDOWS_BLOCK      UpdateTaskSpeedReqSubtaskName
	CONFIGURE_LINUX         UpdateTaskSpeedReqSubtaskName
	CONFIGURE_LINUX_BLOCK   UpdateTaskSpeedReqSubtaskName
	CONFIGURE_LINUX_FILE    UpdateTaskSpeedReqSubtaskName
	CONFIGURE_WINDOWS       UpdateTaskSpeedReqSubtaskName
}

func GetUpdateTaskSpeedReqSubtaskNameEnum() UpdateTaskSpeedReqSubtaskNameEnum {
	return UpdateTaskSpeedReqSubtaskNameEnum{
		CREATE_CLOUD_SERVER: UpdateTaskSpeedReqSubtaskName{
			value: "CREATE_CLOUD_SERVER",
		},
		SSL_CONFIG: UpdateTaskSpeedReqSubtaskName{
			value: "SSL_CONFIG",
		},
		ATTACH_AGENT_IMAGE: UpdateTaskSpeedReqSubtaskName{
			value: "ATTACH_AGENT_IMAGE",
		},
		DETTACH_AGENT_IMAGE: UpdateTaskSpeedReqSubtaskName{
			value: "DETTACH_AGENT_IMAGE",
		},
		FORMAT_DISK_LINUX: UpdateTaskSpeedReqSubtaskName{
			value: "FORMAT_DISK_LINUX",
		},
		FORMAT_DISK_LINUX_FILE: UpdateTaskSpeedReqSubtaskName{
			value: "FORMAT_DISK_LINUX_FILE",
		},
		FORMAT_DISK_LINUX_BLOCK: UpdateTaskSpeedReqSubtaskName{
			value: "FORMAT_DISK_LINUX_BLOCK",
		},
		FORMAT_DISK_WINDOWS: UpdateTaskSpeedReqSubtaskName{
			value: "FORMAT_DISK_WINDOWS",
		},
		MIGRATE_LINUX_FILE: UpdateTaskSpeedReqSubtaskName{
			value: "MIGRATE_LINUX_FILE",
		},
		MIGRATE_LINUX_BLOCK: UpdateTaskSpeedReqSubtaskName{
			value: "MIGRATE_LINUX_BLOCK",
		},
		MIGRATE_WINDOWS_BLOCK: UpdateTaskSpeedReqSubtaskName{
			value: "MIGRATE_WINDOWS_BLOCK",
		},
		CLONE_VM: UpdateTaskSpeedReqSubtaskName{
			value: "CLONE_VM",
		},
		SYNC_LINUX_FILE: UpdateTaskSpeedReqSubtaskName{
			value: "SYNC_LINUX_FILE",
		},
		SYNC_LINUX_BLOCK: UpdateTaskSpeedReqSubtaskName{
			value: "SYNC_LINUX_BLOCK",
		},
		SYNC_WINDOWS_BLOCK: UpdateTaskSpeedReqSubtaskName{
			value: "SYNC_WINDOWS_BLOCK",
		},
		CONFIGURE_LINUX: UpdateTaskSpeedReqSubtaskName{
			value: "CONFIGURE_LINUX",
		},
		CONFIGURE_LINUX_BLOCK: UpdateTaskSpeedReqSubtaskName{
			value: "CONFIGURE_LINUX_BLOCK",
		},
		CONFIGURE_LINUX_FILE: UpdateTaskSpeedReqSubtaskName{
			value: "CONFIGURE_LINUX_FILE",
		},
		CONFIGURE_WINDOWS: UpdateTaskSpeedReqSubtaskName{
			value: "CONFIGURE_WINDOWS",
		},
	}
}

func (c UpdateTaskSpeedReqSubtaskName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateTaskSpeedReqSubtaskName) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
