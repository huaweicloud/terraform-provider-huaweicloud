package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostcheckClusterRequestBody struct {

	// API版本，默认为v3
	ApiVersion string `json:"apiVersion"`

	// 资源类型
	Kind string `json:"kind"`

	Spec *PostcheckSpec `json:"spec"`
}

func (o PostcheckClusterRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostcheckClusterRequestBody struct{}"
	}

	return strings.Join([]string{"PostcheckClusterRequestBody", string(data)}, " ")
}
