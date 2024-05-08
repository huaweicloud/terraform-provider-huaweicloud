package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FileCtime 文件创建时间
type FileCtime struct {
}

func (o FileCtime) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileCtime struct{}"
	}

	return strings.Join([]string{"FileCtime", string(data)}, " ")
}
