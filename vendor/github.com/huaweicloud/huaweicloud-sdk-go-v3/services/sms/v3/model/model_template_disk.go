package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 磁盘模板
type TemplateDisk struct {
	// 磁盘序号，从0开始

	Index int32 `json:"index"`
	// 磁盘名称

	Name string `json:"name"`
	// 磁盘类型，同volumetype字段

	Disktype string `json:"disktype"`
	// 磁盘大小，单位：GB

	Size int64 `json:"size"`
}

func (o TemplateDisk) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TemplateDisk struct{}"
	}

	return strings.Join([]string{"TemplateDisk", string(data)}, " ")
}
