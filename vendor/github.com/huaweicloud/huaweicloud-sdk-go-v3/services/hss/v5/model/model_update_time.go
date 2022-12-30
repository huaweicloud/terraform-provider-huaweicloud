package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 更新时间，毫秒
type UpdateTime struct {
}

func (o UpdateTime) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTime struct{}"
	}

	return strings.Join([]string{"UpdateTime", string(data)}, " ")
}
