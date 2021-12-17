package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// btrfs分区类型
type BtrfsFileSystem struct {
	// 文件系统名称

	Name string `json:"name"`
	// 文件系统标签，若无标签为空字符串

	Label string `json:"label"`
	// 文件系统的uuid

	Uuid string `json:"uuid"`
	// btrfs包含的设备名称

	Device string `json:"device"`
	// 文件系统数据占用大小

	Size int64 `json:"size"`
	// btrfs节点大小

	Nodesize int64 `json:"nodesize"`
	// 扇区大小

	Sectorsize int32 `json:"sectorsize"`
	// 数据配置（RAD）

	DataProfile string `json:"data_profile"`
	// 文件系统配置（RAD）

	SystemProfile string `json:"system_profile"`
	// 元数据配置（RAD）

	MetadataProfile string `json:"metadata_profile"`
	// Btrfs文件系统信息

	GlobalReserve1 string `json:"global_reserve1"`
	// Btrfs卷已使用空间大小

	GVolUsedSize int64 `json:"g_vol_used_size"`
	// 默认子卷ID

	DefaultSubvolid string `json:"default_subvolid"`
	// 默认子卷名称

	DefaultSubvolName string `json:"default_subvol_name"`
	// 默认子卷挂载路径/BTRFS文件系统的挂载路径

	DefaultSubvolMountpath string `json:"default_subvol_mountpath"`
	// 子卷信息

	Subvolumn []BtrfsSubvolumn `json:"subvolumn"`
}

func (o BtrfsFileSystem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BtrfsFileSystem struct{}"
	}

	return strings.Join([]string{"BtrfsFileSystem", string(data)}, " ")
}
