package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 文件属性
type FileAttr struct {
}

func (o FileAttr) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FileAttr struct{}"
	}

	return strings.Join([]string{"FileAttr", string(data)}, " ")
}
