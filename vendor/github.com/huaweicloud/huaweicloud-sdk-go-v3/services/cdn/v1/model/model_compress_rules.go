package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CompressRules struct {

	// GZIP压缩类型（目前只支持 gzip）
	CompressType *string `json:"compress_type,omitempty"`

	// GZIP压缩文件类型（文件后缀竖线分割，如：.js|.html|.css|.xml）
	CompressFileType *string `json:"compress_file_type,omitempty"`
}

func (o CompressRules) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CompressRules struct{}"
	}

	return strings.Join([]string{"CompressRules", string(data)}, " ")
}
