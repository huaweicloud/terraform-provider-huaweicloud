package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 存储盘信息。
type MysqlVolumeInfo struct {
	// 磁盘类型。

	Type string `json:"type"`
	// 已使用磁盘大小，单位GB。

	Size string `json:"size"`
}

func (o MysqlVolumeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlVolumeInfo struct{}"
	}

	return strings.Join([]string{"MysqlVolumeInfo", string(data)}, " ")
}
