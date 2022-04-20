package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BrandBrokens struct {
	// 接收字节数

	RecBytes *[]float64 `json:"recBytes,omitempty"`
	// 发送字节数

	SentBytes *[]float64 `json:"sentBytes,omitempty"`
}

func (o BrandBrokens) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BrandBrokens struct{}"
	}

	return strings.Join([]string{"BrandBrokens", string(data)}, " ")
}
