package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ConfigurationRisks 配置风险项来源
type ConfigurationRisks struct {

	// 组件名称
	Package *string `json:"package,omitempty"`

	// 涉及文件路径
	SourceFile *string `json:"sourceFile,omitempty"`

	// 节点信息
	NodeMsg *string `json:"nodeMsg,omitempty"`

	// 参数值
	Field *string `json:"field,omitempty"`

	// 修改操作类型
	Operation *string `json:"operation,omitempty"`

	// 原始值
	OriginalValue *string `json:"originalValue,omitempty"`

	// 当前值
	Value *string `json:"value,omitempty"`
}

func (o ConfigurationRisks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationRisks struct{}"
	}

	return strings.Join([]string{"ConfigurationRisks", string(data)}, " ")
}
