package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 目的端服务器
type TaskTargetServer struct {
	// 目的端在SMS数据库中的ID

	Id *string `json:"id,omitempty"`
	// 目的端服务器ID，自动创建虚拟机不需要这个参数

	VmId *string `json:"vm_id,omitempty"`
	// 目的端服务器的名称

	Name *string `json:"name,omitempty"`
	// 目的端服务器ip

	Ip *string `json:"ip,omitempty"`
	// 源端服务器的OS类型，分为Windows和Linux，注册必选，更新非必选

	OsType *TaskTargetServerOsType `json:"os_type,omitempty"`
	// 操作系统版本，注册必选，更新非必选

	OsVersion *string `json:"os_version,omitempty"`
	// Windows必选，系统目录

	SystemDir *string `json:"system_dir,omitempty"`
	// 目的端磁盘信息，一般和源端保持一致

	Disks []TargetDisk `json:"disks"`
	// lvm信息，一般和源端保持一致

	VolumeGroups *[]VolumeGroups `json:"volume_groups,omitempty"`
	// Linux 必选，源端的Btrfs信息。如果源端不存在Btrfs，则为[]

	BtrfsList *[]string `json:"btrfs_list,omitempty"`
	// 目的端代理镜像磁盘id

	ImageDiskId *string `json:"image_disk_id,omitempty"`
	// 目的端回滚快照id

	CutoveredSnapshotIds *string `json:"cutovered_snapshot_ids,omitempty"`
}

func (o TaskTargetServer) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskTargetServer struct{}"
	}

	return strings.Join([]string{"TaskTargetServer", string(data)}, " ")
}

type TaskTargetServerOsType struct {
	value string
}

type TaskTargetServerOsTypeEnum struct {
	WINDOWS TaskTargetServerOsType
	LINUX   TaskTargetServerOsType
}

func GetTaskTargetServerOsTypeEnum() TaskTargetServerOsTypeEnum {
	return TaskTargetServerOsTypeEnum{
		WINDOWS: TaskTargetServerOsType{
			value: "WINDOWS",
		},
		LINUX: TaskTargetServerOsType{
			value: "LINUX",
		},
	}
}

func (c TaskTargetServerOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskTargetServerOsType) UnmarshalJSON(b []byte) error {
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
