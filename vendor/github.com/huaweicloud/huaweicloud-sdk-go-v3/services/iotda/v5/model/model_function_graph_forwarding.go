package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FunctionGraphForwarding 函数工作流转发配置信息
type FunctionGraphForwarding struct {

	// **参数说明**：函数的URN（Uniform Resource Name），唯一标识函数。
	FuncUrn string `json:"func_urn"`

	// **参数说明**：函数名称。
	FuncName string `json:"func_name"`
}

func (o FunctionGraphForwarding) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FunctionGraphForwarding struct{}"
	}

	return strings.Join([]string{"FunctionGraphForwarding", string(data)}, " ")
}
