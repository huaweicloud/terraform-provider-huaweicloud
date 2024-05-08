package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Compress 智能压缩。
type Compress struct {

	// 智能压缩开关（on：开启，off：关闭）。
	Status string `json:"status"`

	// 智能压缩类型（gzip：gzip压缩，br：brotli压缩）。
	Type *string `json:"type,omitempty"`

	// 压缩格式，内容总长度不可超过200个字符，  多种格式用“,”分割，每组内容不可超过50个字符， 开启状态下，首次传空时默认值为.js,.html,.css,.xml,.json,.shtml,.htm，否则为上次设置的结果。
	FileType *string `json:"file_type,omitempty"`
}

func (o Compress) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Compress struct{}"
	}

	return strings.Join([]string{"Compress", string(data)}, " ")
}
