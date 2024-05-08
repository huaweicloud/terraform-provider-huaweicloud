package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FileSize 文件大小
type FileSize struct {
}

func (o FileSize) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileSize struct{}"
	}

	return strings.Join([]string{"FileSize", string(data)}, " ")
}
