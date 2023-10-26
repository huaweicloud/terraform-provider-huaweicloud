package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// VolumeFlavor 磁盘规格信息
type VolumeFlavor struct {

	// 引擎版本
	EngineVersion string `json:"engine_version"`

	// 磁盘规格码
	Code string `json:"code"`

	// 磁盘类型
	Type string `json:"type"`

	// 磁盘大小
	Size string `json:"size"`

	// 订购周期
	Period []string `json:"period"`
}

func (o VolumeFlavor) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VolumeFlavor struct{}"
	}

	return strings.Join([]string{"VolumeFlavor", string(data)}, " ")
}
