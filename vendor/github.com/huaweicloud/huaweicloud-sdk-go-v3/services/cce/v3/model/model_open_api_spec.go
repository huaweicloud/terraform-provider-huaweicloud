package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OpenApiSpec struct {
	Spec *OpenApiSpecSpec `json:"spec,omitempty"`
}

func (o OpenApiSpec) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OpenApiSpec struct{}"
	}

	return strings.Join([]string{"OpenApiSpec", string(data)}, " ")
}
