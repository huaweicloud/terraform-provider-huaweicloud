package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FileMtime 文件更新时间
type FileMtime struct {
}

func (o FileMtime) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileMtime struct{}"
	}

	return strings.Join([]string{"FileMtime", string(data)}, " ")
}
