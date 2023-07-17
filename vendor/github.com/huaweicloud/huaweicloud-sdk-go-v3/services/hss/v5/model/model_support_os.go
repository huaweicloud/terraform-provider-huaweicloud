package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SupportOs 支持的操作系统，包含如下:   - Linux ：支持Linux系统   - Windows : 支持Windows系统
type SupportOs struct {
}

func (o SupportOs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SupportOs struct{}"
	}

	return strings.Join([]string{"SupportOs", string(data)}, " ")
}
