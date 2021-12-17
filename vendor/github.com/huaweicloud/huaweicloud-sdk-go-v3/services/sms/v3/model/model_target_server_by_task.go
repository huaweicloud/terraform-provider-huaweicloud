package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 目的端虚拟机信息
type TargetServerByTask struct {
	// btrfs信息，数据从源端获取

	BtrfsList *[]BtrfsFileSystem `json:"btrfs_list,omitempty"`
	// 磁盘信息

	Disks []TargetDisks `json:"disks"`
	// 名称

	Name string `json:"name"`
	// 虚拟机id

	VmId string `json:"vm_id"`
	// 卷组，数据从源端获取

	VolumeGroups *[]VolumeGroups `json:"volume_groups,omitempty"`
}

func (o TargetServerByTask) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetServerByTask struct{}"
	}

	return strings.Join([]string{"TargetServerByTask", string(data)}, " ")
}
