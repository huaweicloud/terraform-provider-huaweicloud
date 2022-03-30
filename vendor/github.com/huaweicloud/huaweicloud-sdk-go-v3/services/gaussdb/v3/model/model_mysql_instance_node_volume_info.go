package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 存储盘信息。
type MysqlInstanceNodeVolumeInfo struct {
	// 磁盘类型。

	Type string `json:"type"`
	// 已使用磁盘大小，单位GB。

	Used string `json:"used"`
	// 包周期购买的存储空间大小，单位GB。

	Size int64 `json:"size"`
}

func (o MysqlInstanceNodeVolumeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlInstanceNodeVolumeInfo struct{}"
	}

	return strings.Join([]string{"MysqlInstanceNodeVolumeInfo", string(data)}, " ")
}
