package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSnapshotConfigRequest Request Object
type DeleteSnapshotConfigRequest struct {

	// 直播流播放域名
	Domain string `json:"domain"`

	// 应用名称
	AppName string `json:"app_name"`
}

func (o DeleteSnapshotConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSnapshotConfigRequest struct{}"
	}

	return strings.Join([]string{"DeleteSnapshotConfigRequest", string(data)}, " ")
}
