package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 文件哈希
type FileHash struct {
}

func (o FileHash) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileHash struct{}"
	}

	return strings.Join([]string{"FileHash", string(data)}, " ")
}
