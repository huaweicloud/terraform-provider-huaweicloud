package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PackageRequestArgs struct {

	// 录制播放相关配置
	Record *[]RecordRequestArgs `json:"record,omitempty"`

	// 时移播放相关配置
	Timeshift *[]TimeshiftRequestArgs `json:"timeshift,omitempty"`

	// 直播播放相关配置
	Live *[]LiveRequestArgs `json:"live,omitempty"`
}

func (o PackageRequestArgs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PackageRequestArgs struct{}"
	}

	return strings.Join([]string{"PackageRequestArgs", string(data)}, " ")
}
