package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// NodeSelector 节点标签选择器，匹配Kubernetes中nodeSelector相关约束
type NodeSelector struct {

	// 标签键
	Key string `json:"key"`

	// 标签值列表
	Value *[]string `json:"value,omitempty"`

	// 标签逻辑运算符
	Operator string `json:"operator"`
}

func (o NodeSelector) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NodeSelector struct{}"
	}

	return strings.Join([]string{"NodeSelector", string(data)}, " ")
}
