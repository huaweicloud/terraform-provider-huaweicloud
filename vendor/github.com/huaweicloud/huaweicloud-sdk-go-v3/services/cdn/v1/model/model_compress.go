package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 智能压缩。
type Compress struct {

	// 智能压缩开关（on：开启，off：关闭）。
	Status string `json:"status"`

	// 智能压缩类型（gzip：gzip压缩，br：br压缩）。
	Type *string `json:"type,omitempty"`
}

func (o Compress) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Compress struct{}"
	}

	return strings.Join([]string{"Compress", string(data)}, " ")
}
