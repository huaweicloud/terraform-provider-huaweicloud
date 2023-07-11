package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// HandleStatus 处理状态，包含如下:   - unhandled ：未处理   - handled : 已处理
type HandleStatus struct {
}

func (o HandleStatus) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HandleStatus struct{}"
	}

	return strings.Join([]string{"HandleStatus", string(data)}, " ")
}
