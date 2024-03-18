package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PrecheckClusterRequestBody struct {

	// API版本，默认为v3
	ApiVersion string `json:"apiVersion"`

	// 资源类型，默认为PreCheckTask
	Kind string `json:"kind"`

	Spec *PrecheckSpec `json:"spec"`
}

func (o PrecheckClusterRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PrecheckClusterRequestBody struct{}"
	}

	return strings.Join([]string{"PrecheckClusterRequestBody", string(data)}, " ")
}
