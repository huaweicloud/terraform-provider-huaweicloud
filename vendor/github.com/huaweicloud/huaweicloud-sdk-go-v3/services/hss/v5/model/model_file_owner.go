package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FileOwner 文件属主
type FileOwner struct {
}

func (o FileOwner) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileOwner struct{}"
	}

	return strings.Join([]string{"FileOwner", string(data)}, " ")
}
