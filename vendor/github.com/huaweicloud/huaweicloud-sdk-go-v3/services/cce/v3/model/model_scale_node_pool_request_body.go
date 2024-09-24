package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ScaleNodePoolRequestBody 伸缩节点池的请求体
type ScaleNodePoolRequestBody struct {

	// API类型，固定值“NodePool”。
	Kind string `json:"kind"`

	// API版本，固定值“v3”。
	ApiVersion string `json:"apiVersion"`

	Spec *ScaleNodePoolSpec `json:"spec"`
}

func (o ScaleNodePoolRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleNodePoolRequestBody struct{}"
	}

	return strings.Join([]string{"ScaleNodePoolRequestBody", string(data)}, " ")
}
