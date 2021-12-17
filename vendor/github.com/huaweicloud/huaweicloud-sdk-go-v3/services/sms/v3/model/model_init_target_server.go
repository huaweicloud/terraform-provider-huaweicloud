package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 推荐的目的端服务器配置
type InitTargetServer struct {
	// 推荐的目的端服务器的磁盘信息

	Disks []DiskIntargetServer `json:"disks"`
}

func (o InitTargetServer) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InitTargetServer struct{}"
	}

	return strings.Join([]string{"InitTargetServer", string(data)}, " ")
}
