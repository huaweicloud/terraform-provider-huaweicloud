package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 文件路径
type FilePath struct {
}

func (o FilePath) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FilePath struct{}"
	}

	return strings.Join([]string{"FilePath", string(data)}, " ")
}
