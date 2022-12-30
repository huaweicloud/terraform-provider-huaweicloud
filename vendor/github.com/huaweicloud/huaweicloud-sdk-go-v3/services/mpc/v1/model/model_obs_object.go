package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ObsObject struct {

	// 对象的key
	FileName *string `json:"file_name,omitempty"`

	// 文件大小
	Size *int64 `json:"size,omitempty"`

	// 文件的最后修改时间
	LastModified *string `json:"last_modified,omitempty"`
}

func (o ObsObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsObject struct{}"
	}

	return strings.Join([]string{"ObsObject", string(data)}, " ")
}
