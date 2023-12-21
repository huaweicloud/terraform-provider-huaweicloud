package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OpenApiSpecSpec 集群访问的地址
type OpenApiSpecSpec struct {
	Eip *EipSpec `json:"eip,omitempty"`

	// 是否动态创建
	IsDynamic *bool `json:"IsDynamic,omitempty"`
}

func (o OpenApiSpecSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OpenApiSpecSpec struct{}"
	}

	return strings.Join([]string{"OpenApiSpecSpec", string(data)}, " ")
}
